package model

import (
	"errors"
	"zoom-api/database"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

// FindUserById id が一致するアカウントを返す
func (m User) FindUserById(id string) (User, error) {
	db := database.Open()
	if err := db.First(&m, id).Error; err != nil {
		return m, err
	}
	return m, nil
}

// FindUserByEmail email が一致するアカウントを返す
func (m User) FindUserByEmail() (User, error) {
	db := database.Open()
	if err := db.Where("email = ?", m.Email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return m, nil
		}
		return m, err
	}
	return m, nil
}

// UserCreate email が重複していない場合にユーザを作成する
func (m User) UserCreate() error {
	db := database.Open()
	user, err := m.FindUserByEmail()
	if err != nil {
		return err
	}
	err = db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
