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

	//"github.com/wmw9/blossom-reposter/pkg/pubsub"
	"github.com/k0kubun/pp"
	"github.com/wMw9/rpdb"
	//    "github.com/jinzhu/gorm"
)

func checkInstagramPost(v *rpdb.User) {
	clearJSON() // Wipe it from last Unmarshal
	jsonPayload = rpdb.ComposeJSONPayload(v, "ig")
	jsonPayload.Type = "post"
	pp.Println(jsonPayload)
	log.Printf("Checking %v's profile json for posts...", jsonPayload.InstagramUsername)

	js := extractJsonFromProfilePage(jsonPayload.InstagramUsername)
	if js == "" {
		log.Printf("%v IG posts is empty: %v", js, jsonPayload.InstagramUsername)
		reportTg(js)
		return
	}
	//log.Printf("\nfull js:\n%s\n\n", js)

	gjspath := fmt.Sprintf(`entry_data.ProfilePage.0.graphql.user.edge_owner_to_timeline_media.edges.@reverse.#(node.taken_at_timestamp>%v).node.shortcode`, jsonPayload.InstagramPostTimestamp) // 592655783)
	shortcode := gjson.Get(js, gjspath).String()
	if shortcode == "" {
		//log.Printf("shortcode: %v", shortcode)
		log.Printf("Cannot get shortcode, skipping...: %v", shortcode)
		return
	}
	log.Printf("Shortcode: %v Post timestamp: %v\n", shortcode, jsonPayload.InstagramPostTimestamp)
	time.Sleep(5 * time.Second)

	jspage, multi, timestamp, caption := extractJsonFromPostPage(shortcode)
	if jspage == "" {
		log.Printf("%v No new IG posts", jsonPayload.InstagramUsername)
		time.Sleep(2 * time.Second)
		return
	}
	jsonPayload.Source = fmt.Sprintf("https://instagram.com/p/%v", shortcode)
	//	log.Printf("post json:\n\n%v %v", jspage, multi)

	extractFilesFromJson(jspage, multi)

	// Prepare JsonPayload files
	jsonPayload.Files = files
	jsonPayload.Timestamp = timestamp
	jsonPayload.Caption = caption
	pp.Println(jsonPayload)
	if sent := sendJSONPayload(); sent {
		log.Printf("Mark it in DB")
		rpdb.UpdateIGPostTimestampDB(db, jsonPayload.Person, jsonPayload.Timestamp)
	}
	log.Printf("Sleep for a sec")
	time.Sleep(1 * time.Second)

}

