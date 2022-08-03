package models

import (
	"time"

	"gorm.io/gorm"
)

type Organization struct {
	ID                  uint32               `gorm:"primarykey"`
	Name                string               `json:"name"`
	OrganizationMembers []OrganizationMember `gorm:"foreignKey:OrganizationId"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}
