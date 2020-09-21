package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/wmw9/blossom-reposter/pkg/pubsub"
)


var s pubsub.ServerSocket

func main() {
	s = pubsub.NewServerSocket()
	fmt.Println(reflect.TypeOf(s))
	s.ServerInit(socketURL, "vk producer")

	initDB()

	for {
		getPersonsDB()
		checkSN()
		log.Printf("⏳ Next run is in 2 minutes...")
		time.Sleep(120 * time.Second)
		//time.Sleep(60 * time.Second)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		reportTg(err)
		log.Printf("%s: %s", msg, err)
	}
}

func checkSN() {
	// Iterate through each person
	for _, v := range persons {
		log.Printf("Checking %s's VK...", v.Person)
		//UsersGet()
		if v.Repost_vk_page_enabled {
			composeJSONPayload(v)
			if v.Check_vk_page == true {
				WallGet(v.Vk_page_id)
			}
		}
		if v.Repost_vk_public_enabled {
			composeJSONPayload(v)
			if v.Check_vk_public == true {
				WallGet(v.Vk_public_id)
			}
		}
		if v.Repost_vk_status_enabled {
			composeJSONPayload(v)
			if v.Check_vk_status == true {
				UsersGet(v.Vk_page_id, v.Vk_status_text)
			}
		}
	}
}

func composeJSONPayload(s *Person) {
	jsonPayload = JsonPayload{}
	files = nil
	jsonPayload.Person = s.Person
	jsonPayload.RepostTelegramEnabled = s.Repost_telegram_enabled
	jsonPayload.RepostTelegramChanID = s.Repost_telegram_chan_id
	jsonPayload.RepostMakabaEnabled = s.Repost_makaba_enabled
	jsonPayload.RepostVkStatusEnabled = s.Repost_vk_status_enabled
	jsonPayload.RepostVkPageEnabled = s.Repost_vk_page_enabled
	jsonPayload.RepostVkPublicEnabled = s.Repost_vk_public_enabled
	jsonPayload.InstagramPostTimestamp = s.Instagram_post_timestamp
	jsonPayload.InstagramStoryTimestamp = s.Instagram_story_timestamp
	jsonPayload.VkPageTimestamp = s.Vk_page_timestamp
	jsonPayload.VkPublicTimestamp = s.Vk_public_timestamp
	jsonPayload.VkStatusTimestamp = s.Vk_status_timestamp
	jsonPayload.InstagramUsername = s.Instagram_username
	jsonPayload.InstagramID = s.Instagram_id
	jsonPayload.VkPageID = s.Vk_page_id
	jsonPayload.DvachBoard = "fag"
	jsonPayload.From = "vk"
	jsonPayload.Files = files
}

func sendJSONPayload() bool {
	body, _ := json.Marshal(&jsonPayload)
	s.ServerSend(body)

	log.Printf(" [x] Sent via tcp socket")

	return true
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}
