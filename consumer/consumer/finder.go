package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
)

func getCatalog() ([]byte, string, string) {
	person := jsonPayload.Person

	var keyword string
	var board string
	switch person {
	case "germanika":
		keyword = "германик"
		board = "fag"
	case "linyasha":
		keyword = "самых ламповых"
		board = "fag"
	case "mellstroy":
		keyword = "MELLSTROY"
		board = "fag"
	case "sharisha":
		//keyword = "не такой как все"
		keyword = " как все"
		board = "fag"
	case "gabi":
		keyword = "самых ламповых"
		board = "fag"
	case "olyashaa":
		keyword = "ляша"
		board = "fag"
	default:
		keyword = `навальный \ролл /ролл`
		board = "test"
	}
	url := fmt.Sprintf("https://2ch.hk/%v/catalog.json", board)
	resp, err := http.Get(url)
	if err != nil {
		reportTg(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	return body, keyword, board
}

func findThread() (string, string) {
	catalogJson, keyword, board := getCatalog()
	threads := gjson.GetBytes(catalogJson, `threads.#.subject`)
	var ind int // Thread index
	for k, v := range threads.Array() {
		if strings.Contains(strings.ToLower(v.String()), keyword) == true {
			fmt.Println("Thread found; Index is:", k, "; subject is:", v)
			ind = k
		}
	}
	gjsonPath := fmt.Sprintf("threads.%v.num", ind)

	num := gjson.GetBytes(catalogJson, gjsonPath).String()
	//fmt.Println("Thread number is:", num)
	return board, num
}
