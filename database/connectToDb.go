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
	dns := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
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
