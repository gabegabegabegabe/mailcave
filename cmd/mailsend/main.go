package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func postMessage(mimeStr string, addr string) error {
	form := url.Values{}
	form.Add("mime_message", mimeStr)

	req, err := http.NewRequest("POST", addr, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	hc := http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("response:\n%s\n", string(bodyBytes))

	return nil
}

func main() {

	mimeFile := flag.String("mimeFile", "mime.msg", "the name of the file containing a MIME message")
	addr := flag.String("mailcaveAddr", "http://localhost:8080/store", "the address of the POST endpoint")
	flag.Parse()

	mimeBytes, err := ioutil.ReadFile(*mimeFile)
	if err != nil {
		panic(err)
	}
	mimeStr := string(mimeBytes)

	err = postMessage(mimeStr, *addr)
	if err != nil {
		panic(err)
	}
}
