package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sreerag_v/Ecom/auth"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

var validate = validator.New()

func Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(404, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(404, gin.H{"err": err.Error()})
		c.Abort()
		return
	}
	otp := VerifyOTP(user.Email)
	result2 := database.InitDB().Create(&user)
	if result2.Error != nil {
		c.JSON(500, gin.H{
			"Status": "False",
			"Error":  result2.Error.Error(),
		})
	} else {
		database.InitDB().Model(&user).Where("email LIKE ?", user.Email).Update("otp", otp)

		c.JSON(200, gin.H{
			"message": "Go to /signup/otpvalidate",
		})
	}
}

type UserLogin struct {
	Email        string
	Password     string
	Block_status bool
}

func LoginUser(c *gin.Context) {
	var ulogin UserLogin
	var user models.User
	if err := c.ShouldBindJSON(&ulogin); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	record := database.InitDB().Raw("select * from users where email=?", ulogin.Email).Scan(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	if !user.Verified{
		db := database.InitDB()
        db.Delete(&user)

        c.JSON(422, gin.H{
            "Error":   "User is not verified. Data deleted.",
            "Message": "Please complete OTP verification to complete registration.",
        })
        return
	}

	if user.Block_status {
		c.JSON(404, gin.H{"msg": "user has been blocked By admin"})
		c.Abort()
		return
	}
	credentialcheck := user.CheckPassword(ulogin.Password)
	if credentialcheck != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		c.Abort()
		return
	}

	str := strconv.Itoa(int(user.ID))
	tokenString, err := auth.GenerateJWT(str)
	fmt.Println(tokenString)
	token := tokenString["access_token"]
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuth", token, 3600*24*30, "", "", false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"email": ulogin.Email, "password": ulogin.Password, "token": tokenString})
}

func UserHome(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "welcome User Home"})

}

func LogoutUser(c *gin.Context) {
	c.SetCookie("UserAuth", "", -1, "", "", false, false)
	c.JSON(200, gin.H{
		"Message": "User Successfully  Log Out",
	})
}
