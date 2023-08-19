package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

func AddBrand(c *gin.Context) {
	var addbrand models.Brand

	if c.Bind(&addbrand) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind JSON data",
		})
		return
	}

	DB := database.InitDB()

	result := DB.Create(&addbrand)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Message":       "New Brand added Successfully",
		"Brand details": addbrand,
	})
}

func ViewBrand(c *gin.Context) {
	var brandData []models.Brand
	DB := database.InitDB()
	result := DB.First(&brandData)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Message": "Brand is empty",
		})
		return
	}
	c.JSON(200, gin.H{
		"Brands data": brandData,
	})
}

func EditBrand(c *gin.Context) {
	bid := c.Param("id")
	id, err := strconv.Atoi(bid)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var editbrands models.Brand
	if c.Bind(&editbrands) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	editbrands.ID = uint(id)
	DB := database.InitDB()

	result := DB.Model(&editbrands).Updates(models.Brand{
		Brands: editbrands.Brands,
	})

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully updated the Brand",
	})
}

func AddCategories(c *gin.Context) {
	type Data struct {
		Category string
	}
	var category Data
	var CategoryData models.Category

	if c.Bind(&category) != nil {
		c.JSON(400, gin.H{
			"Error": "countl not bind the JSON data",
		})
	}

	DB := database.InitDB()
	var count int64
	result := DB.Find(&CategoryData, "category = ?", category.Category).Count(&count)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}
	if count == 0 {
		createData := models.Category{
			Category: category.Category,
		}
		result = DB.Create(&createData)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
		}
		c.JSON(200, gin.H{
			"message":  "Catogery created",
			"Catogery": createData,
		})
	} else {
		c.JSON(400, gin.H{
			"message": "Catogery already exist",
		})
	}
}

// add wishlist
func Wishlist(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))

	if err != nil {
		c.JSON(400, gin.H{
			"err": "cant convert string",
		})
	}
	type dataN struct {
		Product_id uint
	}

	var productId dataN

	if c.Bind(&productId) != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
	}

	DB := database.InitDB()
	// check the product is exist
	var product models.Product
	result := DB.First(&product, productId.Product_id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"err": "product does not exit",
		})
		return
	}
	data := models.Wishlist{
		ProductId: productId.Product_id,
		Userid:    uint(id),
	}

	result = DB.Create(&data)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"message": "Wish list added sucessfully",
	})
}
