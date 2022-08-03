package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	userv1 "zoom-api/gen/user/v1"
	"zoom-api/gen/user/v1/userv1connect"
	"zoom-api/pkg/models"

	zoomv1 "zoom-api/gen/zoom/v1"
	"zoom-api/gen/zoom/v1/zoomv1connect"

	"github.com/bufbuild/connect-go"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type UserServer struct {
	userv1connect.UserServiceHandler
}

func (us *UserServer) CreateUser(
	ctx context.Context,
	req *connect.Request[userv1.CreateUserRequest],
) (*connect.Response[userv1.CreateUserResponse], error) {
	log.Println("Request headers: ", req.Header())
	user := models.User{Auth0Id: req.Msg.Auth0Id, Name: req.Msg.Name, Email: req.Msg.Email, ZoomToken: req.Msg.ZoomToken, ZoomRefreshToken: req.Msg.ZoomRefreshToken}
	err := user.UserCreateOrUpdate()
	if err != nil {
		fmt.Println(err)
	}
	res := connect.NewResponse(&userv1.CreateUserResponse{
		Id: int32(user.ID),
	})
	res.Header().Set("User-Version", "v1")
	return res, nil
}

type ZoomServer struct {
	zoomv1connect.ZoomServiceHandler
}

func (zs *ZoomServer) CreateZoom(
	ctx context.Context,
	req *connect.Request[zoomv1.CreateZoomRequest],
) (*connect.Response[zoomv1.CreateZoomResponse], error) {
	user := models.User{Auth0Id: req.Msg.Auth0Id}
	user, err := user.FindUserByAuth0Id()
	if err != nil {
		panic(err)
	}
	zoomInfo := getZoomInfo(user.ZoomToken)
	res := connect.NewResponse(&zoomv1.CreateZoomResponse{
		Id:        int64(user.ID),
		CreatedAt: user.CreatedAt.String(),
		JoinUrl:   zoomInfo.JoinUrl,
		StartUrl:  zoomInfo.StartUrl,
		Topic:     zoomInfo.Topic,
	})
	res.Header().Set("Zoom-Version", "v1")
	return res, nil
}

type CreateZoomRequestBody struct {
	Type int `json:"type"`
}

func getZoomInfo(zoomToken string) *zoomv1.CreateZoomResponse {
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
		fmt.Println(res.Status)
		panic("zoom-api中の通信エラー")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	response := new(zoomv1.CreateZoomResponse)
	json.Unmarshal(body, &response)

	return response
}

func main() {
	// envファイル読み込み
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	// user
	mux := http.NewServeMux()
	userPath, userH := userv1connect.NewUserServiceHandler(&UserServer{})
	userHandler := cors.AllowAll().Handler(userH)
	mux.Handle(userPath, userHandler)

	// zoom
	zoomPath, zoomH := zoomv1connect.NewZoomServiceHandler(&ZoomServer{})
	zoomHandler := cors.AllowAll().Handler(zoomH)
	mux.Handle(zoomPath, zoomHandler)

	err := http.ListenAndServe(
		"0.0.0.0:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		fmt.Println(err)
	}
}
