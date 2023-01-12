package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type WolConf struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Mac  string `json:"mac"`
}

var jsonConf []WolConf

/**
*  Get Current Path
 */
func getPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return pwd
}

/**
*  Get Config File Path
 */
func getConfPath() string {
	return fmt.Sprintf("%s/%s", getPwd(), "wol.conf")
}

/**
*  Get Config Content
 */
func getConf(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return content
}

/**
*  Get Template Config
 */
func getTplConf() []WolConf {
	return []WolConf{
		{
			Name: "Test",
			IP:   "127.0.0.1",
			Mac:  "AA-BB-CC-DD-EE-FF",
		},
	}
}

/**
*  Get Loaded Config Content
 */
func getLoadConf(wc []WolConf) []string {
	loadConf := []string{}
	for _, v := range wc {
		loadConf = append(loadConf, fmt.Sprintf("Name: \t%s\tIP: \t%s\t Mac:\t%s\t", v.Name, v.IP, v.Mac))
	}
	return loadConf
}

/**
*  Get Mac Address From Args
 */
func getMacFromArgs(wc []WolConf) string {
	mac := os.Args[1]
	for _, v := range wc {
		if v.Name == os.Args[1] {
			mac = v.Mac
		}
	}
	return mac
}

/**
*  Send Wol Magic Packet
 */
func sendMagicPacket(hw net.HardwareAddr) {
	magicPacket := append(bytes.Repeat([]byte{0xff}, 6), bytes.Repeat(hw, 16)...)
	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		return
	}

	fmt.Printf("Sending Magic Packet To %s \n", hw.String())
	conn.Write(magicPacket)
	conn.Close()
}

func main() {
	path := getConfPath()
	conf := getConf(path)
	if len(conf) < 1 {
		jm, err := json.Marshal(getTplConf())
		if err == nil {
			os.WriteFile(path, jm, 0644)
		}
	}
	err := json.Unmarshal(conf, &jsonConf)
	if err != nil {
		jsonConf = []WolConf{}
	}
	if len(os.Args) < 2 {
		fmt.Printf("Usage:\t wol AA-BB-CC-DD-EE-FF\n%s\n", strings.Join(getLoadConf(jsonConf), "\t"))
		return
	}
	hw, err := net.ParseMAC(getMacFromArgs(jsonConf))
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	sendMagicPacket(hw)
}
