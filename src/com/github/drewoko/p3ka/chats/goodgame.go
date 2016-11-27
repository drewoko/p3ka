package chats

import (
	"log"
	"time"
	"sync"
	"encoding/json"

	Core "../core"

	"github.com/gorilla/websocket"
)

type GoodGameSocketStorage struct {
	sync.Mutex
	wsClient *websocket.Conn
}

func (s *GoodGameSocketStorage) writeMessage(request []byte) error{
	s.Lock();
	defer s.Unlock()

	return s.wsClient.WriteMessage(websocket.TextMessage, request)
}

func initGoodGame(messages_channel chan Core.Msg, messages_delete_channel chan Core.Msg, db *Core.DataBase, config *Core.Config) {

	log.Println("Connecting to GoodGame.ru")

	wsClient, _, err := websocket.DefaultDialer.Dial(config.GoodGameHost, nil)

	if err != nil {
		log.Println("Failed to connect to GoodGame.ru", err)
		time.Sleep(time.Second * 5)
		initGoodGame(messages_channel, messages_delete_channel, db, config);
		return
	}

	socket := &GoodGameSocketStorage{
		wsClient: wsClient,
	}

	plainMessageChan := make(chan []byte)
	channelChan := make(chan string)
	quitChat := make(chan bool)

	pingTicker := time.NewTicker(time.Second * 5)
	channelRefresh := time.NewTicker(time.Hour)
	defer func() {
		close(plainMessageChan)
		close(channelChan)
		pingTicker.Stop()
		channelRefresh.Stop()
	}()

	go func() {
		for {
			messageType, message, err := wsClient.ReadMessage()

			if err != nil {
				log.Println("Disconnected from GoodGame.ru")
				quitChat <- true
				time.Sleep(time.Second * 5)
				initGoodGame(messages_channel, messages_delete_channel, db, config);
				return
			}

			if messageType == websocket.TextMessage {
				plainMessageChan <- message
			}
		}
	}()

	counter := 0;

	for {
		var plainMessage []byte
		select {
		case plainMessage = <- plainMessageChan:

			message := GoodGameStruct{}

			json.Unmarshal(plainMessage, &message)

			if message.Type == "welcome" {
				log.Println("Connected to GoodGame.ru")
				joinToSavedChannels(socket, db);
				go requestChannels(socket, 0, config.GoodGameMaxRequestSize)
			} else if message.Type == "channels_list" {
				go processChannels(&counter, socket, config, message, channelChan)
			} else if message.Type == "remove_message" {
				messages_delete_channel <- Core.Msg{Id: int64(message.Data["message_id"].(float64)), Source: "goodgame",}
			} else if message.Type == "message" {
				messages_channel <- Core.Msg{
					Id: int64(message.Data["message_id"].(float64)),
					Text: message.Data["text"].(string),
					Name: message.Data["user_name"].(string),
					Channel: message.Data["channel_id"],
					Source: "goodgame",
				}
			}
		case channel := <-channelChan:
			db.AddRowGGChannel(channel)
		case <- channelRefresh.C:
			go requestChannels(socket, 0, config.GoodGameMaxRequestSize)
		case <- pingTicker.C:
			sendPing(socket)
		case <- quitChat:
			return
		}
	}
}

func processChannels(counter *int, socket *GoodGameSocketStorage, config *Core.Config, message GoodGameStruct, channelChan chan string) {
	var channelInterface interface{}
	intCounter := 0;
	for _, channelInterface = range message.Data["channels"].([]interface{}) {
		channel := channelInterface.(map[string]interface{})["channel_id"].(string)
		channelChan <- channel
		joinToChannel(socket, channel)
		intCounter++;
	}
	*counter = *counter + intCounter

	if intCounter == config.GoodGameMaxRequestSize {
		go requestChannels(socket, intCounter-1, config.GoodGameMaxRequestSize)
	}
}

func joinToSavedChannels(socket *GoodGameSocketStorage, db *Core.DataBase) {
	for _, channel := range db.GetGGChannels() {
		joinToChannel(socket, channel)
	}
}

func sendPing(socket *GoodGameSocketStorage) {
	sentMessage(socket, GoodGameStruct{
		Type: "ping",
		Data: map[string]interface{}{},
	})
}

func joinToChannel(socket *GoodGameSocketStorage, channel interface{}) {
	sentMessage(socket, GoodGameStruct{
		Type: "join",
		Data: map[string]interface{}{"channel_id": channel, "hidden": false},
	})
}

func requestChannels(socket *GoodGameSocketStorage, start int, count int) {
	sentMessage(socket, GoodGameStruct{
		Type: "get_channels_list",
		Data: map[string]interface{}{"start": start, "count": count},
	})
}

func sentMessage(socket *GoodGameSocketStorage, messageStruct GoodGameStruct) {

	request, err := json.Marshal(messageStruct)

	if err != nil {
		log.Println("Failed to create JSON", err)
		return
	}

	err = socket.writeMessage(request)

	if err != nil {
		log.Println(err)
	}
}

type GoodGameStruct struct {
	Type string	`json:"type"`
	Data map[string]interface{} `json:"data"`
}