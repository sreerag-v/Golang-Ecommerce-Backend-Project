package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

type ProfileData struct {
	Firstname   string
	Lastname    string
	Email       string
	PhoneNumber string
}

func ShowUserDetails(c *gin.Context) {
	var userData models.User
	id, err := strconv.Atoi(c.GetString("userid"))

	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	db := database.InitDB()
	result := db.First(&userData, "id = ?", id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": "User not exist",
		})
		return
	}
	c.JSON(200, gin.H{

		"First name":   userData.First_Name,
		"Last name":    userData.Last_Name,
		"Email":        userData.Email,
		"Phone number": userData.Phone,
	})
}

func EditUserProfile(c *gin.Context) {
	var userEnterData ProfileData
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"error": "Data binding error",
		})
		return
	}

	var userData models.User
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	db := database.InitDB()
	result := db.First(&userData, "id = ?", id)
	if result.Error != nil {
		c.JSON(409, gin.H{
			"Error": "User not exist",
		})
		return
	}
	result = db.Model(&userData).Updates(models.User{
		First_Name: userEnterData.Firstname,
		Last_Name:  userEnterData.Lastname,
		Phone:      userEnterData.PhoneNumber,
	})

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Message": "Successfully Updated the profile",
		"Updated data": gin.H{
			"First name": userData.First_Name,
			"Last name":  userData.Last_Name,
			"Phone":      userData.Phone,
		},
	})
}
