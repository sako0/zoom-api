package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type CreateZoomRequestBody struct {
	Type int `json:"type"`
}

type CreateZoomResponse struct {
	Id       int    `json:"id"`
	JoinUrl  string `json:"join_url"`
	Password string `json:"password"`
	StartUrl string `json:"start_url"`
	Topic    string `json:"topic"`
}

func GetZoomInfo(zoomToken string) *CreateZoomResponse {
	requestBody := &CreateZoomRequestBody{
		Type: 1,
	}
	jsonString, err := json.Marshal(requestBody)
	if err != nil {
		panic("Error")
	}
	req, err := http.NewRequest("POST", "https://"+os.Getenv("ZOOM_DOMAIN")+"/users/me/meetings", bytes.NewReader(jsonString))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+zoomToken)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode != 201 {
		panic("zoom-api中の通信エラー")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	response := new(CreateZoomResponse)
	json.Unmarshal(body, &response)

	return response
}
