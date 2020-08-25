package main

import (
	"fmt"
	"time"

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
	Twitch_username         string
	Twitch_id               int64
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
	Announce_twitch_live     bool `gorm:"default:false"`

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

type subscribeBody struct {
	Mode         string `json:"hub.mode"`
	Topic        string `json:"hub.topic"`
	Callback     string `json:"hub.callback"`
	LeaseSeconds int64  `json:"hub.lease_seconds"`
	Secret       string `json:"hub.secret"`
}

type OAuthTokenBody struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AppAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type userFollowsPayload struct {
	Data []struct {
		FollowedAt time.Time `json:"followed_at"`
		FromID     string    `json:"from_id"`
		FromName   string    `json:"from_name"`
		ToID       string    `json:"to_id"`
		ToName     string    `json:"to_name"`
	} `json:"data"`
}

type streamChangedPayload struct {
	Data []struct {
		ID           string        `json:"id"`
		UserID       string        `json:"user_id"`
		UserName     string        `json:"user_name"`
		GameID       string        `json:"game_id"`
		CommunityIds []interface{} `json:"community_ids"`
		Type         string        `json:"type"`
		Title        string        `json:"title"`
		ViewerCount  int           `json:"viewer_count"`
		StartedAt    time.Time     `json:"started_at"`
		Language     string        `json:"language"`
		ThumbnailURL string        `json:"thumbnail_url"`
	} `json:"data"`
}

var persons []*Person
var files []string
var appAccessTkn AppAccessToken
var streamChangedJSON streamChangedPayload
var jsonPayload JsonPayload
var db *gorm.DB
var twitchNotificationID, streamID string

func getPersonsDB() []*Person {
	// SELECT * FROM people WHERE enabled = true;
	db.Where("announce_twitch_live = ?", "true").Find(&persons)
	/*for _, s := range persons {
		log.Println(s)
	}*/
	return persons
}

func findPersonDB(id string) []*Person {
	// SELECT * FROM people WHERE enabled = true;
	db.Where("twitch_id = ?", id).First(&persons)
	/*for _, s := range persons {
		log.Println(s)
	}*/
	return persons
}

func initDB() {
	var err error
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		POSTGRES_HOST, POSTGRES_USER, POSTGRES_DB, POSTGRES_PASSWORD)
	//log.Printf("%s", dbUri)
	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		reportTg("failed to connect database")
		panic("failed to connect database")
	}
	//defer db.Close()
	db.LogMode(true)

	// Migrate the schema
	db.Debug().AutoMigrate(&Person{})
}
