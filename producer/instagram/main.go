package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/wmw9/blossom-reposter/pkg/database"
	"github.com/wmw9/blossom-reposter/pkg/pubsub"
)

var files []string
var status database.Status
var vkPost database.VKPost
var s pubsub.ServerSocket
var db *gorm.DB
var jsonPayload database.JsonPayload
var user []*database.User
var shortcodeJSONArray []database.ShortcodeJSON
var storyJSON []database.StoryJSON
var shortcodeJSON database.ShortcodeJSON

func main() {
	s = pubsub.NewServerSocket()
	s.ServerInit(socketURL, "IG producer")

	db = database.InitDB()

	for {
		checkSN()
		log.Printf("⏳ Next run in 1m ...")
		time.Sleep(720 * time.Second)
	}
}

func checkSN() {
	users := database.GetUsersDB(db)

	// Iterate through each person
	for _, v := range users {
		log.Printf("Checking %s's instagram...", v.Person)

		if v.Check_instagram_post {
			checkInstagramPost(v)7
		}
		if v.Check_instagram_story {
			checkInstagramStory(v)
		}

	}
}

func clearJSON() {
	jsonPayload = database.JsonPayload{}
	files = nil
	storyJSON = nil
}

func sendJSONPayload() bool {
	body, err := json.Marshal(&jsonPayload)
	if err != nil {
		die("sendJSONPayload failed", err)
	}
	s.ServerSend(body)
	return true
}
