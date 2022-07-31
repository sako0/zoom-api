package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	userv1 "zoom-api/gen/user/v1"
	"zoom-api/gen/user/v1/userv1connect"
	"zoom-api/pkg/models"
)

type UserServer struct {
}

func (us *UserServer) CreateUser(
	ctx context.Context,
	req *connect.Request[userv1.CreateUserRequest],
) (*connect.Response[userv1.CreateUserResponse], error) {
	log.Println("Request headers: ", req.Header())
	fmt.Println("connect")
	user := models.User{Auth0Id: req.Msg.Auth0Id, Name: req.Msg.Name, Email: req.Msg.Email, ZoomToken: req.Msg.ZoomToken, ZoomRefreshToken: req.Msg.ZoomRefreshToken}
	err := user.UserCreate()
	if err != nil {
		fmt.Println(err)
	}
	res := connect.NewResponse(&userv1.CreateUserResponse{
		Id: int32(user.ID),
	})
	res.Header().Set("User-Version", "v1")
	return res, nil
}

func main() {
	// envファイル読み込み
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}
	mux := http.NewServeMux()
	fmt.Println("connect")
	// user
	path, h := userv1connect.NewUserServiceHandler(&UserServer{})
	fmt.Println(h)
	handler := cors.AllowAll().Handler(h) // cors
	mux.Handle(path, handler)

	err := http.ListenAndServe(
		"0.0.0.0:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		fmt.Println(err)
	}
}
