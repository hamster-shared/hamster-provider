package utils

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func Download(fullURLFile string, destPath string) error {

	dir := filepath.Dir(destPath)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	// Create blank file
	file, err := os.Create(destPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 30 * time.Second,
	}
	client := http.Client{
		Timeout:   time.Second * 30, // Maximum of 10 secs
		Transport: netTransport,
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d", destPath, size)
	return nil
}
