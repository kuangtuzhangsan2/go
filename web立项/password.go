package main

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null;unique" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"password"`
	Role     string `gorm:"type:varchar(20);not null" json:"role"`
	Passproblem string `gorm:"type:varchar(255);not null" json:"passproblem"`
}

// 将密码存入 Password 字段
func (u *User) SetPassword(pwd string) error {
	u.Password = pwd
	return nil
}

func (u *User) SetPassproblem(pwd string) error {
	u.Passproblem = pwd
	return nil
}

// 比较密码
func (u *User) CheckPassword(pwd string) bool {
	return u.Password == pwd
}

func (u *User) CheckPassproblem(pwd string) bool {
	return u.Passproblem == pwd
}
