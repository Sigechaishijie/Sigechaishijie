package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/88250/gulu"
	"github.com/parnurzeal/gorequest"
)

var logger = gulu.Log.NewLogger(os.Stdout)

const (
	githubUserName = "Sigechaishijie"
)

func main() {
	result := map[string]interface{}{}
	
	request := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("User-Agent", "Profile Bot; +https://github.com/"+githubUserName+"/"+githubUserName)
	
	response, data, errors := request.Get("https://api.example.com/events").EndStruct(&result)
	
	if nil != errors || http.StatusOK != response.StatusCode {
		logger.Fatalf("fetch events failed: %+v, %s", errors, data)
	}
	
	if code, ok := result["code"].(float64); ok && code != 0 {
		logger.Fatalf("fetch events failed: %s", data)
	}

	buf := &bytes.Buffer{}
	buf.WriteString("\n\n")
	updated := time.Now().Format("2006-01-02 15:04:05")
	buf.WriteString("我的近期动态（点个 [Star](https://github.com/" + githubUserName + "/" + githubUserName + ") 将触发自动刷新，最近更新时间：`" + updated + "`）：\n\n")
	
	if dataSlice, ok := result["data"].([]interface{}); ok {
		for _, event := range dataSlice {
			if evt, ok := event.(map[string]interface{}); ok {
				operation := evt["operation"].(string)
				title := evt["title"].(string)
				url := evt["url"].(string)
				content := evt["content"].(string)
				buf.WriteString("* [" + operation + "](" + url + ")：（" + title + "）" + content + "\n")
			}
		}
	}
	buf.WriteString("\n")

	fmt.Println(buf.String())

	readme, err := ioutil.ReadFile("README.md")
	if nil != err {
		logger.Fatalf("read README.md failed: %s", err)
	}

	startFlag := []byte("<!--events start -->")
	startIndex := bytes.Index(readme, startFlag)
	if startIndex == -1 {
		logger.Fatalf("start flag not found in README.md")
	}
	beforeStart := readme[:startIndex+len(startFlag)]
	
	endFlag := []byte("<!--events end -->")
	endIndex := bytes.Index(readme, endFlag)
	if endIndex == -1 {
		logger.Fatalf("end flag not found in README.md")
	}
	afterEnd := readme[endIndex:]

	newReadme := append(beforeStart, buf.Bytes()...)
	newReadme = append(newReadme, afterEnd...)
	
	if err := ioutil.WriteFile("README.md", newReadme, 0644); nil != err {
		logger.Fatalf("write README.md failed: %s", err)
	}
}
