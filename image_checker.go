package main

import (
	"time"

	"github.com/robfig/cron"
	"github.com/parnurzeal/gorequest"
)

func image_checker(db *DataBase, config *Config) {

	c := cron.New()

	c.AddFunc("@daily", func() {
		check_images(db)
	})

}

func check_images(db *DataBase) {

	images := db.get_all()

	request := gorequest.New().Timeout(time.Second*4)

	total := len(images)
	log.Info("Got %s images", total)

	for i, image := range images {

		log.Info("Progress: %s/%s", (i+1), total)

		resp, _, errs :=  request.Head(image["origurl"].(string)).End();
		if(len(errs) != 0 || resp.StatusCode != 200) {
			db.set_deleted(image["id"]);
		}
	}
}