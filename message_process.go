package main

import (
	"time"
	"strconv"

	"github.com/mvdan/xurls"
	"github.com/parnurzeal/gorequest"
)

func message_processor(messages_channel chan Msg, messages_delete_channel chan Msg, db *DataBase, config *Config) {

	for ;; {
		var message Msg
		var deleted_message Msg

		select {
			case message = <-messages_channel:
				go process_message(message, db)
			case deleted_message = <-messages_delete_channel:
				go delete_message(deleted_message, db)
		}
	}
}

func delete_message(message Msg, db *DataBase) {

	msg := db.get_message_by_id(message.Id)
	db.set_deleted((msg["id"]));
}

func process_message(message Msg, db *DataBase)  {
	links := xurls.Strict.FindAllString(message.Text, -1)

	request := gorequest.New().Timeout(time.Second*4)

	for _, link := range links {
		if (!db.is_exists(message.Name, link)) {
			resp, _, errs :=  request.Head(link).End();
			if(len(errs) == 0) {
				if(resp.StatusCode == 200) {
					log.Info("Processing file: ", link, resp)
					if contentType, ok := resp.Header["Content-Type"]; ok {
						if contentLength, ok := resp.Header["Content-Length"]; ok {
							if(len(errs) == 0 && is_file_allowed(contentType[0], contentLength[0])) {
								db.add_row(message.Id, message.Name, link);
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

func is_file_allowed(content_type string, content_length string) bool {

	i64, err := strconv.ParseInt(content_length, 10, 64)

	if(err != nil) {
		return false
	}

	return (i64 < 5000000 && i64 > 100 &&
		(content_type == "image/jpeg" || content_type == "image/png" || content_type == "image/gif"))
}