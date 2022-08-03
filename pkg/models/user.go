package models

import (
	"errors"
	"time"
	"zoom-api/pkg/database"

	"gorm.io/gorm"
)

type User struct {
	ID                  uint32               `gorm:"primarykey"`
	Auth0Id             string               `json:"auth0_id"`
	ZoomUserId          string               `json:"zoom_user_id"`
	Name                string               `json:"name"`
	Email               string               `json:"email"`
	ZoomToken           string               `json:"zoom_token"`
	ZoomRefreshToken    string               `json:"zoom_refresh_token"`
	OrganizationMembers []OrganizationMember `gorm:"foreignKey:UserId"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

// FindUserByAuth0Id Auth0Id が一致するアカウントを返す
func (m User) FindUserByAuth0Id() (User, error) {
	db := database.Open()
	err := db.Where("auth0_id = ?", m.Auth0Id).First(&m).Error
	if err != nil {
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

// UserCreate Auth0Id が重複していない場合にユーザを作成する
func (m User) UserCreateOrUpdate() error {
	db := database.Open()
	user, err := m.FindUserByAuth0Id()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Create(&m).Error
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	err = db.First(&user, user.ID).Updates(User{ZoomToken: m.ZoomToken, ZoomRefreshToken: m.ZoomRefreshToken}).Error
	if err != nil {
		return err
	}
	return nil

}
