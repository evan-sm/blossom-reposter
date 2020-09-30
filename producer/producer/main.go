package main

import (
	"encoding/json"

	//	"fmt"
	"log"
	//	"os"
	//"reflect"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/wMw9/rpdb"
	"github.com/wMw9/rpps"
)

var files []string
var status rpdb.Status
var vkPost rpdb.VKPost
var s rpps.ServerSocket
var db *gorm.DB
var jsonPayload rpdb.JsonPayload

var shortcodeJSONArray []rpdb.ShortcodeJSON
var storyJSON []rpdb.StoryJSON
var shortcodeJSON rpdb.ShortcodeJSON

func main() {
	s = rpps.NewServerSocket()
	s.ServerInit(socketURL, "VK producer")

	db = rpdb.InitDB()

	go func() {
		for {
			time.Sleep(20 * time.Second)

			checkVK()
			checkIG()
			log.Printf("⏳ Next run is VK in 3 m ...\n")
			time.Sleep(180 * time.Second)
			checkVK()
			log.Printf("⏳ Next run is VK in 3 m ...\n")
			time.Sleep(180 * time.Second)
			checkVK()
			log.Printf("⏳ Next run is VK in 3 m ...\n")
			time.Sleep(180 * time.Second)
			checkVK()
			log.Printf("⏳ Next run is VK in 3 m ...\n")
			time.Sleep(180 * time.Second)
			checkVK()
			log.Printf("⏳ Next run is IG in 3 m ...\n")

			time.Sleep(20 * time.Second)
		}
	}()

	forever := make(chan bool)
	<-forever
}

func checkIG() {
	u := rpdb.GetUsersDB(db)
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
	u := rpdb.GetUsersDB(db)
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
	jsonPayload = rpdb.JsonPayload{}
	files = nil
	status = rpdb.Status{}
}

func sendJSONPayload() bool {
	body, err := json.Marshal(&jsonPayload)
	if err != nil {
		die("sendJSONPayload failed", err)
	}
	s.ServerSend(body)
	return true
}
