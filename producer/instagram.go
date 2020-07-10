package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/parnurzeal/gorequest"
	"github.com/tidwall/gjson"
)

func clearJSONPayload() {
	jsonPayload.Timestamp = 0
	jsonPayload.Person = ""
	jsonPayload.Type = ""
	jsonPayload.Source = ""
	jsonPayload.RepostTelegram = false
	jsonPayload.Repost2Ch = false
	jsonPayload.DvachBoard = ""
	jsonPayload.Files = jsonPayload.Files[:0]
	files = files[:0]
	jsonPayload.Caption = ""
}

func extractJsonFromProfilePage(username string) string {
	url := fmt.Sprintf(`https://www.instagram.com/%v`, username)
	log.Printf("url: %v\n", url)
	req := gorequest.New()
	resp, body, errs := req.Get(url).
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36").
		Set("cookie", IGSessionID).
		Retry(4, 1200*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusTooManyRequests).End()
	log.Printf("resp: %v\n\n", resp.Status)
	if errs != nil {
		log.Fatalf("%v\n%v", errs, resp)
	}
	var js string
	doc := soup.HTMLParse(body)
	links := doc.FindAll("script", "type", "text/javascript")
	for _, link := range links {
		if strings.Contains(link.Text(), "window._sharedData = ") {
			js = strings.Replace(link.Text(), "window._sharedData = ", "", 1)
			//log.Printf("js:\n%s", js)
			js = strings.TrimSuffix(js, ";")
			return js
		}
	}
	return js
}

func extractJsonFromPostPage(shortcode string) (string, bool, int64, string) {
	url := fmt.Sprintf("https://www.instagram.com/p/%v", shortcode)
	log.Printf("url: %v\n", url)
	req := gorequest.New()
	resp, body, errs := req.Get(url).
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36").
		Set("cookie", IGSessionID).
		Retry(4, 1200*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusTooManyRequests).End()
	log.Printf("resp: %v\n\n", resp.Status)
	if errs != nil {
		log.Fatalf("%v\n%v", errs, resp)
	}

	var js string
	shrink := fmt.Sprintf("window.__additionalDataLoaded('/p/%v/',", shortcode)
	doc := soup.HTMLParse(body)
	links := doc.FindAll("script", "type", "text/javascript")
	for _, link := range links {
		if strings.Contains(link.Text(), shrink) {
			js = strings.Replace(link.Text(), shrink, "", 1)
			js = strings.TrimSuffix(js, ");")
			//log.Printf("post js:\n%s", js)
		}
	}
	js = gjson.Get(js, "graphql.shortcode_media").String()
	multi := gjson.Get(js, "edge_sidecar_to_children").Exists()
	timestamp := gjson.Get(js, "taken_at_timestamp").Int()
	caption := gjson.Get(js, "edge_media_to_caption.edges.0.node.text").String()
	log.Printf("timestamp: %v\nmulti: %v\n", timestamp, multi)
	return js, multi, timestamp, caption
}

func extractFilesFromJson(js string, multi bool) {
	if multi {
		log.Println("Found IG post with multiple objects")
		r := gjson.Get(js, "edge_sidecar_to_children.edges.#.node")
		if err := json.Unmarshal([]byte(r.String()), &shortcodeJsonArray); err != nil {
			log.Println("err:", err)
		}
		files = nil // clear previous files
		for _, v := range shortcodeJsonArray {
			log.Printf("proccessing %v", v.Shortcode)
			if v.IsVideo {
				log.Printf("📹 mp4 found \"%v\"; appending to files slice", v.VideoURL)
				files = append(files, v.VideoURL)
			} else {
				log.Printf("🖼   jpg found \"%v\"; appending to files slice", v.DisplayURL)
				files = append(files, v.DisplayURL)
			}
		}
	} else {
		log.Println("Found IG post with one object")
		if err := json.Unmarshal([]byte(js), &shortcodeJson); err != nil {
			log.Println("err:", err)
		}
		log.Printf("proccessing %v", shortcodeJson.Shortcode)
		if shortcodeJson.IsVideo {
			log.Printf("📹 mp4 found \"%v\"; appending to files slice", shortcodeJson.VideoURL)
			files = append(files, shortcodeJson.VideoURL)
		} else {
			log.Printf("🖼   jpg found \"%v\"; appending to files slice", shortcodeJson.DisplayURL)
			files = append(files, shortcodeJson.DisplayURL)
		}
	}

}

