package config

import (
	"os"
)

//Конфиг сервиса
func ConfigSetup() {
	os.Setenv("DB_USERNAME", "phunda")
	os.Setenv("DB_PASSWORD", "098908")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "l0wb")
	os.Setenv("DB_SSLMODE", "disable")

	os.Setenv("NATS_HOSTS", "localhost:4223")
	os.Setenv("NATS_CLUSTER_ID", "test-cluster")
	os.Setenv("NATS_CLIENT_ID", "phunda")
	os.Setenv("NATS_SUBJECT", "go.test-phunda")
	os.Setenv("NATS_DURABLE_NAME", "Replica-1")
	os.Setenv("NATS_ACK_WAIT_SECONDS", "30")
}
