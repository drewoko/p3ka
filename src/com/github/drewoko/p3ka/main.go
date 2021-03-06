package main

import (
	"flag"
	"log"
	"runtime"
	"strings"
	"sync"

	"gopkg.in/magiconair/properties.v1"

	Chats "./chats"
	Core "./core"
)

/**
Before running
go-bindata -o core/bindata.go -pkg core static/dist/*
*/

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("Starting P3KA application")

	configurationFile := flag.String("properties", "application.properties", "Properties file")
	flag.Parse()

	propertyFile := properties.MustLoadFile(*configurationFile, properties.UTF8)

	config := &Core.Config{
		Database:               propertyFile.GetString("database", "p3ka.db"),
		Port:                   propertyFile.GetString("port", "8080"),
		Static:                 propertyFile.GetString("static", "/static"),
		BannedUsers:            strings.Split(propertyFile.GetString("banned-users", ""), ","),
		ExcludedUsers:          strings.Split(propertyFile.GetString("exclude-from-rationg", ""), ","),
		Peka2TvHost:            propertyFile.GetString("peka2tv-host", "chat.funstream.tv"),
		Peka2TvPort:            propertyFile.GetInt("peka2tv-port", 80),
		GoodGameHost:           propertyFile.GetString("goodgame-host", "ws://chat.goodgame.ru:8081/chat/websocket"),
		GoodGameMaxRequestSize: propertyFile.GetInt("goodgame-request-size", 50),
		Dev:               propertyFile.GetBool("dev", false),
		HttpResponseLimit: propertyFile.GetInt("http-limit", 51),
	}

	propertyFile.MustFlag(flag.CommandLine)

	messagesInputChannel := make(chan Core.Msg)
	messagesDeleteChannel := make(chan Core.Msg)

	db := new(Core.DataBase).Init(config.Database)

	defer func() {
		close(messagesInputChannel)
		close(messagesDeleteChannel)
		db.Close()
	}()

	for _, bannedUser := range config.BannedUsers {
		db.SetDeletedByUser(bannedUser)
	}

	var wg sync.WaitGroup
	wg.Add(5)

	go Core.MessageProcessor(messagesInputChannel, messagesDeleteChannel, db, config)
	go Core.ImageChecker(db)
	go Core.Web(db, config)

	Chats.InitChats(messagesInputChannel, messagesDeleteChannel, db, config)

	wg.Wait()

}
