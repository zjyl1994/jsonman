package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tidwall/pretty"
)

const (
	promptJSON    = "\033[31;1mJSON >>\033[0m"
	promptResult  = "\033[31;1mRESULT >>\033[0m"
	promptError   = "\033[31;1mERROR >>\033[0m"
	promptMessage = "\033[31;1mMESSAGE >>\033[0m"
	logoText      = `
 _ ____ ____ _  _ _  _ ____ _  _ 
 | [__  |  | |\ | |\/| |__| |\ | 
_| ___] |__| | \| |  | |  | | \| 
`
	promptBanner = "\033[31;1m" + logoText + "\n(C) 2020 zjyl1994.com\033[0m\n"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error param")
		return
	}
	postURL := os.Args[1]
	fmt.Println(promptBanner)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(promptJSON)
		var jsonStr string
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Println(promptError, err.Error())
				}
				return
			}
			trimed := strings.TrimSpace(text)
			if strings.HasPrefix(trimed, "~") { // is command line
				switch trimed {
				case "~exit":
					return
				case "~clear":
					jsonStr = ""
					fmt.Print(promptJSON)
					continue
				}
			}
			jsonStr += text
			if json.Valid([]byte(jsonStr)) {
				break
			}
		}
		fmt.Println(promptMessage, "Posting to "+postURL)
		status, header, content, timeUsed, err := postJSON(postURL, jsonStr)
		if err != nil {
			fmt.Println(promptError, err.Error())
			continue
		}
		fmt.Println(promptResult, status, "\033[31;1min\033[0m", timeUsed.String())
		fmt.Println()
		for k, v := range header {
			for _, vv := range v {
				fmt.Printf("\033[37;1m%s\033[0m: %s\n", k, vv)
			}
		}
		fmt.Println()
		if json.Valid(content) {
			fmt.Println(string(pretty.Color(pretty.Pretty(content), nil)))
		} else {
			fmt.Println(promptError, "Response not valid json")
		}
	}
}

func postJSON(url, content string) (status string, header http.Header, responseBody []byte, timeUsed time.Duration, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(content)))
	if err != nil {
		return "", nil, nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	startTime := time.Now()
	resp, err := client.Do(req)
	timeSince := time.Since(startTime)
	if err != nil {
		return "", nil, nil, timeSince, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return resp.Status, resp.Header, body, timeSince, err
}
