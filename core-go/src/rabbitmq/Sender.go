package rabbitmq

import (
	"encoding/json"
	"github.com/evan-beraben/aptus/core-go/src/messages"
	"github.com/streadway/amqp"
	"log"
)

type SenderMq struct {
	Url      string
	Username string
	Password string
	Queue    string
	Conn     *amqp.Connection
}

func SendMessage(sender *SenderMq, message messages.Message) bool {

	if sender.Conn == nil {
		creds := sender.Username + ":" + sender.Password
		dial := "amqp://" + creds + "@" + sender.Url + ":5672/"

		log.Printf("Dialing %s", dial)
		conn, err := amqp.Dial(dial)

		if err != nil {
			log.Println(err)
			log.Fatal("Failed to connect to to RabbitMq on ", sender.Url)
			return false
		}

		sender.Conn = conn;
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		log.Println("Failed to send message, invalid object.")

		return false
	}

	ch, err := sender.Conn.Channel()
	if err != nil{
		log.Println(err)
		log.Println("Failed to get channel")
		sender.Conn = nil
		return false
	}

	/*
	err = ch.Confirm(false)
	if err != nil {
		log.Println(err)
		log.Fatal("Can not confirm channel")
	}*/

	//confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	err = ch.Publish("", sender.Queue, true, false, amqp.Publishing{Body: []byte(data)})
	if err != nil {
		log.Println(err)
		log.Println("Failed to publish to queue.")
		sender.Conn = nil
		return false
	}
	//confirmResult := <-confirms
/*
	if confirmResult.Ack {
		log.Println("Message ACK")
		return true
	} else {
		log.Println("Message Failed")
		return false
	}
*/
	return true
}


