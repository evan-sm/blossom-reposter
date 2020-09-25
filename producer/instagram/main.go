package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/wmw9/blossom-reposter/pkg/pubsub"
    "github.com/wmw9/blossom-reposter/pkg/database"
    "github.com/jinzhu/gorm"
)

var files []string
var status database.Status
var vkPost database.VKPost
var s pubsub.ServerSocket
var db *gorm.DB
var jsonPayload database.JsonPayload


var user []*database.User
var shortcodeJsonArray []database.ShortcodeJson
var storyJson []database.StoryJson
var shortcodeJson database.ShortcodeJson

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
			checkInstagramPost(v)
		}
		if v.Check_instagram_story {
			checkInstagramStory(v)
		}
		
    }
}

func clearJSON() {
	jsonPayload = database.JsonPayload{}
	files = nil
	storyJson = nil
}


func sendJSONPayload() bool {
	body, err := json.Marshal(&jsonPayload)
	if err != nil {
		die("sendJSONPayload failed", err)
	}
	s.ServerSend(body)
	return true
}

