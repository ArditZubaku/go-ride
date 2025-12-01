package env

var RabbitMQ = struct {
	Username string
	Password string
	URI      string
}{
	Username: "RABBITMQ_DEFAULT_USER",
	Password: "RABBITMQ_DEFAULT_PASS",
	URI:      "RABBITMQ_CONNECTION_STRING",
}

var RabbitMQDefaults = struct {
	Username string
	Password string
	URI      string
}{
	Username: "guest",
	Password: "guest",
	URI:      "amqp://guest:guest@rabbitmq:5672",
}
