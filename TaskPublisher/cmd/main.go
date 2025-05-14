package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/api"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/broker"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/domain/config"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/domain/service"
	"github.com/coffee-realist/TaskManager/TaskPublisher/internal/storage"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

// @title Task Publisher API
// @version 1.0
// @description API для управления публикацией задач
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := initConfig(); err != nil {
		log.Error("Error initializing configs", "error", err)
		os.Exit(1)
	}

	// Инициализация БД
	db := storage.NewSqlConnection(config.DataBaseConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.SSLMode"),
	})

	// Подключение к NATS
	natsConfig := config.NatsConfig{
		Host:       viper.GetString("nats.host"),
		StreamName: viper.GetString("nats.streamName"),
		KVBucket:   viper.GetString("nats.KVBucket"),
	}

	nc, err := nats.Connect(natsConfig.Host,
		nats.MaxReconnects(5),
		nats.ReconnectWait(2*time.Second))
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		os.Exit(1)
	}

	// Ждем подключения
	for i := 0; i < 10; i++ {
		if nc.IsConnected() {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if !nc.IsConnected() {
		log.Error("NATS connection timeout")
		os.Exit(1)
	}

	// Создаем JetStream
	js, err := jetstream.New(nc)
	if err != nil {
		log.Error("Failed to create JetStream context", "error", err)
		nc.Close()
		os.Exit(1)
	}
	jsClient, err := nc.JetStream()
	if err != nil {
		nc.Close()
		log.Error("failed to create JetStream client: %w", err)
	}

	// Инициализация NATS ресурсов
	if err := initNats(js, natsConfig); err != nil {
		log.Error("Failed to initialize NATS resources", "error", err)
		nc.Close()
		os.Exit(1)
	}

	natsBroker, err := broker.NewNatsBrokerService(js, jsClient, nc, natsConfig)
	if err != nil {
		nc.Close()
		log.Error("failed to init NATS broker: %w", err)
	}

	serv, err := service.NewService(db, natsBroker)
	if err != nil {
		log.Error("Error initializing service", "error", err)
	}
	handlers := api.NewHandler(serv)

	srv := new(api.Server)
	go func() {
		if err := srv.Run(viper.GetString("publisher.port"), handlers.InitRoutes()); err != nil {
			log.Error("Server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.ShutDown(ctx); err != nil {
		log.Error("Server Shutdown Failed", "error", err)
	}
	if err := nc.Drain(); err != nil {
		log.Error("NATS connection drain failed", "error", err)
	}
	log.Info("Server exited properly")
}

func initConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	return viper.ReadInConfig()
}

func initNats(js jetstream.JetStream, cfg config.NatsConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Конфигурация потока с хранением сообщений
	streamConfig := jetstream.StreamConfig{
		Name:      cfg.StreamName,
		Subjects:  []string{"created.*", "finished.*"},
		Retention: jetstream.LimitsPolicy,
		Storage:   jetstream.FileStorage,
		MaxAge:    7 * 24 * time.Hour,     // Хранить сообщения 7 дней
		MaxBytes:  1 * 1024 * 1024 * 1024, // 1GB максимальный размер
		Discard:   jetstream.DiscardOld,
	}

	// Создаем или обновляем поток
	stream, err := js.CreateOrUpdateStream(ctx, streamConfig)
	if err != nil {
		return fmt.Errorf("failed to create/update stream: %w", err)
	}

	// Проверяем конфигурацию потока
	info, err := stream.Info(ctx)
	if err != nil {
		return fmt.Errorf("failed to get stream info: %w", err)
	}

	slog.Info("Stream configuration",
		"name", info.Config.Name,
		"subjects", info.Config.Subjects,
		"retention", info.Config.Retention.String(),
		"max_age", info.Config.MaxAge,
		"max_bytes", info.Config.MaxBytes)

	// Создаем KV bucket (если нужно)
	_, err = js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  cfg.KVBucket,
		History: 10,
	})
	if err != nil && !errors.Is(err, jetstream.ErrBucketExists) {
		return fmt.Errorf("failed to create KV bucket: %w", err)
	}

	return nil
}
