package database

import (
	"fmt"
    "reflect"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Person string `gorm:"unique;not null"`

	Enabled               bool `gorm:"default:true"`
	Check_telegram        bool `gorm:"default:false"`
	Check_instagram_post  bool `gorm:"default:true"`
	Check_instagram_story bool `gorm:"default:true"`
	Check_vk_page         bool `gorm:"default:false"`
	Check_vk_status       bool `gorm:"default:false"`
	Check_vk_public       bool `gorm:"default:false"`

	Telegram_chan_id        int64
	Instagram_username      string
	Instagram_id            int64
	Vk_page_id              int64
	Vk_public_id            int64
	Repost_telegram_chan_id int64

	Repost_telegram_enabled  bool `gorm:"default:false"`
	Repost_makaba_enabled    bool `gorm:"default:false"`
	Repost_vk_page_enabled   bool `gorm:"default:false"`
	Repost_vk_public_enabled bool `gorm:"default:false"`
	Repost_vk_status_enabled bool `gorm:"default:false"`

	Telegram_chan_id_timestamp int64 `sql:"DEFAULT:extract(epoch from now())"`
	Instagram_story_timestamp  int64 `sql:"DEFAULT:extract(epoch from now())"`
	Instagram_post_timestamp   int64 `sql:"DEFAULT:extract(epoch from now())"`
	Vk_page_timestamp          int64 `sql:"DEFAULT:extract(epoch from now())"`
	Vk_status_timestamp        int64 `sql:"DEFAULT:extract(epoch from now())"`
	Vk_public_timestamp        int64 `sql:"DEFAULT:extract(epoch from now())"`

	Vk_status_text string
}

//type User struct {
//	gorm.Model
//	Person string `gorm:"unique;not null"; json:"person"`
//
//	Enabled               bool `gorm:"default:true"; json:"enabled"`
//	Check_telegram        bool `gorm:"default:false"; json:"check_telegram"`
//	Check_instagram_post  bool `gorm:"default:true"; json:"check_instagram_post"`
//	Check_instagram_story bool `gorm:"default:true"; json:"check_instagram_story" `
//	Check_vk_page         bool `gorm:"default:false"; json:"check_vk_page" `
//	Check_vk_status       bool `gorm:"default:false"; json:"check_vk_status" `
//	Check_vk_public       bool `gorm:"default:false"; json:"check_vk_public" `
//
//	Telegram_chan_id        int64  `json:"telegram_chan_id"`
//	Instagram_username      string `json:"instagram_username"`
//	Instagram_id            int64  `json:"instagram_id"`
//	Vk_page_id              int64  `json:"vk_page_id"`
//	Vk_public_id            int64  `json:"vk_public_id"`
//	Repost_telegram_chan_id int64  `json:"repost_telegram_chan_id"`
//
//	Repost_telegram_enabled  bool `gorm:"default:false"; json:"repost_telegram_enabled"`
//	Repost_makaba_enabled    bool `gorm:"default:false"; json:"repost_makaba_enabled"`
//	Repost_vk_page_enabled   bool `gorm:"default:false"; json:"repost_vk_page_enabled"`
//	Repost_vk_public_enabled bool `gorm:"default:false"; json:"repost_vk_public_enabled"`
//	Repost_vk_status_enabled bool `gorm:"default:false"; json:"repost_vk_status_enabled"`
//
//	Telegram_chan_id_timestamp int64 `sql:"DEFAULT:extract(epoch from now())"`
//	Instagram_story_timestamp  int64 `sql:"DEFAULT:extract(epoch from now())"`
//	Instagram_post_timestamp   int64 `sql:"DEFAULT:extract(epoch from now())"`
//	Vk_page_timestamp          int64 `sql:"DEFAULT:extract(epoch from now())"`
//	Vk_status_timestamp        int64 `sql:"DEFAULT:extract(epoch from now())"`
//	Vk_public_timestamp        int64 `sql:"DEFAULT:extract(epoch from now())"`
//
//	Vk_status_text string `json:"vk_status_text"`
//}

