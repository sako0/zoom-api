package services

// type UserServer struct {
// 	userv1connect.UserServiceHandler
// }

// func (us *UserServer) CreateUser(
// 	ctx context.Context,
// 	req *connect.Request[userv1.CreateUserRequest],
// ) (*connect.Response[userv1.CreateUserResponse], error) {
// 	log.Println("Request headers: ", req.Header())
// 	user := models.User{Auth0Id: req.Msg.Auth0Id, Name: req.Msg.Name, Email: req.Msg.Email, ZoomToken: req.Msg.ZoomToken, ZoomRefreshToken: req.Msg.ZoomRefreshToken}
// 	err := user.UserCreate()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	if user.ID <= 0 {
// 		fmt.Print("既に存在")
// 	}
// 	res := connect.NewResponse(&userv1.CreateUserResponse{
// 		Id: int32(user.ID),
// 	})
// 	res.Header().Set("User-Version", "v1")
// 	return res, nil
// }
