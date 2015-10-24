package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"net/http/httputil"
	"net/http/cookiejar"
	"strings"
	"strconv"
	"regexp"
)
var ipRange = "192.168.4"
var linkPort1810 = []string{"47","48","49","50","51","52"}
var linkPort1820 = []string{"47","48","49","50","51","52"}
type ResponseAuth struct {
		Redirect	string `json:"redirect"`
		Error		string `json:"error"`
	}

var ipAndMacMapping = map[string]string{}

func main() {
	if !checkSwitch1810() {
		fmt.Println("Error : cannot get 1810")
		return
	}
	if !checkSwitch1820() {
		fmt.Println("Error : cannot get 1820")
		return
	}
	//build configuration file
	saveDhcpConf()
}

func saveDhcpConf(){
	var header = `
############ SETTING ############
######### set dhcp range ########

default-lease-time 600;
max-lease-time 7200;
subnet 192.168.4.0 netmask 255.255.255.0 {
	range 192.168.4.11 192.168.4.200;
	option subnet-mask 255.255.255.0;
	option broadcast-address 192.168.4.255;
	option routers 192.168.4.9;
	option domain-name-servers 8.8.8.8, 8.8.4.4;
	option domain-name "aiyara.lab.sut.ac.th";
} 

######### reserv ip  ########`
var body = ""
for ip,mac := range ipAndMacMapping {
		body = fmt.Sprintf("%s\nport-%s { hardware ethernet %s; fixed-address %s.%s; }", body, ip, mac, ipRange, ip)
}

err := ioutil.WriteFile("./dhcpd.conf", []byte(header+body), 0644)
if err != nil {
	fmt.Printf("error writefile: %s",err)
}
fmt.Println(header + body)
}

func checkSwitch1810() bool {
	var err error
	var password = ""
	// Create cookie.
  	cookieJar, _ := cookiejar.New(nil)

	contentReader := bytes.NewReader([]byte("pwd="+password))
	req, err := http.NewRequest("POST", "http://192.168.4.3/hp_login.html", contentReader)
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

	bodyByLine := strings.Split(string(body),"\n")
	for _,line := range bodyByLine {
		if line == `<td><INPUT class="inputfield" type="password" name="pwd" SIZE="10" MAXLENGTH="128" VALUE=""></td>`  {
			return false
		}
	}

	fmt.Printf("Status: %v\n", resp.Status)


	req, err = http.NewRequest("GET", "http://192.168.4.3/FDBSearch.html", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Cache-Control", "no-cache")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("error Request: %s", err)
	}
	body,err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error ReadAll: %s", err)
	}
	bodyByLine = strings.Split(string(body),"\n")
	//lineAttr := []string{}
	reg, err := regexp.Compile(`<tr><td CLASS="\w*" style="display: none;">|<td CLASS="\w*">|</td>|</td></tr>|</tr>`)
        if err != nil {
            fmt.Println(err)
        }
   
	for _,line := range bodyByLine {
		if strings.HasPrefix(line, `<tr><td CLASS=`) && strings.HasSuffix(line, `</td></tr>`)  {
			mapIpAndPort := reg.ReplaceAllString(line, ",")
			lineAttr := strings.Split(mapIpAndPort,",,")
			if _,err = strconv.Atoi(lineAttr[3]); err == nil {
				fmt.Println(lineAttr[3] + " " + lineAttr[1])
				if !contains(linkPort1810,lineAttr[3]) {
					ipAndMacMapping[lineAttr[3]] = lineAttr[1]
				}
			}
		}
	}

	// Logout to clear session (because it's limited).
	req, err = http.NewRequest("GET", "http://192.168.4.3/hp_login.html", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("error Request: %s", err)
	}
	return true
}

func checkSwitch1820() bool {
	var err error
	credential := map[string]string{
		"username": "admin",
		"password": "",
	}
	// Create cookie.
  	cookieJar, _ := cookiejar.New(nil)

	contentReader := bytes.NewReader([]byte("username="+credential["username"]+"&password="+credential["password"]))
	req, err := http.NewRequest("POST", "http://192.168.4.4/htdocs/login/login.lua", contentReader)
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

	fmt.Printf("Status: %v\n", resp.Status)
	if resp.Status != "200 OK" {
		return false
	}


	req, err = http.NewRequest("GET", "http://192.168.4.4/htdocs/pages/base/mac_address_table.lsp", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Cache-Control", "no-cache")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("error Request: %s", err)
	}
	body,err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error ReadAll: %s", err)
	}
	bodyByLine := strings.Split(string(body),"\n")
	for _,line := range bodyByLine {
		if strings.HasPrefix(line, "['") && strings.HasSuffix(line, "']")  {
			lineAttr := strings.Split(string(line),"', '")
			if ip,err := strconv.Atoi(lineAttr[2]); err == nil {
				ipBox2 := strconv.Itoa(ip+46)
				fmt.Println(ipBox2 + " " + lineAttr[1])
				if !contains(linkPort1820,lineAttr[2]) {
					ipAndMacMapping[ipBox2] = lineAttr[1]
				}

			}
		}
	}

	// Logout to clear session (because it's limited).
	req, err = http.NewRequest("GET", "http://192.168.4.4/htdocs/pages/main/logout.lsp", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("error Request: %s", err)
	}
	return true
}

func contains(slice []string, element string) bool {
    for _, item := range slice {
        if item == element {
            return true
        }
    }
    return false
}