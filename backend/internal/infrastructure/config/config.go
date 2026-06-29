package config

// Конфигурация читается ИЗ ПЕРЕМЕННЫХ ОКРУЖЕНИЯ (12-factor). Полный список
// ключей и значения по умолчанию — в .env.example. Загрузка — config.Load().

type DB struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"wishlist"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Name     string `env:"POSTGRES_DB" envDefault:"wishlist"`
}

type Server struct {
	Host string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	Port int    `env:"SERVER_PORT" envDefault:"8080"`
}

type Worker struct {
	Host string `env:"WORKER_HOST" envDefault:"0.0.0.0"`
	Port int    `env:"WORKER_PORT" envDefault:"8090"`
}

type Telegram struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
}

type RabbitMQ struct {
	Host     string `env:"RABBITMQ_HOST" envDefault:"localhost"`
	Port     int    `env:"RABBITMQ_PORT" envDefault:"5672"`
	User     string `env:"RABBITMQ_USER" envDefault:"guest"`
	Password string `env:"RABBITMQ_PASSWORD,required"`
	Vhost    string `env:"RABBITMQ_VHOST" envDefault:"/"`
}

type Sentry struct {
	DSN         string `env:"SENTRY_DSN"`
	Environment string `env:"SENTRY_ENVIRONMENT" envDefault:"production"`
	Release     string `env:"SENTRY_RELEASE"`
}

type Logging struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
	File  string `env:"LOG_FILE"`
}

type AppConfigStruct struct {
	Server   Server
	Worker   Worker
	Database DB
	Telegram Telegram
	RabbitMQ RabbitMQ
	Sentry   Sentry
	Logging  Logging
}
