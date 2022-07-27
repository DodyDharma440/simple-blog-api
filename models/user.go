package models

import (
	"final-project/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashPw, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPw), []byte(pw))
}

func (u *User) BeforeSave() error {
	hashPw, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashPw)
	return nil
}

func (u *User) LoginCheck(db *gorm.DB) (string, error) {
	user := User{}

	if err := db.Model(User{}).Where("email=?", u.Email).Take(&user).Error; err != nil {
		return "", err
	}

	if err := VerifyPassword(user.Password, u.Password); err != nil {
		return "", err
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}
