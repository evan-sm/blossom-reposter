package main

type JsonPayload struct {
	Timestamp      int64    `json:"timestamp"`
	Person         string   `json:"person"`
	Type           string   `json:"type"`
	From           string   `json:"from"`
	Source         string   `json:"source"`
	RepostTelegram bool     `json:"repost_telegram"`
	TelegramChanID int64    `json:"telegram_chan_id"`
	Repost2Ch      bool     `json:"repost_2ch"`
	DvachBoard     string   `json:"2ch_board"`
	Files          []string `json:"files"`
	Caption        string   `json:"caption"`
}

type ShortcodeJson struct {
	ID           string `json:"id"`
	Shortcode    string `json:"shortcode"`
	MediaPreview string `json:"media_preview"`
	DisplayURL   string `json:"display_url"`
	HasAudio     bool   `json:"has_audio,omitempty"`
	VideoURL     string `json:"video_url,omitempty"`
	IsVideo      bool   `json:"is_video"`
}

type StoryJson struct {
	TakenAt        int64  `json:"taken_at"`
	MediaType      int64  `json:"media_type"`
	Code           string `json:"code"`
	ImageVersions2 struct {
		Candidates []struct {
			Width        int    `json:"width"`
			Height       int    `json:"height"`
			URL          string `json:"url"`
			ScansProfile string `json:"scans_profile"`
		} `json:"candidates"`
	} `json:"image_versions2"`
	OriginalWidth  int         `json:"original_width"`
	OriginalHeight int         `json:"original_height"`
	Caption        interface{} `json:"caption"`
	VideoVersions  []struct {
		Type   int    `json:"type"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		URL    string `json:"url"`
		ID     string `json:"id"`
	} `json:"video_versions,omitempty"`
	StoryCta []struct {
		Links []struct {
			WebURI string `json:"webUri"`
		} `json:"links"`
	} `json:"story_cta"`
}

var persons []*Person
var files []string
var shortcodeJsonArray []ShortcodeJson
var storyJson []StoryJson
var shortcodeJson ShortcodeJson
var jsonPayload JsonPayload
