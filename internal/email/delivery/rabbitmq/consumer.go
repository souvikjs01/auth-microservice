package rabbitmq

import (
	"context"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/souvikjs01/auth-microservice/internal/email"
	"github.com/souvikjs01/auth-microservice/pkg/logger"
	"github.com/streadway/amqp"
)

var (
	IncomingMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "emails_incoming_rebbitmq_messages_total",
		Help: "The total number of incoming rabbitmq messages",
	})

	SuccessMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "emails_success_incoming_rabbitmq_message_total",
		Help: "The total number of successfull mabbitmq messages",
	})

	FailureMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "emails_failure_rabbitmq_message_total",
		Help: "The total number failed rabbitmq messages",
	})
)

// email rabbitmq consumer
type EmailConsumer struct {
	amqpConn *amqp.Connection
	logger   logger.Logger
	emailUC  email.EmailUsecase
}

func NewEmailConsumer(amqpConn *amqp.Connection, logger logger.Logger, emailUC email.EmailUsecase) *EmailConsumer {
	return &EmailConsumer{
		amqpConn: amqpConn,
		logger:   logger,
		emailUC:  emailUC,
	}
}

func (c *EmailConsumer) CreateChannel(exchange, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "amqpConn.Channel")
	}

	c.logger.Infof("Declaring exchange: %s", exchange)

	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Channel.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "ch.QueueDeclare")
	}

	c.logger.Infof("Declaring queue, binding it to exchange - queue: %v, messagecount: %v, consumerCount: %v, exchange: %v, bindingKey: %v", queue.Name, queue.Messages, queue.Consumers, exchange, bindingKey)

	err = ch.QueueBind(queueName, bindingKey, exchange, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Channel.QueueBind")
	}

	c.logger.Infof("Queue bound to exchange, starting to consumer from queue, consumerTag: %v", consumerTag)

	err = ch.Qos(1, 0, false)
	if err != nil {
		return nil, errors.Wrap(err, "Channel.Qos")
	}

	return ch, nil
}

func (c *EmailConsumer) worker(ctx context.Context, messages <-chan amqp.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()

	for delivery := range messages {
		span, ctx := opentracing.StartSpanFromContext(ctx, "EmailConsumer.worker")

		c.logger.Infof("processDeliveries deliveryTag %v", delivery.DeliveryTag)

		IncomingMessages.Inc()

		err := c.emailUC.SendEmail(ctx, delivery)

		if err != nil {
			if err := delivery.Reject(false); err != nil {
				c.logger.Errorf("delivery.rejected: %v", err)
			}
			FailureMessages.Inc()
			c.logger.Errorf("emailUC.SendEmail, failed to process delivery: %v", err)
			span.Finish()
		} else {
			SuccessMessages.Inc()
			err = delivery.Ack(false)
			if err != nil {
				c.logger.Errorf("delivery.Ack: %v", err)
			}
			span.Finish()
		}
	}

	c.logger.Info("Delivery channel closed")
}

// start new rabbitmq consumer
func (c *EmailConsumer) StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {
	ch, err := c.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return err
	}
	defer ch.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go c.worker(ctx, deliveries, wg)
	}

	wg.Wait()

	return nil
}
