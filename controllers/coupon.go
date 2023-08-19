package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

func AddCoupon(c *gin.Context) {

	type data struct {
		CouponCode    string
		Year          uint
		Month         uint
		Day           uint
		DiscountPrice float64
		Expired       time.Time
	}

	var userEnterData data
	var couponData []models.Coupon

	DB := database.InitDB()

	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}

	specificTime := time.Date(int(userEnterData.Year), time.Month(userEnterData.Month), int(userEnterData.Day), 0, 0, 0, 0, time.UTC)

	userEnterData.Expired = specificTime
	var count int64
	result := DB.First(&couponData, "coupon_code = ?", userEnterData.CouponCode).Count(&count)
	if result.Error != nil {
		Data := models.Coupon{
			CouponCode:    userEnterData.CouponCode,
			DiscountPrice: userEnterData.DiscountPrice,
			Expired:       userEnterData.Expired,
		}
		result := DB.Create(&Data)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
		}
		c.JSON(200, gin.H{
			"message": userEnterData,
		})
	} else {
		c.JSON(400, gin.H{
			"message": "Coupon already exist",
		})
	}
}

// chekc the coupen is valid or exist in this database
func CheckCoupon(c *gin.Context) {
	type data struct {
		Coupon string
	}
	var coupon models.Coupon
	var userEnterData data

	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})

	}

	db := database.InitDB()

	var count int64
	result := db.Find(&coupon, "coupon_code = ?", userEnterData.Coupon).Count(&count)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if count == 0 {
		c.JSON(400, gin.H{
			"message": "Coupon not exist",
		})
		return
	}
	currentTime := time.Now()
	expiredData := coupon.Expired

	if currentTime.Before(expiredData) {
		c.JSON(200, gin.H{
			"message": "Coupon valide",
		})
	} else if currentTime.After(expiredData) {
		c.JSON(400, gin.H{
			"message": "Coupon expired",
		})
	}
}

func ApplyCoupen(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error while string conversion",
		})
	}

	type data struct {
		Coupon string
	}

	var CoupenData data
	var Coupon models.Coupon
	var discountPercentage float64
	var discountPrice float64

	if c.Bind(&CoupenData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}

	DB := database.InitDB()
	//checking coupon is existig or not
	var count int64
	result := DB.Find(&Coupon, "coupon_code = ?", CoupenData.Coupon).Count(&count)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if count == 0 {
		c.JSON(400, gin.H{
			"message": "Coupon not exist",
		})
		return

	} else {
		currentTime := time.Now()
		expiredData := Coupon.Expired

		if currentTime.Before(expiredData) {

			c.JSON(200, gin.H{
				"message": "Coupon valide",
			})
			discountPercentage = Coupon.DiscountPrice

		} else if currentTime.After(expiredData) {

			c.JSON(400, gin.H{
				"message": "Coupon expired",
			})
		}
	}

	//fetching the cart details from the table carts
	ViewCart(c)
	//fetching and calculatin the total amount of the cart products
	var totalPrice float64
	result1 := DB.Table("carts").Where("userid = ?", id).Select("SUM(total_price)").Scan(&totalPrice).Error

	//calculating the discount amount
	discountPrice = discountPercentage
	totalPriceAfterDeduct := totalPrice - discountPrice

	if result1 != nil {
		c.JSON(400, gin.H{
			"Error": "Can not fetch total amount",
		})
		return
	}

	c.String(200, "Coupen-Applied")
	DB.Table("carts").Where("userid = ?", id).Update("total_price", totalPriceAfterDeduct)

	var coupon models.Coupon
	if err := DB.Where("coupon_code = ?", CoupenData.Coupon).First(&coupon).Error; err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Delete the coupon record
	if err := DB.Delete(&coupon).Error; err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Price":        totalPrice,
		"Discount":     discountPrice,
		"Total amount": totalPriceAfterDeduct,
	})

}
