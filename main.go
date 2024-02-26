package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"robber/config"
	"robber/module"
	"time"
)

var UserInfo = module.LoginInfo{}
var SolverAddress = ""

func init() {
	config.SetUpEnv()
	// load .env
	UserInfo.Username = os.Getenv("USERNAME")
	UserInfo.Password = os.Getenv("PASSWORD")
	SolverAddress = os.Getenv("CAPTCHA_SOLVER_IP")
	fmt.Println(UserInfo, SolverAddress)
}

func main() {
	// 初始化 HTTP 客戶端
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: time.Second * 30,
	}

	// 獲取驗證碼圖片
	image, err := getValidationImage(client)
	if err != nil {
		panic(err)
	}

	// 儲存圖片到磁盤
	imagePath := "validateCode.jpg"
	err = os.WriteFile(imagePath, image, 0644)
	if err != nil {
		panic(err)
	}

	// 上傳圖片
	code, err := uploadImage(imagePath, fmt.Sprintf("http://%s/recognize-text", SolverAddress))
	if err != nil {
		panic(err)
	}
}

func getValidationImage(client *http.Client) ([]byte, error) {
	req, _ := http.NewRequest("GET", "https://webap.nkust.edu.tw/nkust/validateCode.jsp", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36")
	req.Header.Set("Referer", "https://webap.nkust.edu.tw/nkust/index_main.html?1111")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func uploadImage(imagePath, url string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", file.Name())
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// message will be like { text : "code"}
	type Message struct {
		Text string `json:"text"`
	}
	var message Message
	err = json.NewDecoder(response.Body).Decode(&message)
	if err != nil {
		return "", err
	}

	return message.Text, nil
}
