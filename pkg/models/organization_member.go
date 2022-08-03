package models

import (
	"time"
	"zoom-api/pkg/database"

	"gorm.io/gorm"
)

type OrganizationMember struct {
	ID             uint32 `gorm:"primarykey"`
	UserId         uint32
	User           User
	OrganizationId uint32
	Organization   Organization
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (m OrganizationMember) GetOrganizationMemberListByOrganizationId() ([]OrganizationMember, error) {
	db := database.Open()
	ms := []OrganizationMember{}
	err := db.Preload("User").Where("organization_id = ?", m.OrganizationId).Find(&ms).Error
	if err != nil {
		return nil, err
	}
	return ms, nil

}
