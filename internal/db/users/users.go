package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	tableNameUsers = "users"
)

type User struct {
	ID        uuid.UUID `json:"id"         gorm:"column:id;type:uuid;primaryKey"`
	Email     string    `json:"email"      gorm:"column:email;type:varchar(255);unique;not null"`
	FirstName string    `json:"first_name" gorm:"column:first_name;type:varchar(255)"`
	LastName  string    `json:"last_name"  gorm:"column:last_name;type:varchar(255)"`
	Groups    []string  `json:"groups"     gorm:"column:groups;type:jsonb;serializer:json"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;column:updated_at;not null"`
	LastLogin time.Time `json:"last_login" gorm:"column:last_login;type:timestamp"`
}

func (u *User) TableName() string {
	return tableNameUsers
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func CreateUser(ctx context.Context, db *gorm.DB, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

func UpdateUser(ctx context.Context, db *gorm.DB, userID string, updates *User) error {
	return db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(updates).Error
}

func GetUserByEmail(ctx context.Context, db *gorm.DB, email string) (*User, error) {
	var user User
	if err := db.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(ctx context.Context, db *gorm.DB, id string) (*User, error) {
	var user User
	if err := db.WithContext(ctx).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
