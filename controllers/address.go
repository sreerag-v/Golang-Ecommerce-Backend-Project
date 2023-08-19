package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

func AddAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	// var userName models.User
	var userEnterData models.Address

	if c.ShouldBind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in Binding the JSON",
		})
	}

	DB := database.InitDB()
	DB.Model(&models.Address{}).Where("userid = ?", id).Update("defaultadd", false)
	userEnterData.Userid = uint(id)
	result := DB.Create(&userEnterData)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	DB.Model(&userEnterData).Where("addressid = ?", userEnterData.Addressid).Updates(map[string]interface{}{
		"defaultadd": true,
		// "name":       userName.First_Name,
	})
	c.JSON(200, gin.H{
		"Message": "Address added succesfully",
	})
}

func ShowAddress(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var userAddres models.Address

	DB := database.InitDB()

	var count int64
	result := DB.Raw("SELECT * from addresses WHERE userid = ?", id).Scan(&userAddres).Count(&count)
	if count == 0 {
		c.JSON(500, gin.H{
			"message": "User not found",
		})
	}
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Address": gin.H{
			"Addressid":  userAddres.Addressid,
			"Userid":     userAddres.Userid,
			"Name":       userAddres.Name,
			"Phoneno":    userAddres.Phoneno,
			"Houseno":    userAddres.Houseno,
			"Area":       userAddres.Area,
			"Landmark":   userAddres.Landmark,
			"City":       userAddres.City,
			"Pincode":    userAddres.Pincode,
			"District":   userAddres.District,
			"State":      userAddres.State,
			"Country":    userAddres.Country,
			"Defaultadd": userAddres.Defaultadd,
		},
	})
}
