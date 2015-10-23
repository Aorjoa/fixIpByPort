package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"net/http/httputil"
	"net/http/cookiejar"
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

  	cookieJar, _ := cookiejar.New(nil)

	contentReader := bytes.NewReader([]byte("username="+credential["username"]+"&password="+credential["password"]))
	req, err := http.NewRequest("POST", "http://192.168.1.1/htdocs/login/login.lua", contentReader)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("==== REQ ====\n%s\n=============\n",dump)
	client := &http.Client{
	    Jar: cookieJar,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error Request: %s", err)
	}
	defer resp.Body.Close()

	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error ReadAll: %s", err)
	}

	var responseAuth ResponseAuth
	err = json.Unmarshal(body, &responseAuth)
	if err != nil {
		fmt.Printf("error Unmarshal: %s", err)
	}
	fmt.Printf("Redirect: %v\n", responseAuth.Redirect)
	fmt.Printf("Error: %v\n", responseAuth.Error)
	fmt.Printf("Status: %v\n", resp.Status)


	req, err = http.NewRequest("GET", "http://192.168.1.1/htdocs/pages/base/mac_address_table.lsp", nil)
	if err != nil {
		fmt.Println(err)
	}
	// client := &http.Client{
	//     Jar: cookieJar,
	// }
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("error Request: %s", err)
	}
	body,err = ioutil.ReadAll(resp.Body)
	fmt.Printf("RES: %s\n", body)
//fmt.Printf("Cookie: %v\n", resp.Cookie)
	// resp, err = http.Get("http://192.168.1.1/htdocs/pages/base/mac_address_table.lsp")
	
 //    if err != nil {
 //        fmt.Printf("error MAC table: %s", err)
 //    }
 //    defer resp.Body.Close()
 //    body,err = ioutil.ReadAll(resp.Body)
 //    fmt.Printf("\n%s",body)
}