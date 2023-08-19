package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

func OrderDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	var UserAddress models.Address
	var UserPayment models.Payment
	var UserCart []models.Cart

	DB := database.InitDB()

	result := DB.Find(&UserAddress, "userid = ? AND defaultadd = true", id) // if error record not found
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": err.Error(),
		})
	}

	result = DB.Last(&UserPayment, "user_id = ?", id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
	}

	result = DB.Find(&UserCart, "userid = ?", id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	var oder_item models.Oder_item
	DB.Last(&oder_item, "user_id_no = ?", id)

	for _, UserCart := range UserCart {
		OrderDetails := models.OderDetails{
			Userid:     uint(id),
			AddressId:  UserAddress.Addressid,
			PaymentId:  UserPayment.PaymentId,
			ProductId:  UserCart.ProductId,
			Status:     "pending",
			Quantity:   UserCart.Quantity,
			OderItemId: oder_item.OrderId,
		}
		result = DB.Create(&OrderDetails)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"Message": "Oder Added succesfully",
	})
}

func ShowOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	DB := database.InitDB()
	var userOrder []models.OderDetails
	result := DB.Find(&userOrder, "userid = ?", id)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error":   "Failed to fetch user orders",
			"Message": "Something went wrong. Please try again later.",
		})
		return
	}

	type data struct {
		Name     string
		Phoneno  string
		Houseno  string
		Area     string
		Landmark string
		City     string
		Pincode  string
		District string
		State    string
		Country  string
	}
	var userAddressDatas []data

	result = DB.Raw("SELECT name, phoneno, houseno, area, landmark, city, pincode, district, state, country FROM addresses WHERE defaultadd = true AND userid = ?", id).Scan(&userAddressDatas)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error":   "Failed to fetch user address",
			"Message": "Something went wrong. Please try again later.",
		})
		return
	}

	var orderDetailsList []gin.H
	for _, order := range userOrder {
		var products []models.Product
		result := DB.Find(&products, "product_id = ?", order.ProductId)
		if result.Error != nil {
			c.JSON(500, gin.H{
				"Error":   "Failed to fetch product details for order",
				"Message": "Something went wrong. Please try again later.",
			})
			continue
		}

		// Check if products slice is empty before accessing its elements
		if len(products) == 0 {
			c.JSON(404, gin.H{
				"Error":   "Product not found",
				"Message": fmt.Sprintf("Product not found for order with ID: %d", order.Oderid),
			})
			continue
		}

		orderDetails := gin.H{
			"Order_id":         order.Oderid,
			"Product name":     products[0].ProductName,
			"Price":            products[0].Price,
			"Description":      products[0].Description,
			"Quantity":         order.Quantity,
			"Shipping Address": userAddressDatas,
		}
		orderDetailsList = append(orderDetailsList, orderDetails)
	}

	// Return the list of order details as a single JSON response
	c.JSON(200, orderDetailsList)
}

func CancelOrder(c *gin.Context) {

	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	order_itemID := c.Query("order_itemid")

	var orderDetails models.OderDetails
	var orderItem models.Oder_item
	var wallet models.Wallet

	DB := database.InitDB()

	err = DB.First(&orderItem, order_itemID).Error

	if err != nil {
		c.JSON(400, gin.H{
			"Error": "order id does't exist",
		})
		return
	}

	if orderItem.OrderStatus == "canceled" {
		c.JSON(400, gin.H{
			"Error": "Oder already canceled",
		})
		return
	}

	result := DB.Model(&orderDetails).Where("userid = ? AND oder_item_id = ? ", id, order_itemID).Update("status", "Canceled")
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	result = DB.Model(&orderItem).Where("order_id = ?", order_itemID).Update("order_status", "Canceled")
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	// var check models.Oder_item
	// var payment string
	// DB.Where("order_id = ?", order_itemID).Select("payment_m").Find(&check).Scan(&payment)
	// // adding the balance amount into wallet
	// fmt.Println("payment", payment)
	// if payment != "COD" {
	// 	fmt.Println("....................")

	if orderItem.PaymentM != "COD" {
		result = DB.Where("user_id", id).First(&wallet)
		if result.Error != nil {
			walletData := models.Wallet{
				UserId: uint(id),
				Amount: float64(orderItem.TotalAmount),
			}
			result = DB.Create(&walletData)
			if result.Error != nil {
				c.JSON(400, gin.H{
					"Error": result.Error.Error(),
				})
				return
			}
		} else {
			totalAmount := wallet.Amount + float64(orderItem.TotalAmount)
			fmt.Println("this is the added amount : ", totalAmount)

			result = DB.Model(&wallet).Where("user_id = ?", id).Update("amount", totalAmount)
			if result.Error != nil {
				c.JSON(400, gin.H{
					"Error": "Error occurd while adding the amoutn into the wallet",
				})
				return
			}
		}

		wHistory := models.WalletHistory{
			UserId:         uint(id),
			Amount:         float64(orderItem.TotalAmount),
			TransctionType: "Credit",
			Date:           time.Now(),
		}

		result = DB.Create(&wHistory)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": "Error occurd while adding the amoutn into the wallet",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"Massage": "Order canceld",
	})
}

func ReturnOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	oderid, err := strconv.Atoi(c.Query("orderid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	var order models.OderDetails
	var oderItem models.Oder_item
	DB := database.InitDB()

	result := DB.Model(&order).Where("userid = ? AND oder_item_id  = ?", id, oderid).Update("status", "Return product")
	if result.Error != nil {

		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	result = DB.Model(&oderItem).Where("user_id_no = ? AND order_id = ?", id, oderid).Update("order_status", "Return product")
	if result.Error != nil {

		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Massage": "Product Return",
	})
}
