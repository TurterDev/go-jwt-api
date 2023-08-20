package orm

import "gorm.io/gorm"

// สร้าง Table ใน Database
type User struct {
	gorm.Model
	Username string
	Password string
	Fullname string
	Avatar   string
}
