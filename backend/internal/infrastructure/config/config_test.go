package config

import (
	"os"
	"testing"
)

func setRequired(t *testing.T) {
	t.Helper()
	t.Setenv("POSTGRES_PASSWORD", "pgpass")
	t.Setenv("RABBITMQ_PASSWORD", "rmqpass")
	t.Setenv("TELEGRAM_BOT_TOKEN", "bot-token")
}

func TestLoad_AppliesDefaultsAndReadsRequired(t *testing.T) {
	setRequired(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// defaults
	if cfg.Server.Host != "0.0.0.0" || cfg.Server.Port != 8080 {
		t.Fatalf("server defaults wrong: %+v", cfg.Server)
	}
	if cfg.Database.Host != "localhost" || cfg.Database.Port != 5432 {
		t.Fatalf("db defaults wrong: %+v", cfg.Database)
	}
	if cfg.RabbitMQ.Vhost != "/" {
		t.Fatalf("rabbitmq vhost default wrong: %q", cfg.RabbitMQ.Vhost)
	}
	// required values read from env
	if cfg.Database.Password != "pgpass" {
		t.Fatalf("db password not read: %q", cfg.Database.Password)
	}
	if cfg.Telegram.BotToken != "bot-token" {
		t.Fatalf("bot token not read: %q", cfg.Telegram.BotToken)
	}
}

func TestLoad_OverridesFromEnv(t *testing.T) {
	setRequired(t)
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("POSTGRES_HOST", "db.internal")
	t.Setenv("RABBITMQ_VHOST", "wishlist")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Fatalf("server port override wrong: %d", cfg.Server.Port)
	}
	if cfg.Database.Host != "db.internal" {
		t.Fatalf("db host override wrong: %q", cfg.Database.Host)
	}
	if cfg.RabbitMQ.Vhost != "wishlist" {
		t.Fatalf("vhost override wrong: %q", cfg.RabbitMQ.Vhost)
	}
}

func TestLoad_MissingRequiredErrors(t *testing.T) {
	t.Setenv("POSTGRES_PASSWORD", "pgpass")
	t.Setenv("RABBITMQ_PASSWORD", "rmqpass")
	_ = os.Unsetenv("TELEGRAM_BOT_TOKEN")

	if _, err := Load(); err == nil {
		t.Fatal("expected error when TELEGRAM_BOT_TOKEN is missing")
	}
}
