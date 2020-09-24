package main

import (
//	"fmt"
//	"reflect"

    "github.com/jinzhu/gorm"
	"github.com/wmw9/blossom-reposter/pkg/pubsub"
    "github.com/wmw9/blossom-reposter/pkg/database"
)


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

var files []string
var status Status
var vkPost VKPost
var s pubsub.ServerSocket
var db *gorm.DB
var jsonPayload database.JsonPayload
