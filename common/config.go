package utils

import "os"

type CONFIG struct {
	ServiceName string
	DotLogs     string
}

type RABBITCONFIG struct {
	Uri          string
	StorageQueue string
	GatewayQueue string
}

var Config = CONFIG{
	ServiceName: getEnv("SERVICE_NAME", "gateway"),
	DotLogs:     getEnv("DOT_LOGS", "logs"),
}

var RabbitConfig = RABBITCONFIG{
	Uri:          getEnv("RABBIT_URI", "amqp://admin:admin@localhost:5672/"),
	StorageQueue: getEnv("STORAGE_QUEUE", "storage"),
	GatewayQueue: getEnv("GATEWAY_QUEUE", "gateway"),
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
