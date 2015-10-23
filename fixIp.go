package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
//	"strings"
//	"net/http/httptest"
//	"net/url"
//	"github.com/parnurzeal/gorequest"
)

type ResponseAuth struct {
		Redirect	string `json:"redirect"`
		Error		string `json:"error"`
	}

func main() {
	var err error
	credential := map[string]interface{}{
		"username": "admin",
		"password": "",
	}

	credentialJson, _ := json.Marshal(credential)
	contentReader := bytes.NewReader(credentialJson)
	req, _ := http.NewRequest("POST", "http://192.168.1.1/htdocs/login/login.lua", contentReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Notes","Get MAC address form switch port")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error ReadAll:", err)
	}
	var responseAuth ResponseAuth
	err = json.Unmarshal(body, &responseAuth)
	if err != nil {
		fmt.Println("error Unmarshal:", err)
	}
	fmt.Printf("Redirect: %v\n", responseAuth.Redirect)
	fmt.Printf("Error: %v\n", responseAuth.Error)
}