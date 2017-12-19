package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

func processSingleFile(mimeFile string, addr string) {

	mimeBytes, err := ioutil.ReadFile(mimeFile)
	if err != nil {
		panic(err)
	}
	mimeStr := string(mimeBytes)

	err = postMessage(mimeStr, addr)
	if err != nil {
		panic(err)
	}
}

func processTarball(tarFileName string, addr string) {

	tarBytes, err := ioutil.ReadFile(tarFileName)
	if err != nil {
		panic(err)
	}

	// Open the tar archive for reading.
	// take from somewhere
	r := bytes.NewReader(tarBytes)
	tr := tar.NewReader(r)

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Contents of %s:\n", hdr.Name)

		buf := new(bytes.Buffer)
		buf.ReadFrom(tr)
		mimeStr := buf.String()
		processSingleFile(mimeStr, addr)
	}
}

func main() {

	mode := flag.String("mode", "file", "mode: file (single file) or tar (tarball of files)")
	mimeFile := flag.String("mimeFile", "mime.msg", "the name of the file containing a MIME message")
	tarFile := flag.String("tarFile", "emails.tar.gz", "a tarball containing MIME messages")
	addr := flag.String("mailcaveAddr", "http://localhost:8080/store", "the address of the POST endpoint")
	flag.Parse()

	if *mode == "file" {
		processSingleFile(*mimeFile, *addr)
	} else {
		processTarball(*tarFile, *addr)
	}
}
