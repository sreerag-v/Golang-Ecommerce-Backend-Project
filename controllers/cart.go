package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

func AddToCart(c *gin.Context) {
	type data struct {
		Product_id uint
		Quantity   uint
	}

	var bindData data
	var productData models.Product

	
	if c.Bind(&bindData) != nil {
		c.JSON(400, gin.H{
			"Bad Request": "Could not bind the JSON data",
		})
		return
	}

	id, err := strconv.Atoi(c.GetString("userid"))

	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	DB := database.InitDB()

	result := DB.First(&productData, bindData.Product_id)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"Message": "Product not exist",
		})
		return
	}

	// cheklcing stoke qunatity is exist

	if bindData.Quantity > productData.Stock {
		c.JSON(404, gin.H{
			"Message": "Out of Stock",
		})
		return
	}

	var sum uint
	var Price uint

	
	err = DB.Table("carts").Where("product_id = ? AND userid = ? ", bindData.Product_id, id).Select("quantity", "total_price").Row().Scan(&sum, &Price)
	fmt.Println("this is the erro : ", err)

	if err != nil {
		totalPrice := productData.Price * bindData.Quantity

		cartitems := models.Cart{
			ProductId:  bindData.Product_id,
			Quantity:   bindData.Quantity,
			Price:      productData.Price,
			TotalPrice: totalPrice,
			Userid:     uint(id),
		}

		result := DB.Create(&cartitems)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"Message": "Added to the Cart Successfull",
		})
		return
	}

	totalQuantity := sum + bindData.Quantity
	totalPrice := productData.Price * totalQuantity

	result = DB.Model(&models.Cart{}).Where("product_id = ? AND userid = ? ", bindData.Product_id, id).Updates(map[string]interface{}{"quantity": totalQuantity, "total_price": totalPrice})
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Message": "Quantity added Successfully",
	})
}

// view cart items user id
func ViewCart(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))

	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	type cartData struct {
		Product_name string
		Quantity     uint
		TotalPrice   uint
		Price        uint
	}

	var datas []cartData

	DB := database.InitDB()

	result := DB.Table("carts").
		Select("products.product_name, carts.quantity, carts.price, carts.total_price").
		Joins("INNER JOIN products ON products.product_id=carts.product_id").Where("userid = ?", id).Scan(&datas)

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	fmt.Println("this is carts data : ", datas)
	if datas != nil {
		c.JSON(200, gin.H{
			"cart items": datas,
		})
	} else {
		c.JSON(404, gin.H{
			"Message": "Cart is empty",
		})
	}
}

// delete cart item
func DeleteCart(c *gin.Context){
	id:=c.Param("id")
	userid, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	DB:=database.InitDB()

	result := DB.Exec("delete from carts where id= ? AND userid = ?", id, userid)
	count := result.RowsAffected
	if count == 0 {
		c.JSON(400, gin.H{
			"Message": "Cart not exist",
		})
		return
	}
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Cart Items": "Delete successfully",
	})
}
