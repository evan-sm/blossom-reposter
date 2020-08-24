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
	Check_vk_public       bool `gorm:"default:false"`

	Telegram_chan_id   int64
	Instagram_username string
	Instagram_id       int64
	Vk_page_id         int64
	Vk_public_id       int64

	Repost_telegram_chan_id int64
	Repost_telegram_enabled bool `gorm:"default:false"`
	Repost_makaba_enabled   bool `gorm:"default:false"`

	Telegram_chan_id_timestamp int64 `sql:"DEFAULT:extract(epoch from now())"`
	Instagram_story_timestamp  int64 `sql:"DEFAULT:extract(epoch from now())"`
	Instagram_post_timestamp   int64 `sql:"DEFAULT:extract(epoch from now())"`
	Vk_page_timestamp          int64 `sql:"DEFAULT:extract(epoch from now())"`
	Vk_public_timestamp        int64 `sql:"DEFAULT:extract(epoch from now())"`
}

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
