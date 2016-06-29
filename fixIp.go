package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"bufio"
	"strings"
	"strconv"	
	"regexp"
	"golang.org/x/crypto/ssh"
)

//Setting
var ipSw = "192.168.4.254"
var ipDhcpRange = "192.168.4"
var linkPortSw1 = []int{25,26,27,28}
var linkPortSw2 = []int{49,50,51,52}
var newIpLv = 24
var config = &ssh.ClientConfig{
    User: "cisco",
    Auth: []ssh.AuthMethod{
        ssh.Password("sut@1234"),
    },
    Config: ssh.Config{
			Ciphers: []string{"aes128-cbc"}, // not currently supported
	},
}

var ipAndMacMapping = map[int]string{}

func main() {
	if !checkSwitchSg500() {
		fmt.Println("Error : cannot get SG500")
		return
	}
	//build configuration file
	saveDhcpConf()
}


func checkSwitchSg500() bool {
		client, err := ssh.Dial("tcp", ipSw+":22", config)
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
			ssh.TTY_OP_ISPEED: 115200,
			ssh.TTY_OP_OSPEED: 115200,
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
			stdin.Write([]byte("term datadump\n"))
			stdin.Write([]byte("show mac address-table\n"))
			scanner := bufio.NewScanner(stdout)
			printMacTable := false
			regMacAddr := regexp.MustCompile(`([0-9a-f]{2}[:-]){5}([0-9a-f]{2})\s+[a-z]{2}\d+/\d+/\d+`)
			for scanner.Scan() {
				s := scanner.Text()
				findMac := regMacAddr.FindString(s)
				if len(findMac) > 0 {
					lineMacIpSplit := strings.Fields(findMac)
					getMac := lineMacIpSplit[0]
					splitVal := strings.Split(lineMacIpSplit[1],"/")
					getIp,_ := strconv.Atoi(splitVal[2])
					getIp = getIp+100;
					 
					if(splitVal[0] == "gi2"){
						getIp = getIp+newIpLv
						if !contains(linkPortSw2,getIp){
							ipAndMacMapping[getIp] = getMac
						}
					}else{
						if !contains(linkPortSw1,getIp){
							ipAndMacMapping[getIp] = getMac
						}
					}
				
				//lineBuffer = append(lineBuffer,s)
				}else if (strings.HasPrefix(s,"  Vlan        Mac Address         Port       Type    ")){
					printMacTable = true
				}
				if (len([]byte(s)) == 0 && printMacTable) {
					stdin.Write([]byte("exit\n"))
					client.Close()
					session.Close()
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}
	return true
}

func contains(slice []int, element int) bool {
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
	range 192.168.4.101 192.168.4.200;
	option subnet-mask 255.255.255.0;
	option broadcast-address 192.168.4.255;
	option routers 192.168.4.9;
	option domain-name-servers 8.8.8.8, 8.8.4.4;
	option domain-name "aiyara.lab.sut.ac.th";
} 

######### reserv ip  ########`
var body = ""
for ip,mac := range ipAndMacMapping {
		body = fmt.Sprintf("%s\nport-%d { hardware ethernet %s; fixed-address %s.%d; }", body, ip, mac, ipDhcpRange, ip)
}

err := ioutil.WriteFile("./dhcpd.conf", []byte(header+body), 0644)
if err != nil {
	fmt.Printf("error writefile: %s",err)
}
fmt.Println(header + body)
}
