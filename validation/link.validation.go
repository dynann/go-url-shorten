package validation

import (
	"fmt"
	"net/http"
	"time"
)

func LinkValidation(url string) bool {
	client := http.Client{
		Timeout: 10 * time.Second, 
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	//head is faster than get
	response, err := client.Do(req)
	fmt.Println("Error ==> ✅✅✅",err)
	if err != nil {
		fmt.Print(err)
		return false
	}
	defer response.Body.Close()
	return response.StatusCode >= 200 && response.StatusCode < 400
}