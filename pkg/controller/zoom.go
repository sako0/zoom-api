package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"zoom-api/model"

	"github.com/labstack/echo/v4"
)

type CreateZoomRequestBody struct {
	Type int `json:"type"`
}

type CreateZoomResponse struct {
	Id       int    `json:"id"`
	JoinUrl  string `json:"join_url"`
	Password string `json:"password"`
	StartUrl string `json:"start_url"`
}

func CreateZoom() echo.HandlerFunc {
	return func(c echo.Context) error {
		zoom := new(model.Zoom)
		if err := c.Bind(zoom); err != nil {
			return err
		}
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
		req.Header.Add("Authorization", "Bearer "+zoom.Token)

		client := new(http.Client)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode != 201 {
			c.JSON(res.StatusCode, res.Status)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		response := new(CreateZoomResponse)
		json.Unmarshal(body, &response)

		return c.JSON(http.StatusOK, response)
	}
}

type Meeting struct {
	Id        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	JoinUrl   string `json:"join_url"`
	StartTime string `json:"start_time"`
	Topic     string `json:"topic"`
}
type MyZoomListResponse struct {
	Meetings []*Meeting `json:"meetings"`
}

func MyZoomList() echo.HandlerFunc {
	return func(c echo.Context) error {
		req, err := http.NewRequest("GET", "https://"+os.Getenv("ZOOM_DOMAIN")+"/users/me/meetings?type=live", nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.QueryParam("token"))
		client := new(http.Client)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode != 200 {
			fmt.Println(res.Status)
			c.JSON(res.StatusCode, res.Status)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		response := new(MyZoomListResponse)
		json.Unmarshal(body, &response)
		return c.JSON(http.StatusOK, response)
	}

}
