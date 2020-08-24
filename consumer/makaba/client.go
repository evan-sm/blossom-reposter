package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func customClient() (*http.Client, bool) {
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	auth := CurrentUsercode.PasscodeAuth()
	if auth == false {
		log.Println("Failed to authorize passcode. Skip.")
		reportTg("Failed to authorize passcode. Skip.")
	}
	cookie := &http.Cookie{
		Name:   "passcode_auth",
		Value:  CurrentUsercode.Usercode,
		Path:   "/",
		Domain: "2ch.hk",
	}
	cookies = append(cookies, cookie)
	u, _ := url.Parse(postingUrl)
	jar.SetCookies(u, cookies)
	//log.Println(jar.Cookies(u))
	client := &http.Client{
		Jar: jar,
	}
	return client, auth
}
