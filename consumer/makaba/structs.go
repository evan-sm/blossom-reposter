package main

type JsonPayload struct {
	Timestamp               int64    `json:"timestamp"`
	InstagramPostTimestamp  int64    `json:"instagram_post_timestamp"`
	InstagramStoryTimestamp int64    `json:"instagram_story_timestamp"`
	VkPageTimestamp         int64    `json:"vk_page_timestamp"`
	VkPublicTimestamp       int64    `json:"vk_public_timestamp"`
	VkStatusTimestamp       int64    `json:"vk_status_timestamp"`
	Person                  string   `json:"person"`
	TwitchUsername          string   `json:"twitch_username"`
	TwitchID                int64    `json:"twitch_id"`
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
	AnnounceTwitchLive      bool     `json:"announce_twitch_live"`
	RepostTelegramChanID    int64    `json:"repost_telegram_chan_id"`
	VkPageID                int64    `json:"vk_page_id"`
	VkPublicID              int64    `json:"vk_public_id"`
	DvachBoard              string   `json:"2ch_board"`
	Files                   []string `json:"files"`
	Language                string   `json:"language"`
	Caption                 string   `json:"caption"`
}

type Passcode struct {
	Usercode string
	Error    bool
}

var CurrentUsercode Passcode = Passcode{
	Usercode: "",
	Error:    false,
}

var jsonPayload JsonPayload
