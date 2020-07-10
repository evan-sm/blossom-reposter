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

type Passcode struct {
	Usercode string
	Error    bool
}

var CurrentUsercode Passcode = Passcode{
	Usercode: "",
	Error:    false,
}

var jsonPayload JsonPayload
