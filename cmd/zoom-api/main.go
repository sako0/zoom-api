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
	userInfo, err := getUserInfo(req.Msg.ZoomToken)
	if err != nil {
		fmt.Println("getUserInfoでエラー")
	}
	user := models.User{Auth0Id: req.Msg.Auth0Id, ZoomUserId: userInfo.Id, Name: req.Msg.Name, Email: req.Msg.Email, ZoomToken: req.Msg.ZoomToken, ZoomRefreshToken: req.Msg.ZoomRefreshToken}
	err = user.UserCreateOrUpdate()
	if err != nil {
		fmt.Println(err)
	}
	res := connect.NewResponse(&userv1.CreateUserResponse{
		Id: int32(user.ID),
	})
	res.Header().Set("User-Version", "v1")
	return res, nil
}

type GetUserInfoResponse struct {
	Id string `json:"id"`
}

func getUserInfo(zoomToken string) (*GetUserInfoResponse, error) {
	req, err := http.NewRequest("GET", "https://"+os.Getenv("ZOOM_DOMAIN")+"/users/me", nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+zoomToken)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if res.StatusCode != 200 {
		fmt.Println(res.Status)
		panic("/api/v2/users/me取得中の通信エラー")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	response := new(GetUserInfoResponse)
	json.Unmarshal(body, &response)

	return response, nil
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
		Id:        uint32(user.ID),
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

func (zs *ZoomServer) GetZoomList(
	ctx context.Context,
	req *connect.Request[zoomv1.GetZoomListRequest],
) (*connect.Response[zoomv1.GetZoomListResponse], error) {
	om := models.OrganizationMember{OrganizationId: 1}
	members, err := om.GetOrganizationMemberListByOrganizationId()
	if err != nil {
		panic("GetZoomListのmembersが取得できていない")
	}

	zoomInfoList := []*zoomv1.ZoomInfo{}
	// user := models.User{}
	// db := database.Open()
	for _, member := range members {
		// err := db.First(&user, member.UserId).Error
		fmt.Println(member.User.Name)
		if err != nil {
			fmt.Println("DBの読み込みに失敗")
			return nil, err
		}
		zoomList, err := getZoomListByZoomUserId(member.User.ZoomUserId, req.Msg.GetZoomToken())
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		info := &zoomv1.ZoomInfo{ZoomUserId: member.User.ZoomUserId, ZoomMeetingList: zoomList}
		zoomInfoList = append(zoomInfoList, info)
	}
	res := connect.NewResponse(&zoomv1.GetZoomListResponse{ZoomList: zoomInfoList})
	res.Header().Set("Zoom-Version", "v1")
	return res, nil
}

type ZoomMeeting struct {
	Id        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	JoinUrl   string `json:"join_url"`
	StartUrl  string `json:"start_url"`
	Topic     string `json:"topic"`
	HostId    string `json:"host_id"`
}
type GetZoomListByZoomUserIdResponse struct {
	Meetings []*ZoomMeeting `json:"meetings"`
}

func getZoomListByZoomUserId(zoomUserId string, zoomToken string) ([]*zoomv1.ZoomMeetingInfo, error) {
	req, err := http.NewRequest("GET", "https://"+os.Getenv("ZOOM_DOMAIN")+"/users/"+zoomUserId+"/meetings?type=live", nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+zoomToken)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if res.StatusCode != 200 {
		fmt.Println(res.Status)
		panic("/meetings?type=live取得中の通信エラー")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := new(GetZoomListByZoomUserIdResponse)
	json.Unmarshal(body, &response)
	zoomMeetingInfoList := []*zoomv1.ZoomMeetingInfo{}
	for _, meeting := range response.Meetings {
		info := &zoomv1.ZoomMeetingInfo{Id: uint32(meeting.Id), CreatedAt: meeting.CreatedAt, JoinUrl: meeting.JoinUrl, StartUrl: meeting.StartUrl, Topic: meeting.Topic, HostId: meeting.HostId}
		zoomMeetingInfoList = append(zoomMeetingInfoList, info)
	}

	return zoomMeetingInfoList, err

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
