package main

import (
	"time"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

func init_listener(messages_channel chan Msg, messages_delete_channel chan Msg) {
	init_ws(messages_channel, messages_delete_channel)
}

func init_ws(messages_channel chan Msg, messages_delete_channel chan Msg) {

	log.Info("Trying to connect to Funstream.tv WS")

	ws_client, err := gosocketio.Dial(
		gosocketio.GetUrl("chat.funstream.tv", 80, false),
		&transport.WebsocketTransport{
			PingInterval:   5 * time.Second,
			PingTimeout:    10 * time.Second,
			ReceiveTimeout: 10 * time.Second,
			SendTimeout:    10 * time.Second,
			BufferSize:     1024 * 32,
		})

	if(err != nil) {
		log.Info("Failed to connect to funstreams.tv WS. Reason: ", err)
	}

	ws_client.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Info("Funstream.tv WS connected")

		err := ws_client.Emit("/chat/join", &Channel{"all"})

		if(err != nil) {
			log.Fatal("Failed to join channel. Reason: ", err)
		}
	})

	ws_client.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		ws_client.Close()
		log.Fatal("Funstream.tv WS disconnected")
	})

	ws_client.On("/chat/message", func(h *gosocketio.Channel, args Message) {
		messages_channel <- Msg{Id: args.Id, Text:args.Text, Name:args.From.Name, Channel:args.Channel}
	})

	ws_client.On("/chat/message/remove", func(h *gosocketio.Channel, args Message) {
		messages_delete_channel <- Msg{Id: args.Id}
	})
}

type Message struct {
	Id      int64    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	From	struct{
		Name	string `json:"name"`
			} `json:"from"`
}

type Channel struct {
	Channel string `json:"channel"`
}
