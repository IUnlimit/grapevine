package database

import (
	"log"

	"github.com/illtamer/grapevine/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash default password: %v", err)
	}

	admin := model.User{
		Username:     "admin",
		PasswordHash: string(hash),
	}
	if err := DB.Create(&admin).Error; err != nil {
		log.Fatalf("failed to seed admin user: %v", err)
	}
	log.Println("default admin user created (admin/admin123)")
}
