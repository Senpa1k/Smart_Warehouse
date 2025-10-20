package main

import (
	"log"
	"time"

	"github.com/Senpa1k/Smart_Warehouse/internal/config"
	"github.com/Senpa1k/Smart_Warehouse/internal/models"
)

// type Env struct{
// 	Jwt_secret string   getenv()
// }

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("can not open db")
	}

	var user models.Users
	var userp models.Users = models.Users{Name: "keru",
		Email:        "new",
		PasswordHash: "feq",
		Role:         "admin"}

	db.Create(&userp)
	for {
		time.Sleep(1 * time.Second)
		db.First(&user, 1)
		log.Println(user.Name)
	}

}
