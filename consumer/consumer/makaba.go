package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

func repost2ch() bool {
	log.Println("Post is ready and will be sent in 20 secs.")
	time.Sleep(20 * time.Second)
	board, thread := findThread()
	log.Printf("https://2ch.hk/%v/res/%v.html", board, thread)
	valuesBase := prepareBase(board, thread)
	valuesFiles := prepareFiles()

	client, ok := customClient()
	if ok == false {
		return false
	}
	//fmt.Println("valuesFiles type is:", reflect.TypeOf(valuesFiles))
	err, success := makabaPost(client, postingUrl, valuesBase, valuesFiles)
	if err != nil {
		log.Println(err)
		reportTg(err)
	}
	return success
}

func makabaPost(client *http.Client, url string, valuesBase map[string]io.Reader, valuesFiles map[string]io.Reader) (err error, success bool) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range valuesBase {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if fw, err = w.CreateFormField(key); err != nil {
			return
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err, false
		}

	}
	for key, r := range valuesFiles {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if fw, err = w.CreateFormFile(key, ""); err != nil {
			return
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err, false
		}

	}
	w.Close()

	// Prepare handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Высрать в тред
	res, err := client.Do(req)
	if err != nil {
		log.Println("client.Do(req) error:", err)
		reportTg(err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error:", err)
		reportTg(err)
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if result["Error"] != nil {
		log.Println("Makaba post error:", result)
		reportTg(result)
	}
	log.Println(result)
	if result["Error"] == nil {
		log.Println("Successfully made post 👌🏻")
		success = true
	}
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		reportTg(err)
	}
	return err, success
}