type JsonPayload struct {
	Timestamp               int64    `json:"timestamp"`
	InstagramPostTimestamp  int64    `json:"instagram_post_timestamp"`
	InstagramStoryTimestamp int64    `json:"instagram_story_timestamp"`
	VkPageTimestamp         int64    `json:"vk_page_timestamp"`
	VkPublicTimestamp       int64    `json:"vk_public_timestamp"`
	VkStatusTimestamp       int64    `json:"vk_status_timestamp"`
	Person                  string   `json:"person"`             //
	InstagramUsername       string   `json:"instagram_username"` //
	InstagramID             int64    `json:"instagram_id"`       //
	Type                    string   `json:"type"`
	From                    string   `json:"from"`
	Source                  string   `json:"source"`
	TelegramChanID          int64    `json:"telegram_chan_id"` //
	RepostMakabaEnabled     bool     `json:"repost_makaba_enabled"`
	RepostTelegramEnabled   bool     `json:"repost_telegram_enabled"`
	RepostVkStatusEnabled   bool     `json:"repost_vk_status_enabled"`
	RepostVkPageEnabled     bool     `json:"repost_vk_page_enabled"`
	RepostVkPublicEnabled   bool     `json:"repost_vk_public_enabled"`
	RepostTelegramChanID    int64    `json:"repost_telegram_chan_id"` //
	VkPageID                int64    `json:"vk_page_id"`              //
	VkPublicID              int64    `json:"vk_public_id"`            //
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

type Status struct {
	Response []struct {
		Status string `json:"status"`
	} `json:"response"`
}

type VKPost struct {
	ID     int `json:"id"`
	FromID int `json:"from_id"`
	// OwnerID     int    `json:"owner_id"`
	Date        int64  `json:"date"`
	Text        string `json:"text"`
	Attachments []struct {
		Type  string `json:"type"`
		Photo struct {
			Sizes []struct {
				Height int    `json:"height"`
				URL    string `json:"url"`
				Type   string `json:"type"`
				Width  int    `json:"width"`
			} `json:"sizes"`
		} `json:"photo,omitempty"`
	} `json:"attachments"`
}

var users []*User
var files []string
var jsonPayload JsonPayload
var status Status
var vkPost VKPost
var db *gorm.DB


func GetUsersDB(db *gorm.DB) []*User {
	// SELECT * FROM people WHERE enabled = true;
	//db.Where("enabled = ?", "true").Find(&persons)
	/*for _, s := range persons {
		log.Println(s)
	}*/

	db.Find(&users, "enabled = true")
	return users
}

func UpdateVKStatusDB(db *gorm.DB, id int64, status string) {
    db.Model(&users).Where("vk_page_id = ?", id).
    Update("vk_status_text", status)
}

func UpdateIGPostTimestampDB(db *gorm.DB, person string, date int64) {
    db.Model(&users).Where("person = ?", person).
    Update("instagram_post_timestamp", date)
}

func UpdateIGStoryTimestampDB(db *gorm.DB, person string, date int64) {
    db.Model(&users).Where("person = ?", person).
    Update("instagram_story_timestamp", date)
}

func UpdateVKTimestampDB(db *gorm.DB, id int64, person string, date int64) {
	var col string

	if id > 0 {
		col = "vk_page_timestamp"
	} else if id < 0 {
		col = "vk_public_timestamp"
	}
	db.Model(&users).Where("person = ?", person).
		Update(col, date)

}

func InitDB() *gorm.DB {
	var err error
	dbUri := fmt.Sprintf("host=%s user=%s port=%s dbname=%s sslmode=disable password=%s",
		POSTGRES_HOST, POSTGRES_USER, POSTGRES_PORT, POSTGRES_DB, POSTGRES_PASSWORD_2)
	fmt.Printf("%s", dbUri)
	db, err = gorm.Open("postgres", dbUri)
	if err != nil {
		panic("failed to connect database")
	}
    fmt.Printf("%v", reflect.TypeOf(db))
	//defer db.Close()
	db.LogMode(true)

	// Migrate the schema
	db.Debug().AutoMigrate(&User{})
    return db
}


func ComposeJSONPayload(s *User, from string) JsonPayload {
	jsonPayload = JsonPayload{}
	files = nil

	jsonPayload.Person = s.Person
	jsonPayload.RepostTelegramEnabled = s.Repost_telegram_enabled
	jsonPayload.RepostTelegramChanID = s.Repost_telegram_chan_id
	jsonPayload.RepostMakabaEnabled = s.Repost_makaba_enabled
	jsonPayload.RepostVkStatusEnabled = s.Repost_vk_status_enabled
	jsonPayload.RepostVkPageEnabled = s.Repost_vk_page_enabled
	jsonPayload.RepostVkPublicEnabled = s.Repost_vk_public_enabled
	jsonPayload.InstagramPostTimestamp = s.Instagram_post_timestamp
	jsonPayload.InstagramStoryTimestamp = s.Instagram_story_timestamp
	jsonPayload.VkPageTimestamp = s.Vk_page_timestamp
	jsonPayload.VkPublicTimestamp = s.Vk_public_timestamp
	jsonPayload.VkStatusTimestamp = s.Vk_status_timestamp
	jsonPayload.InstagramUsername = s.Instagram_username
	jsonPayload.InstagramID = s.Instagram_id
	jsonPayload.VkPageID = s.Vk_page_id
	jsonPayload.DvachBoard = "fag"
	jsonPayload.From = from // "vk"
	//jsonPayload.Files = files

    return jsonPayload
}
