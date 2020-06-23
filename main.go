package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
	"strings"
)

// Get create a get request
func Get(query, url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// req.Header.Add("User-Agent", "GoodBoy")
	// req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")

	req.URL.RawQuery = query
	trt := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: trt}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getLinkFBVideo(videoURL string) {
	q := url.Values{}
	byteBody, err := Get(q.Encode(), videoURL)
	if err != nil {
		fmt.Println("Error when do request: ", err)
		return
	}

	// regular expression pattern
	r := regexp.MustCompile(`(?P<hdSrc>hd_src):"(?P<hdLink>.*?)",(?P<sdSrc>sd_src):"(?P<sdLink>.*?)"`)
	if strings.Contains(string(byteBody), "hd_src:null") {
		r = regexp.MustCompile(`(?P<hdSrc>hd_src):(?P<hdLink>.*?),(?P<sdSrc>sd_src):"(?P<sdLink>.*?)"`)
	}
	// get capturing group functionality
	names := r.SubexpNames()
	// Use FindAllStringSubmatch for many results
	result := r.FindStringSubmatch(string(byteBody))
	m := map[string]string{}

	if len(result) == 0 {
		fmt.Println("HD Link:", "not found")
		fmt.Println("SD Link:", "not found")
		return
	}

	for i, n := range result {
			m[names[i]] = n
	}
	
	fmt.Println("HD Link:", m["hdLink"])
	fmt.Println("SD Link:", m["sdLink"])
}

func printGopher() {
	// Embed files into a Go executable
	gopherByte := []byte{32,32,32,32,32,32,32,32,32,44,95,45,45,45,126,126,126,126,126,45,45,45,45,46,95,32,32,32,32,32,32,32,32,32,13,10,32,32,95,44,44,95,44,42,94,95,95,95,95,32,32,32,32,32,32,95,95,95,95,95,96,96,42,103,42,92,34,42,44,32,32,32,32,32,32,32,32,65,117,116,104,111,114,58,32,86,73,73,70,84,87,13,10,32,47,32,95,95,47,32,47,39,32,32,32,32,32,94,46,32,32,47,32,32,32,32,32,32,92,32,94,64,113,32,32,32,102,32,32,32,32,32,32,32,69,90,32,71,101,116,32,76,105,110,107,32,68,111,119,110,108,111,97,100,32,86,105,100,101,111,32,70,97,99,101,98,111,111,107,13,10,91,32,32,64,102,32,124,32,64,41,41,32,32,32,32,124,32,32,124,32,64,41,41,32,32,32,108,32,32,48,32,95,47,32,32,13,10,32,92,96,47,32,32,32,92,126,95,95,95,95,32,47,32,95,95,32,92,95,95,95,95,95,47,32,32,32,32,92,32,32,32,13,10,32,32,124,32,32,32,32,32,32,32,32,32,32,32,95,108,95,95,108,95,32,32,32,32,32,32,32,32,32,32,32,73,32,32,32,13,10,32,32,125,32,32,32,32,32,32,32,32,32,32,91,95,95,95,95,95,95,93,32,32,32,32,32,32,32,32,32,32,32,73,32,32,13,10,32,32,93,32,32,32,32,32,32,32,32,32,32,32,32,124,32,124,32,124,32,32,32,32,32,32,32,32,32,32,32,32,124,32,32,13,10,32,32,93,32,32,32,32,32,32,32,32,32,32,32,32,32,126,32,126,32,32,32,32,32,32,32,32,32,32,32,32,32,124,32,32,13,10,32,32,124,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,124,32,32,32,13,10,32,32,32,124,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,32,124,32,32,32}
	fmt.Println(string(gopherByte))
}

func main() {
	url := flag.String("url", "", "URL video you want to download")

	flag.Parse()
	if *url == "" {
		fmt.Println("-url <URL video> is required")
		os.Exit(1)
	}

	videoURL := *url

	printGopher()
	fmt.Println("Finding download link of video: ", videoURL)

	getLinkFBVideo(videoURL)
}