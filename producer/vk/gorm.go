package main

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
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

var persons []*Person
var files []string
var jsonPayload JsonPayload
var status Status
var vkPost VKPost
var db *gorm.DB

func convStrInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func getPersonsDB() []*Person {
	// SELECT * FROM people WHERE enabled = true;
	db.Where("enabled = ?", "true").Find(&persons)
	/*for _, s := range persons {
		log.Println(s)
	}*/
	return persons
}

func initDB() {
	var err error
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		POSTGRES_HOST, POSTGRES_USER, POSTGRES_DB, POSTGRES_PASSWORD)
	//log.Printf("%s", dbUri)
	db, err = gorm.Open("postgres", dbUri)
	if err != nil {
		reportTg("failed to connect database")
		panic("failed to connect database")
	}
	//defer db.Close()
	db.LogMode(true)

	// Migrate the schema
	db.Debug().AutoMigrate(&Person{})
}
