package controller

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"zoom-api/model"

	"github.com/labstack/echo/v4"
)

type LoginResponse struct {
	Email string
	Name  string
}

//ユーザをgetする
func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		// ssl認証を無視する
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		auth := new(model.Auth)
		if err := c.Bind(auth); err != nil {
			return err
		}
		// ヘッダーについて定義
		req, _ := http.NewRequest("GET", "https://"+os.Getenv("AUTH0_DOMAIN")+"/api/v2/users/"+auth.UserId, nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+auth.Token)

		client := new(http.Client)
		res, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)
		response := new(LoginResponse)
		json.Unmarshal(body, &response)
		user := new(model.User)
		user.Email = response.Email
		user.Name = response.Name
		targetUser, err := user.FindUserByEmail()
		if err != nil {
			c.JSON(http.StatusInternalServerError, "FindUserByEmail Error")
		}
		if targetUser.ID == 0 {
			err := targetUser.UserCreate()
			if err != nil {
				c.JSON(http.StatusInternalServerError, "UserCreate Error")
			}
		}
		return c.JSON(http.StatusOK, targetUser.ID)
	}
}

// func decode_jwt_proc(str_token string) string {
// 	aaa := strings.Split(str_token, ".")
// 	str_bbb := strings.Replace(aaa[1], "-", "+", -1)
// 	str_ccc := strings.Replace(str_bbb, "_", "/", -1)

// 	llx := len(str_ccc)
// 	nnx := ((4 - llx%4) % 4)
// 	ssx := strings.Repeat("=", nnx)
// 	str_ddd := strings.Join([]string{str_ccc, ssx}, "")
// 	ppp, err := b64.StdEncoding.DecodeString(str_ddd)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "*** error *** StdEncoding.DecodeString ***\n")
// 		fmt.Println("error:", err)
// 		return "error"
// 	}

// 	uEnc := b64.URLEncoding.EncodeToString([]byte(ppp))
// 	decode, _ := b64.URLEncoding.DecodeString(uEnc)

// 	return string(decode)
// }
