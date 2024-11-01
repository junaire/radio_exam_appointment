package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

type LoginResponse struct {
	Ret  int    `json:"ret"`
	Data string `json:"data"`
}

type ExamPlanResponse struct {
	Data struct {
		List []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"list"`
	} `json:"data"`
}

type ReserveResponse struct {
	Ret  int    `json:"ret"`
	Data string `json:"data"`
}

func login(client *http.Client, username, password string) {
	fmt.Println("Try to login...")
	data := url.Values{
		"username": {username},
		"password": {password},
		"type":     {"account"},
	}

	req, err := http.NewRequest("POST", "https://xt.bjwxdxh.org.cn/loginapi/login", bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResp)

	if os.Getenv("DEBUG") != "" {
		jsonData, err := json.Marshal(loginResp)
		if err != nil {
			fmt.Println("Error:", err)
			log.Fatalf("Error sending request: %v", err)
		}
		fmt.Println(string(jsonData))
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Fail to login")
	}

	if loginResp.Ret != 1 {
		log.Fatal("Fail to login, check your username/password")
	}

	fmt.Println("Login succeed!")
}

func getPossiblePlan(client *http.Client) string {
	fmt.Println("Try to select exam plan...")
	req, err := http.NewRequest("GET", "https://xt.bjwxdxh.org.cn/memberapi/getExamPlanCanSelected?pageIndex=1", nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	var planResp ExamPlanResponse
	json.NewDecoder(resp.Body).Decode(&planResp)

	if os.Getenv("DEBUG") != "" {
		jsonData, err := json.Marshal(planResp)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		fmt.Println(string(jsonData))
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Fail to select exam plan")
	}

	planID := planResp.Data.List[0].ID
	fmt.Println("Selected:", planResp.Data.List[0].Title)
	return planID
}

func reserve(client *http.Client, planID string) {
	fmt.Println("Try to reserve the exam...")
	data := url.Values{
		"id": {planID},
	}

	req, err := http.NewRequest("POST", "https://xt.bjwxdxh.org.cn/memberapi/addMemberExamPlan", bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	var reserveResp ReserveResponse
	json.NewDecoder(resp.Body).Decode(&reserveResp)
	if os.Getenv("DEBUG") != "" {
		jsonData, err := json.Marshal(reserveResp)
		if err != nil {
			log.Fatalf("Error sending request: %v", err)
		}
		fmt.Println(string(jsonData))
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Fail to reserve the exam!")
	}
	if reserveResp.Ret != 1 {
		log.Fatal("Fail to reserve the exam")

	}

	fmt.Println("OK!")
}

func setCommonHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("DNT", "1")
	req.Header.Add("Origin", "https://xt.bjwxdxh.org.cn")
	req.Header.Add("Referer", "https://xt.bjwxdxh.org.cn/static/member/")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	req.Header.Add("sec-ch-ua", `"Chromium";v="130", "Google Chrome";v="130", "Not?A_Brand";v="99"`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", `"Linux"`)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./main <username> <password>")
		return
	}

	username := os.Args[1]
	password := os.Args[2]
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}
	login(client, username, password)
	planID := getPossiblePlan(client)
	reserve(client, planID)
}
