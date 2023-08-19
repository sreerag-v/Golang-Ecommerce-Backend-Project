package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/auth"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

type AdminLogins struct {
	Email    string
	Password string
}

func AdminSignup(c *gin.Context) {
	var admin models.Admin
	var count uint
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
		c.Abort()
		return
	}

	database.InitDB().Raw("select count(*) from admins where email=?", admin.Email).Scan(&count)
	if count > 0 {
		c.JSON(400, gin.H{
			"status": "false",
			"msg":    "an admin with same email already exists",
		})
		c.Abort()
		return
	}

	if err := admin.HashPassword(admin.Password); err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
	}
	record := database.InitDB().Create(&admin)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"error": record.Error.Error(),
		})
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"msg":    "Admin Created",
	})
}

func AdminLogin(c *gin.Context) {
	var u AdminLogins
	var admin models.Admin
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	record := database.InitDB().Raw("select * from admins where email=?", u.Email).Scan(&admin)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	credentialcheck := admin.CheckPassword(u.Password)
	if credentialcheck != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		c.Abort()
		return
	}
	tokenstring, err := auth.GenerateJWT(u.Email)
	token := tokenstring["access_token"]
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Adminjwt", token, 3600*24*30, "", "", false, true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      true,
		"message":     "ok",
		"tokenstring": tokenstring,
	})
}

func AdminHome(c *gin.Context) {
	fmt.Println("hai")
	c.JSON(202, gin.H{
		"msg": "Welcome to Admin Panel",
	})
}

func AdminLogout(c *gin.Context) {

}

// Admin Controlling

type Userdata struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Phone     int
}

func UserData(c *gin.Context) {
	count, err1 := strconv.Atoi(c.Query("count"))
	PageN, err2 := strconv.Atoi(c.Query("page"))
	err3 := errors.Join(err1, err2)

	if err3 != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	limit := count
	offset := (PageN - 1) * limit
	var userData []Userdata
	db := database.InitDB()

	result := db.Table("users").Select("id,first_name, last_name, email,phone").
		Limit(limit).
		Offset(offset).
		Scan(&userData)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Message": "Could not find the users",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User data": userData,
	})
}

type Order struct {
	UserIdNo    uint
	OrderId     uint
	OrderStatus string
	PaymentM    string
	TotalAmount uint
	CreatedAt   time.Time
}

func ShowAllOrders(c *gin.Context) {
	var Orders []Order

	DB := database.InitDB()

	result := DB.Table("oder_items").Select("order_id,user_id_no, order_status, payment_m, total_amount, created_at").
		Order("user_id_no ASC").
		Scan(&Orders)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Message": "Could not find the orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"All Users Orders": Orders,
	})
}

type OrderNew struct {
	OrderId     uint
	OrderStatus string
	PaymentM    string
	TotalAmount uint
	CreatedAt   time.Time
}

func ShowOrderById(c *gin.Context) {
	param := c.Param("id")

	var Orders []OrderNew
	DB := database.InitDB()

	result := DB.Table("oder_items").Where("user_id_no = ?", param).Select("order_id,order_status, payment_m, total_amount, created_at").
		Order("order_id ASC").
		Scan(&Orders)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		" User Orders": Orders,
	})
}

func BlockUser(c *gin.Context) {
	params := c.Param("id")
	var user models.User
	database.InitDB().Raw("UPDATE users SET block_status=true where id=?", params).Scan(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "Blocked successfully"})
}
func UnBlockUser(c *gin.Context) {
	params := c.Param("id")
	var user models.User
	database.InitDB().Raw("UPDATE users SET block_status=false where id=?", params).Scan(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "Unblocked successfully"})
}
