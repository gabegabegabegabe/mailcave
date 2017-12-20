package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

func processTarGZ(tarGZFileName string, addr string) {

	tgzf, err := os.Open(tarGZFileName)
	if err != nil {
		panic(err)
	}

	gzf, err := gzip.NewReader(tgzf)
	if err != nil {
		panic(err)
	}

	tr := tar.NewReader(gzf)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(tr)
			if err != nil {
				panic(err)
			}
			mimeStr := buf.String()
			postMessage(mimeStr, addr)
		default:
			panic(fmt.Errorf("unknown file type (%v) in tar file", header.Typeflag))
		}
	}
}

func main() {

	mode := flag.String("mode", "file", "mode: file (single file) or tar.gz (tarball of files)")
	mimeFile := flag.String("mimeFile", "mime.msg", "the name of the file containing a MIME message")
	tarGZFile := flag.String("tarGZFile", "emails.tar.gz", "a tar.gz file containing MIME messages")
	addr := flag.String("mailcaveAddr", "http://localhost:8080/store", "the address of the POST endpoint")
	flag.Parse()

	if *mode == "file" {
		processSingleFile(*mimeFile, *addr)
	} else {
		processTarGZ(*tarGZFile, *addr)
	}
}
