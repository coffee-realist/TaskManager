package config

type DataBaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

type NatsConfig struct {
	Host       string
	StreamName string
	KVBucket   string
}
