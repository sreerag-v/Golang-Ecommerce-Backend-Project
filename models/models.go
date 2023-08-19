package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	// "gorm.io/gorm"
)

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey;unique"  `
	First_Name   string `json:"first_name"  gorm:"not null" validate:"required,min=2,max=50"  `
	Last_Name    string `json:"last_name"    gorm:"not null"    validate:"required,min=1,max=50"  `
	Email        string `json:"email"   gorm:"not null;unique"  validate:"email,required"`
	Password     string `json:"password" gorm:"not null"  validate:"required"`
	Phone        string `json:"phone"   gorm:"not null;unique" validate:"required"`
	Otp          string `JSON:"otp"`
	Block_status bool   `json:"block_status " gorm:"not null"   `
	Verified     bool   `json:"verified " gorm:"not null"   `
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Admin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (admin *Admin) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	admin.Password = string(bytes)
	return nil
}
func (admin *Admin) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

type Product struct {
	gorm.Model
	ProductId   uint     `json:"product_id" gorm:"autoIncrement"  `
	ProductName string   `json:"product_name" gorm:"not null"  `
	Price       uint     `json:"price" gorm:"not null"  `
	Image       string   `json:"image" gorm:"not null"  `
	Stock       uint     `json:"stock"  `
	Color       string   `json:"color" gorm:"not null"  `
	Description string   `json:"description"   `
	Brand       Brand    `gorm:"ForeignKey:BrandId"`
	BrandId     uint     `json:"brand_id"`
	Catogery    Category `gorm:"ForeignKey:CatogeryId"`
	CatogeryId  uint     `json:"Catogery_id"`
}

type Brand struct {
	ID     uint   `json:"id" gorm:"primaryKey"  `
	Brands string `json:"brands" gorm:"not null"  `
}

type Category struct {
	ID         uint   `json:"id" gorm:"primaryKey"  `
	CategoryID uint   `json:"category_id" gorm:"autoIncrement" `
	Category   string `json:"category" `
}

type Cart struct {
	gorm.Model
	Product    Product `gorm:"ForeignKey:ProductId"`
	ProductId  uint
	Quantity   uint
	Price      uint
	TotalPrice uint
	Userid     uint
	User       User `gorm:"ForeignKey:Userid"`
}

type Address struct {
	Addressid uint `JSON:"addressid" gorm:"primarykey;unique"`

	User   User `gorm:"ForeignKey:Userid"`
	Userid uint `JSON:"uid"`

	Name       string `JSON:"name" gorm:"not null"`
	Phoneno    string `JSON:"phoneno" gorm:"not null"`
	Houseno    string `JSON:"houseno" gorm:"not null"`
	Area       string `JSON:"area" gorm:"not null"`
	Landmark   string `JSON:"landmark" gorm:"not null"`
	City       string `JSON:"city" gorm:"not null"`
	Pincode    string `JSON:"pincode" gorm:"not null"`
	District   string `JSON:"district" gorm:"not null"`
	State      string `JSON:"state" gorm:"not null"`
	Country    string `JSON:"country" gorm:"not null"`
	Defaultadd bool   `JSON:"defaultadd" gorm:"default:false"`
}

type Wallet struct {
	Id     uint
	User   User `gorm:"ForeignKey:UserId"`
	UserId uint
	Amount float64
}

type WalletHistory struct {
	Id             uint `JSON:"Id" gorm:"primarykey"`
	User           User `gorm:"ForeignKey:UserId"`
	UserId         uint
	Amount         float64
	TransctionType string
	Date           time.Time
}

type Coupon struct {
	ID            int
	CouponCode    string
	DiscountPrice float64
	CreatedAt     time.Time
	Expired       time.Time
}

type Wishlist struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	Userid    uint
	User      User    `gorm:"ForeignKey:Userid"`
	Product   Product `gorm:"ForeignKey:ProductId"`
	ProductId uint
}
