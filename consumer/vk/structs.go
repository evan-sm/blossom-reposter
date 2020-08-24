package main

type JsonPayload struct {
	Timestamp             int64    `json:"timestamp"`
	Person                string   `json:"person"`
	Type                  string   `json:"type"`
	From                  string   `json:"from"`
	Source                string   `json:"source"`
	RepostTelegram        bool     `json:"repost_telegram"`
	TelegramChanID        int64    `json:"telegram_chan_id"`
	RepostMakabaEnabled   bool     `json:"repost_makaba_enabled"`
	RepostTelegramEnabled bool     `json:"repost_telegram_enabled"`
	RepostTelegramChanID  int64    `json:"repost_telegram_chan_id"`
	DvachBoard            string   `json:"2ch_board"`
	Files                 []string `json:"files"`
	Caption               string   `json:"caption"`
}

var jsonPayload JsonPayload
