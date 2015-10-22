package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
//	"net/http/httptest"
//	"net/url"
//	"github.com/parnurzeal/gorequest"
)

func main() {
	credential := map[string]interface{}{
		"username": "admin",
		"password": "",
	}

	credentialJson, _ := json.Marshal(credential)
	contentReader := bytes.NewReader(credentialJson)
	req, _ := http.NewRequest("POST", "http://192.168.1.1/htdocs/login/login.lua", contentReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notes","Get MAC address form switch port")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}