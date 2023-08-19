package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sreerag_v/Ecom/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() *gorm.DB {

	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db_URL := os.Getenv("DNS")
	db, err := gorm.Open(postgres.Open(db_URL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nConnected to DATABASE: ", db.Name())
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Admin{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Brand{})
	db.AutoMigrate(&models.Address{})
	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.Payment{})
	db.AutoMigrate(&models.Oder_item{})
	db.AutoMigrate(&models.OderDetails{})
	db.AutoMigrate(&models.RazorPay{})
	db.AutoMigrate(&models.Wallet{})
	db.AutoMigrate(&models.WalletHistory{})
	db.AutoMigrate(&models.Coupon{})
	db.AutoMigrate(&models.Wishlist{})

	return db
}
