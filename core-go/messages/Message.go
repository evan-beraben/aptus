package messages

import(
	"github.com/satori/go.uuid"
)


type Message struct{
	Id string 			`json:"Id"`
	RoutingKey string 	`json:"RoutingKey"`
	Type string 		`json:"Type"`
	Body string			`json:"Body"`
}

func Build(routeKey string, messageType string, body string) Message{
	return Message{Id: uuid.Must(uuid.NewV4()).String(), RoutingKey: routeKey, Type: messageType, Body: body}
}