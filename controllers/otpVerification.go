package controllers

import (
	"crypto/rand"
	"fmt"

	"net/http"
	"net/smtp"
	"os"

	"math/big"
	"strconv"

	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func VerifyOTP(email string) string {
	Otp, err := getRandNum()
	if err != nil {
		panic(err)
	}

	sendMail(email, Otp)
	return Otp
}

func sendMail(email string, otp string) {

	fmt.Println("Email : ", email, " otp :", otp)
	// Sender data.
	from := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.

	//  Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(otp+" is your otp"))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getRandNum() (string, error) {
	nBig, e := rand.Int(rand.Reader, big.NewInt(8999))
	if e != nil {
		return "", e
	}
	return strconv.FormatInt(nBig.Int64()+1000, 10), nil
}

//<-------Otp validtioin------------->

func OtpValidation(c *gin.Context) {
	type User_otp struct {
		Otp   string
		Email string
	}
	var user_otp User_otp
	var userData models.User
	if c.Bind(&user_otp) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON Data",
		})
		return
	}
	db := database.InitDB()
	result := db.First(&userData, "otp LIKE ? AND email LIKE ?", user_otp.Otp, user_otp.Email)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		// db.Last(&userData).Delete(&userData)
		c.JSON(422, gin.H{
			"Error":   "Wrong OTP Register Once agian",
			"Message": "Goto /signup/otpvalidate",
		})
		return
	}
	if result.Error == nil {
		db := database.InitDB()
		db.Model(&userData).Where("otp = ? AND email = ?", user_otp.Otp, user_otp.Email).Update("verified", true)
	
		c.JSON(200, gin.H{
			"Message": "OTP Validated. New User Successfully Registered. Go to the Login Page",
		})
		return
	}	

	c.JSON(200, gin.H{
		"Message": "New User Successfully Registered.Go to Login Page",
	})
}

func GenrateOtpForForgotPassword(c *gin.Context){
	type UserEnterData struct{
		Email  string
	}

	var data UserEnterData

	if c.Bind(&data)!=nil{
		c.JSON(400, gin.H{
			"Error": "Error when the data binding",
		})
		return
	}

	otp:=VerifyOTP(data.Email)

	DB:=database.InitDB()

	var userData models.User
	
	result:=DB.Model(userData).Where("email = ?",data.Email).Update("otp",otp)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "otp send go to /changepassword",
	})
}

func ChangePassword(c *gin.Context){
	type userEnterData struct {
		Email           string
		Otp             string
		Password        string
		ConfirmPassword string
	}

	var data userEnterData
	var userData models.User

	if c.Bind(&data)!=nil{
		c.JSON(400, gin.H{
			"Error": "Error when the data binding",
		})
		return
	}
	if data.Password!=data.ConfirmPassword{
		c.JSON(400, gin.H{
			"Error": "Password not match",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "False",
			"Error":  "Hashing password error",
		})
		return
	}
	DB:=database.InitDB()
	result := DB.Find(&userData, "email = ?", data.Email)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	}
	if data.Otp != userData.Otp {
		c.JSON(400, gin.H{
			"Error": "Invalide otp",
		})
		return
	}
	result = DB.Model(&userData).Where("email = ?", data.Email).Update("password", hash)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Password Change successfully",
	})
}