func checkInstagramPost(person string, username string, t int64) {
	clearJSONPayload()

	jsonPayload.Person = person
	jsonPayload.Type = "post"
	jsonPayload.RepostTelegram = false
	jsonPayload.Repost2Ch = true
	jsonPayload.DvachBoard = "fag"

	log.Printf("Checking %v's profile json for posts...", username)

	js := extractJsonFromProfilePage(username)
	if js == "" {
		log.Printf("js is empty, something went wrong: %v", js)
		reportTg(js)
		return
	}
	//log.Printf("\nfull js:\n%s\n\n", js)

	gjspath := fmt.Sprintf(`entry_data.ProfilePage.0.graphql.user.edge_owner_to_timeline_media.edges.@reverse.#(node.taken_at_timestamp>%v).node.shortcode`, t) // 592655783)
	shortcode := gjson.Get(js, gjspath).String()
	if shortcode == "" {
		log.Printf("shortcode: %v", shortcode)
		log.Printf("Couldn't get shortcode, skipping...: %v", shortcode)
		return
	}
	log.Printf("shortcode: %v", shortcode)
	log.Printf("t: %v", t)
	log.Printf("Sleep for 10 sec\n")
	time.Sleep(10 * time.Second)

	jspage, multi, timestamp, caption := extractJsonFromPostPage(shortcode)
	if jspage == "" {
		log.Printf("%v new posts not found", username)
		time.Sleep(2 * time.Second)
		return
	}
	jsonPayload.Source = fmt.Sprintf("https://instagram.com/p/%v", shortcode)
	log.Printf("post json:\n\n%v %v", jspage, multi)

	extractFilesFromJson(jspage, multi)

	// Prepare JsonPayload files
	jsonPayload.Files = files
	jsonPayload.Timestamp = timestamp
	jsonPayload.Caption = caption
	log.Println(jsonPayload)
	sent := sendJsonPayload()
	if sent {
		log.Printf("JsonPayload sent to RabbitMQ, updating DB last post timestamp")
		db.Model(&persons).Where("person = ?", jsonPayload.Person).Update("instagram_post_timestamp", jsonPayload.Timestamp)
	}
	log.Printf("Sleep for a few secs for rate-limits")
	time.Sleep(1 * time.Second)

}

func checkInstagramStory(person string, username string, id int64, t int64) {
	clearJSONPayload()

	jsonPayload.Person = person
	jsonPayload.Type = "story"
	jsonPayload.RepostTelegram = false
	jsonPayload.Repost2Ch = true
	jsonPayload.DvachBoard = "fag"
	jsonPayload.Source = fmt.Sprintf("https://instagram.com/stories/%v", username)
	log.Printf("Checking %v's stories... person: %v; timestamp: %v", username, person, t)
	url := fmt.Sprintf(`https://i.instagram.com/api/v1/feed/user/%v/story/`, id)
	log.Printf("url: %v", url)
	req := gorequest.New()
	resp, body, errs := req.Get(url).
		Set("user-agent", "Instagram 10.26.0 (iPhone7,2; iOS 10_1_1; en_US; en-US; scale=2.00; gamut=normal; 750x1334) AppleWebKit/420+").
		Set("cookie", IGSessionID).
		Retry(4, 1200*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusTooManyRequests).End()
	//log.Printf("%v", body)
	log.Printf("resp: %v", resp.Status)
	if errs != nil {
		log.Printf("%v %v %v", url, resp.Status, errs)
		reportTg(errs)
	}

	jsonPayload.Timestamp = gjson.Get(body, "reel.latest_reel_media").Int()
	log.Printf("Latest story timestamp: %v", jsonPayload.Timestamp)
	log.Printf("Last story timestamp from DB: %v; Post timestamp: %v ", t, jsonPayload.Timestamp)
	if t >= jsonPayload.Timestamp {
		log.Printf("📭 New stories not found. Sleep for a few secs for rate-limits. ")
		time.Sleep(10 * time.Second)
		return
	}
	result := gjson.Get(body, "reel.items")
	//log.Printf("%v", result.String())
	if err := json.Unmarshal([]byte(result.String()), &storyJson); err != nil {
		log.Println("err:", err)
	}
	files = nil // clear previous files
	for _, v := range storyJson {
		log.Printf("proccessing %v: DB timestamp: %v vs Story timestamp: %v",
			v.Code, t, v.TakenAt)
		log.Printf("Files count: %v", len(files))
		if t < v.TakenAt && len(files) != 4 {
			jsonPayload.Timestamp = v.TakenAt
			if v.StoryCta != nil {
				if v.StoryCta[0].Links[0].WebURI != "" {
					jsonPayload.Caption = ""
					jsonPayload.Caption = fmt.Sprintf("Swipe up ⤴️: %v", v.StoryCta[0].Links[0].WebURI)
					log.Printf("%v", jsonPayload.Caption)
				}
			}
			if v.MediaType == 1 {
				log.Printf("🖼   jpg found \"%v\"; appending to files slice",
					v.ImageVersions2.Candidates[0].URL)
				files = append(files, v.ImageVersions2.Candidates[0].URL)
			}
			if v.MediaType == 2 {
				log.Printf("📹 mp4 found \"%v\"; appending to files slice",
					v.VideoVersions[0].URL)
				files = append(files, v.VideoVersions[0].URL)
			}
		}
	}
	// Prepare JsonPayload files
	jsonPayload.Files = files
	log.Println(jsonPayload)
	sent := sendJsonPayload()
	if sent {
		log.Printf("JsonPayload sent to RabbitMQ, updating DB last post timestamp")
		db.Model(&persons).Where("person = ?", jsonPayload.Person).
			Update("instagram_story_timestamp", jsonPayload.Timestamp)
	}
	log.Printf("Sleep for a few secs for rate-limits")
	time.Sleep(10 * time.Second)

}
