package config

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Worker struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Telegram struct {
	BotToken string `yaml:"bot_token"`
}

type RabbitMQ struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Vhost    string `yaml:"vhost"`
}

type Sentry struct {
	DSN         string `yaml:"dsn"`
	Environment string `yaml:"environment"`
	Release     string `yaml:"release"`
}

type Logging struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}
type AppConfigStruct struct {
	Server   Server   `yaml:"server"`
	Worker   Worker   `yaml:"worker"`
	Database DB       `yaml:"database"`
	Telegram Telegram `yaml:"telegram"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
	Sentry   Sentry   `yaml:"sentry"`
	Logging  Logging  `yaml:"logging"`
}
