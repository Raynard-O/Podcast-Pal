package library

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Get request handler
func Getl(url, listen_api, FFullurl string) (*bytes.Buffer, int64, error) {

	client := http.Client{
		//Timeout: time.Duration(20 * time.Second),
	}
	// fetch listen note api key from env
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}
	// set header
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("X-ListenAPI-Key", listen_api)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	request, err := client.Get(FFullurl)
	defer request.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Panic(err)
	}
	size := int64(len(b))
	l := bytes.NewBuffer(b)

	return l, size, nil
}

// Get request handler
func Get(url, listen_api string) ([]byte, int64, error) {
	var size int64
	response, err := http.Get(url)
	response.Header.Set("Content-type", "application/json")
	response.Header.Set("X-ListenAPI-Key", listen_api)
	if err != nil {
		//fmt.Printf("%s", err)
		//os.Exit(1)
		log.Fatal(err.Error())
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		size = int64(len(contents))
		if err != nil {
			//fmt.Printf("%s", err)
			//os.Exit(1)
			log.Fatal(err.Error())
		}
		return contents, 0, nil
	}
	return nil, size, nil
}

// Post request handler
func Post(url string, jsonData string, auth string) string {
	var jsonStr = []byte(jsonData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

// Delete request handler
func Delete(url string, jsonData string) string {
	var jsonStr = []byte(jsonData)

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
