package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `json:"id" gorm:"type:char(36);not null;primary_key;unique_index"`
	FirstName string         `json:"first_name" gorm:"type:char(50);not null"`
	LastName  string         `json:"last_name" gorm:"char(50);not null"`
	Email     string         `json:"email"  gorm:"type:varchar(50);not null"`
	Phone     string         `json:"phone" gorm:"type:char(20);not null"`
	CreateAt  *time.Time     `json:"-"`
	UpdateAt  *time.Time     `json:"-"`
	Deleted   gorm.DeletedAt `json:"-"`
}
