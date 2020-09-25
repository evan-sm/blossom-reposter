package main

import (
	"encoding/json"
//	"fmt"
	"log"
//	"os"
	//"reflect"
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

func main() {
	s = pubsub.NewServerSocket()
	s.ServerInit(socketURL, "VK producer")

    db = database.InitDB()

	for {
		checkSN()
		log.Printf("⏳ Next run in 1m ...")
		time.Sleep(60 * time.Second)
	}
}

func checkSN() {
    users := database.GetUsersDB(db)
	
    // Iterate through each person
	for _, v := range users {
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

