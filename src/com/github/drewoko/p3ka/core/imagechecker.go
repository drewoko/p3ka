package core

import (
	"log"
	"time"

	"github.com/parnurzeal/gorequest"
)

func ImageChecker(db *DataBase) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovering ImageChecker", r)
		}
	}()

	startImageChecker(db)
}

func startImageChecker(db *DataBase)  {

	for ;; {
		log.Println("Starting image checking")

		checkImages(db)

		time.Sleep(time.Hour)
	}
}

func checkImages(db *DataBase) {

	images := db.GetAll()

	request := gorequest.New().Timeout(time.Second*4)

	total := len(images)
	log.Printf("Got %d images", total)

	for _, image := range images {
		if(image["url"] != nil) {
			resp, _, _ :=  request.Head(image["url"].(string)).End();
			if(resp != nil && resp.StatusCode != 200) {
				db.SetDeleted(image["id"])
			};
		} else {
			db.SetDeleted(image["id"]);
		}
	}
	log.Println("Image checking finished")
}