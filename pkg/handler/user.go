package handler

import (
	"net/http"
	"zoom-api/gen/user/v1/userv1connect"
)

type UserServer struct {
	userv1connect.UserServiceHandler
}

func userHander() (string, http.Handler) {
	path, h := userv1connect.NewUserServiceHandler(&UserServer{})
	return path, h
}
