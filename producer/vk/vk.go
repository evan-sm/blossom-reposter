package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/k0kubun/pp"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
    "github.com/wmw9/blossom-reposter/pkg/database"
)

// UsersGet used to retrieve status
func UsersGet(id int64, s string, v *database.User) {
	clearJSON() // Wipe it from last Unmarshal 
    jsonPayload = database.ComposeJSONPayload(v, "vk")
	query := fmt.Sprintf(`{"v": "%v", "fields":"%v", "user_id": "%v", "access_token": "%v"}`,
		vkAPIVersion, vkUsersGetFields, id, vkAccessTkn)
	resp, _, errs := gorequest.New().Get(vkUsersGetURL).Query(query).EndStruct(&status) // Get status and unmarshal into struct
	if errs != nil {
		log.Fatalf("%v\n%v", resp, errs)
	}

	if status.Response[0].Status == "" {
		log.Println("Skip. ⏩ Empty status")
		return
	}

	if s == status.Response[0].Status {
		log.Printf("[status] Skip. ⏩ Old status: %v", s)
		return
	}
	log.Printf("🆕 New status: %v", status.Response[0].Status)

	jsonPayload.Caption = status.Response[0].Status
	jsonPayload.Type = "status"
	jsonPayload.Source = fmt.Sprintf("https://vk.com/id%v", jsonPayload.VkPageID)
	pp.Print(jsonPayload)
    if sent := sendJSONPayload(); sent {
		log.Printf("Mark it in DB")
        database.UpdateVKStatusDB(db, id, status.Response[0].Status)
	}
	time.Sleep(5 * time.Second)
}

// WallGet used to retrieve posts
func WallGet(id int64, v *database.User) {
	var pathSyntax string

	clearJSON() // Wipe it from last Unmarshal 
    jsonPayload = database.ComposeJSONPayload(v, "vk")

	// Call VK API wall.get method to retrieve 5 last posts
	query := fmt.Sprintf(`{"v": "%v", "filter":"%v", "owner_id": "%v", "count":"5", "access_token": "%v"}`,
		vkAPIVersion, vkWallGetFilter, id, vkAccessTkn)
	resp, body, errs := gorequest.New().Get(vkWallGetURL).Query(query).End()
	if errs != nil {
		log.Fatalf("%v\n%v\n%v", resp, errs, body)
	}
	// log.Printf("body: %v\n\n", body)

	// Drop old posts we don't need by checking it against "vk_page_timestamp" column from postgres "people" table
	// If there's more than 1 new posts then get the last one from the end [-1] and leave others for the next cycle
	if id > 0 {
		pathSyntax = fmt.Sprintf("response.items.#(date>%v)#|0", jsonPayload.VkPageTimestamp)
	} else if id < 0 {
		pathSyntax = fmt.Sprintf("response.items.#(date>%v)#|0", jsonPayload.VkPublicTimestamp)
	}
	js := gjson.Get(body, pathSyntax)

	// Unmarshal our new post into vkPost struct
	if err := json.Unmarshal([]byte(js.String()), &vkPost); err != nil {
		if id > 0 {
			log.Println("[page] Skip. ⏩ Still quiet. Err:", err)
		} else if id < 0 {
			log.Println("[public] Skip. ⏩ Still quiet. Err:", err)
		}
		time.Sleep(1 * time.Second)
		return
	}

	// Skip cycle if post is a repost.
	if gjson.Get(js.String(), "copy_history").Exists() {
		log.Println("Repost. Skip.")
        database.UpdateVKTimestampDB(db, id, jsonPayload.Person, vkPost.Date)
		//updateVkTimestamp(id)
		return
	}

	// Skip cycle if both js["attachments"] not present and "js["text"]" is empty.
	if txt := gjson.Get(js.String(), "text").String(); txt == "" && !gjson.Get(js.String(), "attachments").Exists() {
		log.Println("Skip. js['attachments'] not present and js['text'] is empty.")
        database.UpdateVKTimestampDB(db, id, jsonPayload.Person, vkPost.Date)
		//updateVkTimestamp(id)
		return
	}

	log.Printf("js: %v\n\n", js)
	// Prepare payload
	for _, v := range vkPost.Attachments {
		if v.Type == "photo" { // Only photo types
			var width int
			var url string
			for _, v := range v.Photo.Sizes { // Iterate through each sizes to get max hi-res jpeg
				if v.Width > width {
					width = v.Width
					url = v.URL
				}
			}
			log.Printf("%vp @ %v", width, url)
			files = append(files, url) // add .jpg to slice
		}
	}
    log.Printf("\nWallGet: %v\n", jsonPayload)
	jsonPayload.Files = files
	jsonPayload.Timestamp = vkPost.Date
	jsonPayload.Caption = vkPost.Text
	jsonPayload.Type = "post"
	jsonPayload.Source = fmt.Sprintf("https://vk.com/wall%v_%v", vkPost.FromID, vkPost.ID)
	pp.Print(vkPost)
	pp.Print(jsonPayload)
    if sent := sendJSONPayload(); sent {
        database.UpdateVKTimestampDB(db, id, jsonPayload.Person, vkPost.Date)
    }
	time.Sleep(5 * time.Second)
}

//func updateVkTimestamp(id int64) {
//	var col string
//
//	if id > 0 {
//		col = "vk_page_timestamp"
//	} else if id < 0 {
//		col = "vk_public_timestamp"
//	}
//	db.Model(&users).Where("person = ?", jsonPayload.Person).
//		Update(col, vkPost.Date)
//
//}
