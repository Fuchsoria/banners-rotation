package simpleproducer

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

var errPublish = errors.New("cannot publish message because channel isn't declared")

type AMQPMessage struct {
	Type         string `json:"type"`
	SlotID       string `json:"slot_id"`
	BannerID     string `json:"banner_id"`
	SocialDemoID string `json:"social_demo_id"`
	Date         string `json:"date"`
}

type RMQConnection interface {
	Channel() (*amqp.Channel, error)
}

type Producer struct {
	name    string
	conn    RMQConnection
	channel *amqp.Channel
}

func New(name string, conn RMQConnection) *Producer {
	return &Producer{name: name, conn: conn}
}

func (p *Producer) Connect() error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("cannot get channel, %w", err)
	}

	p.channel = ch

	_, err = ch.QueueDeclare(p.name, false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("cannot create queue, %w", err)
	}

	return nil
}

func (p *Producer) Publish(message AMQPMessage) error {
	if p.channel != nil {
		bytes, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("cannot marshall message, %w", err)
		}

		err = p.channel.Publish(
			"",     // exchange
			p.name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        bytes,
			})
		if err != nil {
			return fmt.Errorf("cannot publish message, %w", err)
		}

		return nil
	}

	return errPublish
}
