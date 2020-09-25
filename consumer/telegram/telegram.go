package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
    "time"
    "net/http"

	"github.com/k0kubun/pp"
	"github.com/parnurzeal/gorequest"
)

type InputMedia struct {
	Type      string `json:"type"`
	Media     string `json:"media"`
	Caption   string `json:"caption"`
	ParseMode string `json:"parse_mode"`
}

// Query for telegram sendMessage method
type SendMessagePayload struct {
	ChatID                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}


func reportTg(e interface{}) {
	log.Printf("%v", e)
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", tgBotTkn)
	query := fmt.Sprintf(`{"chat_id": %d, "text":"%v"}`, tgChanErr, e)
	gorequest.New().Get(url).Send(query).End()
}

func repostTg() {
	log.Println("Reposting to telegram")
	var caption string

	//log.Printf("Caption length: %v", len(jsonPayload.Caption))
	if len(jsonPayload.Caption) > 1004 {
		log.Printf("%v characters. Caption is long.  I'll cut it!", len(jsonPayload.Caption))
		caption = fmt.Sprintf("%v ... <a href=\"%v\">more</a>", jsonPayload.Caption[:1004], jsonPayload.Source)
	} else {
		caption = jsonPayload.Caption
	}
	log.Printf("Caption length: %v", len(caption))

	switch jsonPayload.From {
	case "vk":
		switch jsonPayload.Type {
		case "status":
			log.Printf("Got VK status: %v\nFrom: %v", caption, jsonPayload.Source)
			caption = fmt.Sprintf("✏️<b><a href=\"%v\">VK статус</a>:</b> %v", jsonPayload.Source, caption)
			log.Printf("Caption: %v", caption)
			sendMessage(caption)
		case "post":
			log.Printf("Git VK post: %v\nFrom: %v", caption, jsonPayload.Source)
			caption = fmt.Sprintf("%v\n\n🔗<a href=\"%v\">VK Post</a>", caption, jsonPayload.Source)
			sendMediaGroup(caption)
		case "public":
			log.Printf("Got VK public post: %v\nFrom: %v", caption, jsonPayload.Source)
			caption = fmt.Sprintf("%v\n\n🔗<a href=\"%v\">VK Public</a>", caption, jsonPayload.Source)
			sendMediaGroup(caption)
		}
	case "ig":
		log.Println("Got post from IG")
		switch jsonPayload.Type {
		case "story":
			if jsonPayload.Caption == "" {
				caption = fmt.Sprintf("🔗<a href=\"%v\">IG Story</a>", jsonPayload.Source)
			} else {
				caption = fmt.Sprintf("%v\n\n🔗<a href=\"%v\">IG Story</a>", caption, jsonPayload.Source)
			}
			sendMediaGroup(caption)
		case "post":
			if jsonPayload.Caption == "" {
				caption = fmt.Sprintf("🔗<a href=\"%v\">IG Post</a>", jsonPayload.Source)
			} else {
				caption = fmt.Sprintf("%v\n\n🔗<a href=\"%v\">IG Post</a>", caption, jsonPayload.Source)
			}
			sendMediaGroup(caption)
		}
	case "twitch":
		log.Println("Got post from twitch")
		switch jsonPayload.Type {
		case "live":
			if jsonPayload.Caption == "" {
				caption = fmt.Sprintf("🔗<a href=\"%v\">Стрим запустился!</a>", jsonPayload.Source)
			} else {
                caption = fmt.Sprintf("Заголовок стрима: \"%v\"\n\n🔗<a href=\"%v\">Стрим запустился!</a>", caption, jsonPayload.Source)
			}
			sendMediaGroup(caption)
		}
	}

}

func sendMessage(s string) {
	tgID := repostTo()
	query := SendMessagePayload{ParseMode: "HTML", DisableWebPagePreview: true}
	query.ChatID = tgID
	query.Text = s
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", tgBotTkn)

	//log.Printf("%v", sjson)
	//log.Printf("%v", string(sjson))
	log.Printf("%v", query)
	resp, body, errs := gorequest.New().Get(url).Query(query).End()
	log.Printf("%v\n%v\n%v", errs, resp, body)

}

func sendMediaGroup(caption string) {
	m := composeInputMedia(caption)
	mjson, _ := json.Marshal(m)
	tgID := repostTo()
	payload := fmt.Sprintf(`{"chat_id": %v, "media": %v}`, tgID, string(mjson))
	url := fmt.Sprintf("https://api.telegram.org/bot%v/sendMediaGroup", tgBotTkn)

	resp, body, errs := gorequest.New().Post(url).Send(payload).
        Retry(3, 5 * time.Second, http.StatusBadRequest, http.StatusInternalServerError).End()
	log.Printf("%v\n%v\n%v", errs, resp, body)
}

func composeInputMedia(caption string) []InputMedia {
	var m []InputMedia

	for _, v := range jsonPayload.Files {
		if strings.Contains(v, ".mp4") {
			if len(m) == 0 {
				m = append(m, InputMedia{"video", v, caption, "HTML"})
				log.Printf("📹 .mp4 found")
			} else {
				m = append(m, InputMedia{"video", v, "", "HTML"})
				log.Printf("📹 .mp4 found")
			}

		}
		if strings.Contains(v, ".jpg") {
			if len(m) == 0 {
				m = append(m, InputMedia{"photo", v, caption, "HTML"})
				log.Printf("🖼 .jpg found")
			} else {
				m = append(m, InputMedia{"photo", v, "", "HTML"})
				log.Printf("🖼 .jpg found")
			}
		}

	}
	pp.Print(m)
	return m
}

func repostTo() int64 {
	var tgID int64
	if jsonPayload.RepostTelegramChanID == 0 {
		tgID = -1001424937175
	} else {
		tgID = jsonPayload.RepostTelegramChanID
	}
	return tgID
}
