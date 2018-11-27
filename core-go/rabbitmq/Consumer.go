package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct{
	Url string
	Username string
	Password string
	Queue string
	OnReceive func(data []byte)
}

func NewConsumer(url string, user string, password string, queue string, onRecevie func(data []byte)) (*Consumer){
	var c Consumer

	c.Url = url
	c.Username = user
	c.Password = password
	c.Queue = queue
	c.OnReceive = onRecevie

	return &c
}



func StartReceive(c *Consumer){
	creds := c.Username + ":" + c.Password
	dial := "amqp://" + creds + "@" + c.Url + ":5672/"

	log.Printf("Dialing %s", dial)
	conn,err := amqp.Dial(dial)

	if err != nil{
		log.Fatal("Failed to connect to to RabbitMq on " + c.Url)
		return
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil{
		log.Fatal("Failed to open channel")
		return;
	}

	defer ch.Close()
	//ch.QueueDeclare(c.Queue, false, false, false,false,nil)

	msgs, err := ch.Consume(c.Queue,"",true,false,false, false, nil)
	if err != nil {
		log.Println("Faield to consume, creating queue")
		ch.QueueDeclare(c.Queue, true, false, false,false,nil)
		msgs, err = ch.Consume(c.Queue,"",true,false,false, false, nil)
		if err != nil {
			log.Printf("Failed to consume: %s" + c.Queue)
			return
		}
	}

	log.Printf("Waiting for messages on %s", c.Queue)

	for d := range msgs{
			log.Printf("Recevied message: %s", d.Body)
			c.OnReceive(d.Body)
	}

}