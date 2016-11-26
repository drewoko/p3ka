package core

import (
	"log"
	"time"
	"strconv"
	"strings"

	"github.com/mvdan/xurls"
	"github.com/parnurzeal/gorequest"
)

func MessageProcessor(messagesInputChannel chan Msg, messagesDeleteChannel chan Msg, db *DataBase, config *Config) {

	defer func() {
		if r := recover(); r != nil {
			log.Print("Recovering MessageProcessor", r)
		}
	}()

	for ;; {
		var message Msg
		var deleted_message Msg

		select {
			case message = <-messagesInputChannel:
				go processMessage(message, db, config)
			case deleted_message = <-messagesDeleteChannel:
				go deleteMessage(deleted_message, db)
		}
	}
}

func deleteMessage(message Msg, db *DataBase) {

	msg := db.GetMessageById(message.Id)
	db.SetDeleted((msg["id"]));
}

func processMessage(message Msg, db *DataBase, config *Config)  {

	request := gorequest.New().Timeout(time.Second*4)

	if(!ContainsString(config.BannedUsers, message.Name)) {

		links := xurls.Strict.FindAllString(message.Text, -1)

		if(len(links) != 0) {
			log.Print("Received message: ", message)
		}

		for _, link := range links {
			if (!strings.Contains(link, "#noP3KA") && !db.IsExists(message.Name, link)) {
				resp, _, errs :=  request.Head(link).End();
				if(len(errs) == 0) {
					if(resp.StatusCode == 200) {
						log.Print("Processing file: ", link, resp)
						if contentType, ok := resp.Header["Content-Type"]; ok {
							if contentLength, ok := resp.Header["Content-Length"]; ok {
								if(len(errs) == 0 && isFileAllowed(contentType[0], contentLength[0])) {

									db.AddRow(message.Id, message.Name, link, isMature(message.Text));
								}
							}
						}
					}
				} else {
					log.Print("Feiled to process image. Reason: ", errs[0])
				}
			}
		}
	}
}

func isMature(text string) bool {
	//sorry for that
	text = strings.ToLower(text)
	return strings.Contains(text, "+18") || strings.Contains(text, "18+") || strings.Contains(text, "+21")  || strings.Contains(text, "21+") || strings.Contains(text, "[spoiler")
}

func isFileAllowed(content_type string, content_length string) bool {

	i64, err := strconv.ParseInt(content_length, 10, 64)

	if(err != nil) {
		return false
	}

	return (i64 < 5000000 && i64 > 100 &&
		(content_type == "image/jpeg" || content_type == "image/png" || content_type == "image/gif"))
}