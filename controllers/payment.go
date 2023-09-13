package controllers

import (
	// "crypto/hmac"
	// "crypto/sha256"
	// "encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

func DeleteCartItems(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}
	// var cartData models.Cart
	DB := database.InitDB()
	// result := db.Where("userid = ?", id).Delete(&cartData)
	result := DB.Exec("delete from carts where userid = ?", id)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
}

func CashOnDelivery(c *gin.Context) {
	// fetch the user with the token
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
		return
	}

	var cartData models.Cart

	DB := database.InitDB()

	// fetching the data from the table carts by id
	result := DB.First(&cartData, "userid = ?", id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Message": "Cart is empty",
		})
		return
	}

	//fteching the total amount from thr table carts
	var total_amount float64
	result = DB.Table("carts").Where("userid = ?", id).Select("SUM(total_price)").Scan(&total_amount)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "Error fetching the total amount from the table carts",
		})
		return
	}

	todaysDate := time.Now()
	paymentData := models.Payment{
		PaymentMethod: "COD",
		Totalamount:   uint(total_amount),
		Date:          todaysDate,
		Status:        "pending",
		UserId:        uint(id),
	}

	result = DB.Create(&paymentData)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	var addressData models.Address
	result = DB.First(&addressData, "userid = ?", id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "address not exist",
		})
		return
	}

	orderData := models.Oder_item{
		UserIdNo:    uint(id),
		TotalAmount: uint(total_amount),
		PaymentId:   paymentData.PaymentId,
		AddId:       addressData.Addressid,
		PaymentM:    paymentData.PaymentMethod,
		OrderStatus: "success",
	}

	result = DB.Create(&orderData)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": "Creating the table order",
		})
		return
	}

	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error retrieving order details",
		})
		return
	}
	OrderDetails(c)
	response := gin.H{
		"Message": "Payment Method COD",
		"Status":  "True",
	}

	c.JSON(200, response)

	DeleteCartItems(c)
}

func Razorpay(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	DB := database.InitDB()

	var userdata models.User
	// fetch the user id
	result := DB.Find(&userdata, "id = ?", id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	// fetch the total price from the table carts
	var amount uint
	row := DB.Table("carts").Where("userid = ?", id).Select("SUM(total_price)").Row()
	err = row.Scan(&amount)

	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
	}
	//Sending the payment details to Razorpay
	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET"))
	data := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	//Creating the payment details to client order
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err,
		})
		return
	}

	//To rendering the html page with user&payment details
	value := body["id"]

	c.HTML(200, "app.html", gin.H{
		"userid":     userdata.ID,
		"totalprice": amount,
		"paymentid":  value,
	})
}

// when the Razorpay payment is completed this funcion will work
func RazorpaySuccess(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	DB := database.InitDB()

	//fetching the payment details from Razorpay
	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")

	totalamount := c.Query("total")

	//Creating table razorpay  using the data from Razorpay
	Rpay := models.RazorPay{
		UserID:          uint(userID),
		RazorPaymentId:  paymentid,
		Signature:       signature,
		RazorPayOrderID: orderid,
		AmountPaid:      totalamount,
	}

	fmt.Println(Rpay)
	result := DB.Create(&Rpay)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	//....................//

	// 	var transactionData models.RazorPay
	//     var data = transactionData.RazorPayOrderID + "|" + transactionData.RazorPaymentId
	//     var KEY = []byte(os.Getenv("RAZORPAY__SECRET"))
	//    //
	//   hash := hmac.New(sha256.New, []byte(KEY))
	//   hash.Write([]byte(data))
	//   genarated_signature := hex.EncodeToString(hash.Sum(nil))

	//    if transactionData.Signature != genarated_signature {
	// 	  c.JSON(200, gin.H{
	// 		"Error": "Transaction not verified",
	// 	     })
	// 	   deleteRazorpay(transactionData)
	// 	   return
	//     } else {
	// 	   c.JSON(200, gin.H{
	// 		"Message": "Signature Check verified",
	// 	     })
	//     }

	//............................//

	todyDate := time.Now()
	method := "Razor Pay"
	status := "pending"

	//converting to string total amount veriable
	totalprice, err := strconv.Atoi(totalamount)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion--",
			"err":   err.Error(),
		})
		return
	}

	//Creating payment table
	paymentdata := models.Payment{
		UserId:        uint(userID),
		PaymentMethod: method,
		Status:        status,
		Date:          todyDate,
		Totalamount:   uint(totalprice),
	}
	result1 := DB.Create(&paymentdata)
	if result1.Error != nil {
		c.JSON(400, gin.H{
			"Error": result1.Error.Error(),
		})
		return
	}

	var addressData models.Address
	result = DB.First(&addressData, "userid = ?", userID)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	pid := paymentdata.PaymentId

	oderData := models.Oder_item{
		UserIdNo:    uint(userID),
		TotalAmount: uint(totalprice),
		PaymentId:   pid,
		AddId:       addressData.Addressid,
		OrderStatus: "success",
		PaymentM:    method,
	}

	result = DB.Create(&oderData)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":    true,
		"paymentid": pid,
	})
	OrderDetails(c)
	DeleteCartItems(c)
}

// When the payment is successfull this function will work
func Success(c *gin.Context) {

	pid, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})
}
