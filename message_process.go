package main

import (
	"time"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/mvdan/xurls"
	"github.com/parnurzeal/gorequest"
)

func message_processor(messages_channel chan Msg, messages_delete_channel chan Msg, db *DataBase, config *Config) {

	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovering", r)
		}
	}()

	for ;; {
		var message Msg
		var deleted_message Msg

		select {
			case message = <-messages_channel:
				go process_message(message, db, config)
			case deleted_message = <-messages_delete_channel:
				go delete_message(deleted_message, db)
		}
	}
}

func delete_message(message Msg, db *DataBase) {

	msg := db.get_message_by_id(message.Id)
	db.set_deleted((msg["id"]));
}

func process_message(message Msg, db *DataBase, config *Config)  {

	request := gorequest.New().Timeout(time.Second*4)

	if(is_channel_public(message.Channel, request) && !contains_string(config.BannedUsers, message.Name)) {

		links := xurls.Strict.FindAllString(message.Text, -1)

		if(len(links) != 0) {
			log.Info("Received message: ", message)
		}

		for _, link := range links {
			if (!strings.Contains(link, "#noP3KA") && !db.is_exists(message.Name, link)) {
				resp, _, errs :=  request.Head(link).End();
				if(len(errs) == 0) {
					if(resp.StatusCode == 200) {
						log.Info("Processing file: ", link, resp)
						if contentType, ok := resp.Header["Content-Type"]; ok {
							if contentLength, ok := resp.Header["Content-Length"]; ok {
								if(len(errs) == 0 && is_file_allowed(contentType[0], contentLength[0])) {

									db.add_row(message.Id, message.Name, link, is_mature(message.Text));
								}
							}
						}
					}
				} else {
					log.Info("Feiled to process image. Reason: ", errs[0])
				}
			}
		}
	}
}

func is_mature(text string) bool {
	//sorry for that
	text = strings.ToLower(text)
	return strings.Contains(text, "+18") || strings.Contains(text, "18+") || strings.Contains(text, "+21")  || strings.Contains(text, "21+") || strings.Contains(text, "[spoiler")
}

func is_channel_public(channel string, request *gorequest.SuperAgent) bool {
	if(strings.Contains(channel, "room/")) {
		var room RoomResponse

		resp, body, errs := request.Post("http://funstream.tv/api/room").
			Send(`{"roomId":`+strings.Replace(channel, "room/", "", -1)+`}`).
			End()

		if(len(errs) == 0 && resp.StatusCode == 200) {
			err_un := json.Unmarshal([]byte(body), &room)
			if(err_un == nil) {
				return room.Mode == "public"
			}
		}
	}

	return true;
}

func is_file_allowed(content_type string, content_length string) bool {

	i64, err := strconv.ParseInt(content_length, 10, 64)

	if(err != nil) {
		return false
	}

	return (i64 < 5000000 && i64 > 100 &&
		(content_type == "image/jpeg" || content_type == "image/png" || content_type == "image/gif"))
}

type RoomResponse struct  {
	Mode string `json:"mode"`
}

func contains_string(s []string, e string) bool {
	for _, a := range s {
		if strings.ToLower(a) == strings.ToLower(e) {
			return true
		}
	}
	return false
}