package config

// Config holds the application configuration
type Config struct {
	RabbitMQURL  string
	Queue1       string
	Queue2       string
	HTTPService1 string
	HTTPService2 string
}

func LoadConfig() *Config {
	return &Config{
		RabbitMQURL:  "amqp://user:password@localhost:5672/myvhost",
		Queue1:       "dummy_data_queue_1",
		Queue2:       "dummy_data_queue_2",
		HTTPService1: "http://localhost:8081/api1",
		HTTPService2: "http://localhost:8081/api2",
	}
}
