package main

type JsonPayload struct {
	Timestamp               int64    `json:"timestamp"`
	InstagramPostTimestamp  int64    `json:"instagram_post_timestamp"`
	InstagramStoryTimestamp int64    `json:"instagram_story_timestamp"`
	VkPageTimestamp         int64    `json:"vk_page_timestamp"`
	VkPublicTimestamp       int64    `json:"vk_public_timestamp"`
	VkStatusTimestamp       int64    `json:"vk_status_timestamp"`
	Person                  string   `json:"person"`
	InstagramUsername       string   `json:"instagram_username"`
	InstagramID             int64    `json:"instagram_id"`
	Type                    string   `json:"type"`
	From                    string   `json:"from"`
	Source                  string   `json:"source"`
	TelegramChanID          int64    `json:"telegram_chan_id"`
	RepostMakabaEnabled     bool     `json:"repost_makaba_enabled"`
	RepostTelegramEnabled   bool     `json:"repost_telegram_enabled"`
	RepostVkStatusEnabled   bool     `json:"repost_vk_status_enabled"`
	RepostVkPageEnabled     bool     `json:"repost_vk_page_enabled"`
	RepostVkPublicEnabled   bool     `json:"repost_vk_public_enabled"`
	RepostTelegramChanID    int64    `json:"repost_telegram_chan_id"`
	VkPageID                int64    `json:"vk_page_id"`
	VkPublicID              int64    `json:"vk_public_id"`
	DvachBoard              string   `json:"2ch_board"`
	Files                   []string `json:"files"`
	Caption                 string   `json:"caption"`
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
