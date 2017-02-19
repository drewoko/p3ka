package chats

import (
	"log"
	"time"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	Core "../core"
)

func initPeka2Tv(messages_channel chan Core.Msg, messages_delete_channel chan Core.Msg, config *Core.Config) {

	log.Println("Trying to connect to Funstream.tv WS")

	ws_client, err := gosocketio.Dial(
		gosocketio.GetUrl(config.Peka2TvHost, 80, false),
		&transport.WebsocketTransport{
			PingInterval:   5 * time.Second,
			PingTimeout:    10 * time.Second,
			ReceiveTimeout: 10 * time.Second,
			SendTimeout:    10 * time.Second,
			BufferSize:     1024 * 32,
		})

	if err != nil {
		log.Println("Failed to connect to funstreams.tv WS. Reason: ", err)

		time.Sleep(time.Second * 10)
		initPeka2Tv(messages_channel, messages_delete_channel, config)
		return
	}

	quitChan := make(chan bool)

	ws_client.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Funstream.tv WS connected")

		err := ws_client.Emit("/chat/join", struct {
			Channel string `json:"channel"`
		}{Channel: "all"})

		if err != nil {
			log.Fatal("Failed to join channel. Reason: ", err)
		}
	})

	ws_client.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Println("Disconnected from Peka2.tv")
		quitChan <- true
	})

	ws_client.On("/chat/message", func(h *gosocketio.Channel, args Message) {
		messages_channel <- Core.Msg{
			Id:      args.Id,
			Text:    args.Text,
			Name:    args.From.Name,
			Channel: args.Channel,
			Source:  "peka2tv",
		}
	})

	ws_client.On("/chat/message/remove", func(h *gosocketio.Channel, args Message) {
		messages_delete_channel <- Core.Msg{Id: args.Id, Source: "peka2tv"}
	})

	for {
		select {
		case <-quitChan:
			close(quitChan)
			initPeka2Tv(messages_channel, messages_delete_channel, config)
			return
		}
	}
}

type Message struct {
	Id      int64  `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	From    struct {
		Name string `json:"name"`
	} `json:"from"`
}
