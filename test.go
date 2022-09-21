package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type wcSendcontent struct {
	Content string `json:"content"`
}

type WcSendMsg struct {
	MsgType string        `json:"msgtype"`
	Text    wcSendcontent `json:"text"`
}

// ExecShell ...
func ExecShell(command string, arg ...string) (out string, err error) {
	var Stdout []byte
	cmd := exec.Command(command, arg...)
	Stdout, err = cmd.CombinedOutput()
	out = string(Stdout)
	return
}

// Repo ...
func Repo() (repo string, err error) {
	var (
		out string
	)
	if out, err = ExecShell("/bin/sh", "-c", "git remote -v"); err != nil {
		return
	}
	if repo = out[strings.Index(out, ":")+1 : strings.Index(out, ".git")]; repo == "" {
		err = fmt.Errorf("not found, %s", out)
		return
	}
	return
}

// Branch ...
func Branch() (branch string, err error) {
	var (
		out string
	)
	if out, err = ExecShell("/bin/sh", "-c", "git branch"); err != nil {
		return
	}
	list := strings.Split(out, "\n")
	for _, v := range list {
		if strings.HasPrefix(v, "*") {
			branch = v[strings.Index(v, "*")+2:]
			return
		}
	}
	err = fmt.Errorf("not found, %s", out)
	return
}

//git diff --name-only HEAD~ HEAD
func main() {
	//var Stdout []byte/
	cmd, _ := ExecShell("git", "diff", "--name-only", "HEAD~", "HEAD")
	//Stdout, _ = cmd.CombinedOutput()
	//out := string(Stdout)
	fmt.Println(cmd)
	//SendCardMsg()
	SendCardMsg(cmd)
}

//获取到

//企业微信应用消息提醒方法如下
func SendCardMsg(fliename string) (WcSendMsg, error) {

	flie_diff := fmt.Sprintf("%s 词条文件发生改变", fliename)

	req := WcSendMsg{MsgType: "text", Text: wcSendcontent{Content: flie_diff}}

	sendurl := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=cd75d3a3-0899-4c63-a1a7-fe578784b9e2"
	data, err := httpPostJson(sendurl, req)
	if err != nil {
		log.Println(err)
		return WcSendMsg{MsgType: "", Text: wcSendcontent{Content: ""}}, err
	}
	return data, nil
}

func httpPostJson(url string, data WcSendMsg) (WcSendMsg, error) {
	res, err := json.Marshal(data)
	if err != nil {
		return WcSendMsg{MsgType: "", Text: wcSendcontent{Content: ""}}, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(res))
	if err != nil {
		return WcSendMsg{MsgType: "", Text: wcSendcontent{Content: ""}}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return WcSendMsg{MsgType: "", Text: wcSendcontent{Content: ""}}, err
	}
	return data, nil
}
