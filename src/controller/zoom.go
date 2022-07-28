package controller

import (
	"bytes"
	"crypto/tls"
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
		// ssl認証を無視する
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

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
		fmt.Println(string(body))
		response := new(CreateZoomResponse)
		json.Unmarshal(body, &response)

		return c.JSON(http.StatusOK, response)
	}
}
