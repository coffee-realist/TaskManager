package main

import (
	"TaskPublisher/internal/api"
	"TaskPublisher/internal/broker"
	"TaskPublisher/internal/domain/config"
	"TaskPublisher/internal/domain/service"
	"TaskPublisher/internal/storage"
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := initConfig(); err != nil {
		log.Error("Error initializing configs", "error", err)
	}

	db := storage.NewSqlConnection(config.DataBaseConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.SSLMode"),
	})

	natsConfig := config.NatsConfig{
		Host:       viper.GetString("nats.host"),
		StreamName: viper.GetString("nats.streamName"),
		KVBucket:   viper.GetString("nats.KVBucket"),
	}
	nc, err := nats.Connect(natsConfig.Host)
	if err != nil {
		log.Error("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		log.Error("failed to create JetStream context: %w", err)
	}

	if err := initNats(js, natsConfig); err != nil {
		nc.Close()
		log.Error("Failed to initialize NATS resources", "error", err)
		os.Exit(1)
	}

	natsBroker, err := broker.NewNatsBrokerService(js, nc, natsConfig)
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
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
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
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	return viper.ReadInConfig()
}

func initNats(js jetstream.JetStream, cfg config.NatsConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:      cfg.StreamName,
		Subjects:  []string{"created.*", "finished.*"},
		Retention: jetstream.LimitsPolicy,
		Storage:   jetstream.FileStorage,
	})
	if err != nil && !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	_, err = js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:  cfg.KVBucket,
		History: 10,
	})
	if err != nil && !errors.Is(err, jetstream.ErrBucketExists) {
		return fmt.Errorf("failed to create KV bucket: %w", err)
	}

	return nil
}