func checkInstagramStory(v *rpdb.User) {
	clearJSON() // Wipe it from last Unmarshal
	jsonPayload = rpdb.ComposeJSONPayload(v, "ig")

	jsonPayload.Type = "story"
	//jsonPayload.Caption = ""

	jsonPayload.Source = fmt.Sprintf("https://instagram.com/stories/%v", jsonPayload.InstagramUsername)
	log.Printf("Checking %v's stories... person: %v; timestamp: %v", jsonPayload.InstagramUsername, jsonPayload.Person, jsonPayload.InstagramStoryTimestamp)
	url := fmt.Sprintf(`https://i.instagram.com/api/v1/feed/user/%v/story/`, jsonPayload.InstagramID)
	//log.Printf("url: %v", url)
	req := gorequest.New()
	resp, body, errs := req.Get(url).
		Set("user-agent", "Instagram 10.26.0 (iPhone7,2; iOS 10_1_1; en_US; en-US; scale=2.00; gamut=normal; 750x1334) AppleWebKit/420+").
		Set("cookie", IGSessionID).
		Retry(4, 1200*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusTooManyRequests).End()
	//log.Printf("%v", body)
	//log.Printf("resp: %v", resp.Status)
	if errs != nil {
		log.Printf("%v %v %v", url, resp.Status, errs)
		reportTg(errs)
	}

	jsonPayload.Timestamp = gjson.Get(body, "reel.latest_reel_media").Int()
	log.Printf("Last story timestamp from DB: %v; Post timestamp: %v ", jsonPayload.InstagramStoryTimestamp, jsonPayload.Timestamp)
	if jsonPayload.InstagramStoryTimestamp >= jsonPayload.Timestamp {
		log.Printf("📭 New stories not found.")
		time.Sleep(10 * time.Second)
		return
	}
	result := gjson.Get(body, "reel.items")
	//log.Printf("%v", result.String())
	if err := json.Unmarshal([]byte(result.String()), &storyJSON); err != nil {
		log.Println("err:", err)
	}
	// files = nil // clear previous files
	for _, v := range storyJSON {
		log.Printf("proccessing %v: DB timestamp: %v vs Story timestamp: %v",
			v.Code, jsonPayload.InstagramStoryTimestamp, v.TakenAt)
		log.Printf("Files count: %v", len(files))
		if jsonPayload.InstagramStoryTimestamp < v.TakenAt && len(files) != 4 {
			jsonPayload.Timestamp = v.TakenAt
			if v.StoryCta != nil {
				if v.StoryCta[0].Links[0].WebURI != "" {
					jsonPayload.Caption = ""
					jsonPayload.Caption = fmt.Sprintf("Swipe up ⤴️: %v", v.StoryCta[0].Links[0].WebURI)
					log.Printf("%v", jsonPayload.Caption)
				}
			}
			if v.MediaType == 1 {
				log.Printf("🖼   jpg found \"%v\"; appending to []files",
					v.ImageVersions2.Candidates[0].URL)
				files = append(files, v.ImageVersions2.Candidates[0].URL)
			}
			if v.MediaType == 2 {
				log.Printf("📹 mp4 found \"%v\"; appending to []files",
					v.VideoVersions[0].URL)
				files = append(files, v.VideoVersions[0].URL)
			}
		}
	}
	// Prepare JsonPayload files
	jsonPayload.Files = files
	pp.Println(jsonPayload)
	if sent := sendJSONPayload(); sent {
		log.Printf("Mark it in DB")
		rpdb.UpdateIGStoryTimestampDB(db, jsonPayload.Person, jsonPayload.Timestamp)
	}
	time.Sleep(5 * time.Second)

}

func extractJsonFromProfilePage(username string) string {
	url := fmt.Sprintf(`https://www.instagram.com/%v`, username)
	//log.Printf("url: %v\n", url)
	req := gorequest.New()
	resp, body, errs := req.Get(url).
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36").
		Set("cookie", IGSessionID).
		Retry(4, 1200*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusTooManyRequests).End()
	//log.Printf("resp: %v\n\n", resp.Status)
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
	//log.Printf("url: %v\n", url)
	req := gorequest.New()
	resp, body, errs := req.Get(url).
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36").
		Set("cookie", IGSessionID).
		Retry(4, 1200*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusTooManyRequests).End()
	//log.Printf("resp: %v\n\n", resp.Status)
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
		if err := json.Unmarshal([]byte(r.String()), &shortcodeJSONArray); err != nil {
			log.Println("err:", err)
		}
		files = nil // clear previous files
		for _, v := range shortcodeJSONArray {
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
		if err := json.Unmarshal([]byte(js), &shortcodeJSON); err != nil {
			log.Println("err:", err)
		}
		log.Printf("proccessing %v", shortcodeJSON.Shortcode)
		if shortcodeJSON.IsVideo {
			log.Printf("📹 mp4 found \"%v\"; appending to files slice", shortcodeJSON.VideoURL)
			files = append(files, shortcodeJSON.VideoURL)
		} else {
			log.Printf("🖼   jpg found \"%v\"; appending to files slice", shortcodeJSON.DisplayURL)
			files = append(files, shortcodeJSON.DisplayURL)
		}
	}

}
