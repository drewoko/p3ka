package main

import (
	"os"
	"sync"
	"flag"

	"github.com/op/go-logging"
	"github.com/magiconair/properties"
	"strings"
)

var (
	log = logging.MustGetLogger("P3KA")
)

func main() {

	init_logger()

	log.Info("Starting P3KA application")

	prop_file := flag.String("properties", "application.properties", "Properties file")
	flag.Parse()

	p := properties.MustLoadFile(*prop_file, properties.UTF8)

	config := &Config {
		Database: p.GetString("database", "p3ka.db"),
		Port: p.GetString("port", "8080"),
		Static: p.GetString("static", "./static"),
		BannedUsers: strings.Split(p.GetString("banned-users", ""), ","),
		ExcludedUsers: strings.Split(p.GetString("exclude-from-rationg", ""), ","),
	}

	var wg sync.WaitGroup
	messages_channel := make(chan Msg)
	messages_delete_channel := make(chan Msg)

	db := new(DataBase).init(config.Database);
	defer db.db.Close()

	wg.Add(3)

	go init_listener(messages_channel, messages_delete_channel)
	go message_processor(messages_channel, messages_delete_channel, db, config)
	go web(db, config)
	go image_checker(db, config)

	wg.Wait()
}

func init_logger() {

	logging.SetBackend(
		logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(
			`%{color}%{time} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
		)),
	)
}

type Msg struct {
	Id int64
	Name string
	Text string
	Channel string
}

type Config struct {
	Database string
	Port string
	Static string
	BannedUsers []string
	ExcludedUsers []string
}