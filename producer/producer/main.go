package main

import (
	"encoding/json"

	//	"fmt"
	"log"
	//	"os"
	//"reflect"
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

var shortcodeJSONArray []database.ShortcodeJSON
var storyJSON []database.StoryJSON
var shortcodeJSON database.ShortcodeJSON

func main() {
	s = pubsub.NewServerSocket()
	s.ServerInit(socketURL, "VK producer")

	db = database.InitDB()

	go func() {
		for {
			time.Sleep(3 * time.Second)
			checkVK()
			log.Printf("⏳ Next run in 1 m ...")
			time.Sleep(60 * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(20 * time.Second)
			checkIG()
			log.Printf("⏳ Next run in 12 m ...")
			time.Sleep(720 * time.Second)
		}
	}()

	forever := make(chan bool)
	<-forever
}

func checkIG() {
	u := database.GetUsersDB(db)
	for _, v := range u {
		log.Printf("Checking %s's instagram...", v.Person)

		if v.Check_instagram_post {
			checkInstagramPost(v)
		}
		if v.Check_instagram_story {
			checkInstagramStory(v)
		}

	}
}

func checkVK() {
	u := database.GetUsersDB(db)
	for _, v := range u {
		log.Printf("Checking %s's VK...", v.Person)
		if v.Repost_vk_page_enabled {
			if v.Check_vk_page {
				WallGet(v.Vk_page_id, v)
			}
		}
		if v.Repost_vk_public_enabled {
			if v.Check_vk_public {
				WallGet(v.Vk_public_id, v)
			}
		}
		if v.Repost_vk_status_enabled {
			if v.Check_vk_status {
				UsersGet(v.Vk_page_id, v.Vk_status_text, v)
			}
		}
	}
}

func clearJSON() {
	jsonPayload = database.JsonPayload{}
	files = nil
	status = database.Status{}
}

func sendJSONPayload() bool {
	body, err := json.Marshal(&jsonPayload)
	if err != nil {
		die("sendJSONPayload failed", err)
	}
	s.ServerSend(body)
	return true
}
