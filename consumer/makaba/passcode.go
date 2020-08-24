package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// PasscodeAuth is used to authorize your passcode to get usercode. Used to bypass captcha
func (c *Passcode) PasscodeAuth() bool {
	formData := url.Values{
		"json":     {"1"},
		"task":     {"auth"},
		"usercode": {passcode}}
	resp, err := http.PostForm(makabaUrl, formData)
	if err != nil {
		reportTg(err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		reportTg(err)
		return false
	}
	//log.Println(string(body))

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		reportTg(err)
		return false
	}
	//log.Println(result)
	if result["result"].(float64) == 0 {
		reportTg(result["description"])
		return false
	}
	if result["result"].(float64) == 1 {
		hash := fmt.Sprint(result["hash"])
		log.Println("✅ Got passcode_auth:", result["hash"])
		c.Usercode = hash
		c.Error = false
		return true
	}

	return false
}
