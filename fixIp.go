package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"net/http/httputil"
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
	credential := map[string]string{
		"username": "admin",
		"password": "",
	}

	contentReader := bytes.NewReader([]byte("username="+credential["username"]+"&password="+credential["password"]))
	req, _ := http.NewRequest("POST", "http://192.168.1.1/htdocs/login/login.lua", contentReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("%s",dump)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Status: %v\n", resp.Status)
}