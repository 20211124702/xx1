package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Password  string         `json:"-" gorm:"not null;size:255"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100"`
	Nickname  string         `json:"nickname" gorm:"size:50"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Status    int            `json:"status" gorm:"default:1"` // 1:正常 0:禁用
	RoleID    int            `json:"role_id"`
	Role      Role           `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Role struct {
	ID          int            `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null;size:50"`
	Description string         `json:"description" gorm:"size:255"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Menu struct {
	ID        int            `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"not null;size:50"`
	Path      string         `json:"path" gorm:"size:100"`
	Component string         `json:"component" gorm:"size:100"`
	Icon      string         `json:"icon" gorm:"size:50"`
	Sort      int            `json:"sort" gorm:"default:0"`
	ParentID  *int           `json:"parent_id"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
