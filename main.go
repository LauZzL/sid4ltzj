package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/bytedance/sonic"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"github.com/qtgolang/SunnyNet/public"
	"github.com/zhangboqi/cfmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var NodeURL = "http://localhost:3000/"

func main() {
	args := os.Args
	commandLine(args)
}

func commandLine(args []string) {
	arg := args[1]
	switch arg {
	case "-i":
		startSunnyNet()
		break
	case "-q":
		stopSunnyNet()
		break
	}
	time.Sleep(24 * time.Hour)
}

func Request(uri string, data string) string {
	url := NodeURL + uri
	method := "POST"
	payload := strings.NewReader(`{"data":"` + data + `"}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err.Error()
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	return string(body)
}

func stopSunnyNet() {
	fmt.Println("Stop SunnyNet")
	s := SunnyNet.NewSunny()
	err := s.Error
	if err != nil {
		panic(err)
	}
	s.SetIeProxy(true)
	s.Close()
	fmt.Println("Stop SunnyNet Done.")
}

func startSunnyNet() {
	fmt.Println("Check Node Serve")
	for {
		resp := Request("x", "")
		if resp == "1" {
			break
		}
		cfmt.Yprintln("Node serve not running, look at https://github.com/LauZzL/leitingzhanji")
		time.Sleep(3 * time.Second)
	}
	cfmt.Gprintln("Check Node Serve Done.")

	fmt.Println("Start SunnyNet")
	s := SunnyNet.NewSunny()
	s.SetGoCallback(HttpCallback, TcpCallback, WSCallback, UdpCallback)
	Port := 2024
	s = s.SetPort(Port).Start()
	// 开启随机TLS指纹
	s.SetRandomTLS(true)
	// 安装证书
	s.InstallCert()
	s.SetIeProxy(false)
	err := s.Error
	if err != nil {
		panic(err)
	}
	cfmt.Gprintln("Start SunnyNet Done.")
}

func HttpCallback(Conn *SunnyNet.HttpConn) {

	URL := Conn.Request.URL.String()
	urlEquals := URL == "https://wxmini.jj5agame.com/p.f"
	if Conn.Type == public.HttpSendRequest && urlEquals {
		if Conn.Request.Body != nil {
			Body, _ := io.ReadAll(Conn.Request.Body)
			_ = Conn.Request.Body.Close()
			// 将body进行base64编码
			bs64Body := base64.StdEncoding.EncodeToString(Body)
			Conn.Request.Body = io.NopCloser(bytes.NewBuffer(Body))
			decryptResult := Request("d", bs64Body)
			root, _ := sonic.GetFromString(decryptResult)
			sid, _ := root.Get("head").Get("sid").String()
			uid, _ := root.Get("head").Get("uid").String()

			if sid != "" && uid != "" {
				cfmt.Gprintln("Copy to clipboard: sid:" + sid + " uid:" + uid)
				_ = clipboard.WriteAll(`sid:` + sid + ` uid:` + uid)
			}
		}
	}

}
func WSCallback(Conn *SunnyNet.WsConn) {
}
func TcpCallback(Conn *SunnyNet.TcpConn) {
}
func UdpCallback(Conn *SunnyNet.UDPConn) {
}
