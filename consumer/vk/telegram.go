package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"strings"
)

type InputMedia struct {
	Type    string `json:"type"`
	Media   string `json:"media"`
	Caption string `json:"caption"`
}

func reportTg(e interface{}) {
	log.Printf("%v", e)
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", tgBotTkn)
	query := fmt.Sprintf(`{"chat_id": %d, "text":"%v"}`, tgChanErr, e)
	gorequest.New().Get(url).Send(query).End()
}

func repostTg() {
	log.Println("Reposting to telegram")
	var m []InputMedia

	for _, v := range jsonPayload.Files {
		if strings.Contains(v, ".mp4") {
			m = append(m, InputMedia{"video", v, jsonPayload.Source})
			log.Printf("📹 .mp4 found")
		}
		if strings.Contains(v, ".jpg") {
			m = append(m, InputMedia{"photo", v, jsonPayload.Source})
			log.Printf("🖼 .jpg found")
		}
		//log.Println(v)
	}

	mjson, _ := json.Marshal(m)

	payload := fmt.Sprintf(`{"chat_id": -1001424937175, "media": %v}`, string(mjson))
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMediaGroup", tgBotTkn)

	gorequest.New().Post(url).
		Send(payload).
		End()
	//log.Printf("%v\n%v\n%v", errs, resp, body)
}
