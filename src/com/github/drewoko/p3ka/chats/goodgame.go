package chats

import (
	"log"
	"time"
	"encoding/json"

	Core "../core"

	"github.com/gorilla/websocket"
)

func initGoodGame(messages_channel chan Core.Msg, messages_delete_channel chan Core.Msg, db *Core.DataBase, config *Core.Config) {

	log.Println("Connecting to GoodGame.ru")

	wsClient, _, err := websocket.DefaultDialer.Dial(config.GoodGameHost, nil)

	if err != nil {
		log.Println("Failed to connect to GoodGame.ru", err)
		time.Sleep(time.Second * 5)
		initGoodGame(messages_channel, messages_delete_channel, db, config);
		return
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
				joinToSavedChannels(wsClient, db);
				go requestChannels(wsClient, 0, config.GoodGameMaxRequestSize)
			} else if message.Type == "channels_list" {
				go processChannels(&counter, wsClient, config, message, channelChan)
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
			go requestChannels(wsClient, 0, config.GoodGameMaxRequestSize)
		case <- pingTicker.C:
			sendPing(wsClient)
		case <- quitChat:
			return
		}
	}
}

func processChannels(counter *int, wsClient *websocket.Conn, config *Core.Config, message GoodGameStruct, channelChan chan string) {
	var channelInterface interface{}
	intCounter := 0;
	for _, channelInterface = range message.Data["channels"].([]interface{}) {
		channel := channelInterface.(map[string]interface{})["channel_id"].(string)
		channelChan <- channel
		joinToChannel(wsClient, channel)
		intCounter++;
	}
	*counter = *counter + intCounter

	if intCounter == config.GoodGameMaxRequestSize {
		go requestChannels(wsClient, intCounter-1, config.GoodGameMaxRequestSize)
	}
}

func joinToSavedChannels(wsClient *websocket.Conn, db *Core.DataBase) {
	for _, channel := range db.GetGGChannels() {
		joinToChannel(wsClient, channel)
	}
}

func sendPing(wsClient *websocket.Conn) {
	sentMessage(wsClient, GoodGameStruct{
		Type: "ping",
		Data: map[string]interface{}{},
	})
}

func joinToChannel(wsClient *websocket.Conn, channel interface{}) {
	sentMessage(wsClient, GoodGameStruct{
		Type: "join",
		Data: map[string]interface{}{"channel_id": channel, "hidden": false},
	})
}

func requestChannels(wsClient *websocket.Conn, start int, count int) {
	sentMessage(wsClient, GoodGameStruct{
		Type: "get_channels_list",
		Data: map[string]interface{}{"start": start, "count": count},
	})
}

func sentMessage(wsClient *websocket.Conn, messageStruct GoodGameStruct) {

	request, err := json.Marshal(messageStruct)

	if err != nil {
		log.Println("Failed to create JSON", err)
		return
	}
	err = wsClient.WriteMessage(websocket.TextMessage, request)

	if err != nil {
		log.Println(err)
	}
}

type GoodGameStruct struct {
	Type string	`json:"type"`
	Data map[string]interface{} `json:"data"`
}