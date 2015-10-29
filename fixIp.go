package main

import (
	// "bytes"
	"fmt"
	"log"
	"io/ioutil"
//	"io"
	"os"
	"bufio"
//	"strings"	
//	"strconv"
//	"regexp"
	"golang.org/x/crypto/ssh"
)

//Setting
var ipSw = []string{"192.168.4.3","192.168.4.4"}
var ipDhcpRange = "192.168.4"
var linkPortSw = []string{"49","50","51","52"}
var config = &ssh.ClientConfig{
    User: "cisco",
    Auth: []ssh.AuthMethod{
        ssh.Password("sut@1234"),
    },
    Config: ssh.Config{
			Ciphers: []string{"aes128-cbc"}, // not currently supported
	},
}

var ipAndMacMapping = map[string]string{}

func main() {
	if !checkSwitchSg500() {
		fmt.Println("Error : cannot get SG500")
		return
	}
	//build configuration file
	saveDhcpConf()
}


func checkSwitchSg500() bool {
	for _,ip := range ipSw {
		client, err := ssh.Dial("tcp", ip+":22", config)
		if err != nil {
			panic("Failed to dial: " + err.Error())
		}
		defer client.Close()
		session, err := client.NewSession()
		if err != nil {
			log.Fatalf("unable to create session: %s", err)
		}
		defer session.Close()
		// Set up terminal modes
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		// Request pseudo terminal
			if err := session.RequestPty("vt100", 0, 200, modes); err != nil {
				log.Fatalf("request for pseudo terminal failed: %s", err)
			}
			stdin, err := session.StdinPipe()
			if err != nil {
				log.Fatalf("Unable to setup stdin for session: %v\n", err)
			}

			stdout, err := session.StdoutPipe()
			if err != nil {
				log.Fatalf("Unable to setup stdout for session: %v\n", err)
			}

		// Start remote shell
			if err := session.Shell(); err != nil {
				log.Fatalf("failed to start shell: %s", err)
			}
			stdin.Write([]byte("show mac address-table\r\n"))
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Println(scanner.Text()) // Println will add back the final '\n'
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}

	}
	fmt.Println("\n<<< SSH SESSION EXPIRED >>>")
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
		body = fmt.Sprintf("%s\nport-%s { hardware ethernet %s; fixed-address %s.%s; }", body, ip, mac, ipDhcpRange, ip)
}

err := ioutil.WriteFile("./dhcpd.conf", []byte(header+body), 0644)
if err != nil {
	fmt.Printf("error writefile: %s",err)
}
fmt.Println(header + body)
}