package main

import (
	"net/http"
	"time"
)

func linkValidation(url string) bool{
	client := http.Client {
		Timeout: 5 * time.Second, // set limit 5 seconds
	}

	//head is faster than get
	response, err := client.Head(url)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode >= 200 && response.StatusCode < 400
}