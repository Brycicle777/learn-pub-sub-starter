module github.com/bootdotdev/learn-pub-sub-starter

go 1.23.6

require github.com/rabbitmq/amqp091-go v1.10.0 // indirect

require internal/pubsub v0.0.0

replace internal/pubsub => ./internal/pubsub

require internal/routing v0.0.0

replace internal/routing => ./internal/routing