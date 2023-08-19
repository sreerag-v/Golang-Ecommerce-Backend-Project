package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/controllers"
	"github.com/sreerag_v/Ecom/middleware"
)

func UserRoutes(ctx *gin.Engine) {

	user := ctx.Group("/user")

	{
		user.POST("/signup", controllers.Signup)
		user.POST("/signup/otp", controllers.OtpValidation)

		user.POST("/login", controllers.LoginUser)
		user.GET("/logout", middleware.UserAuth(), controllers.LogoutUser)

		user.GET("/view-products", middleware.UserAuth(), controllers.ProductView)
		user.GET("/view-products/:id", middleware.UserAuth(), controllers.GetProductByID)

		user.GET("/viewprofile", middleware.UserAuth(), controllers.ShowUserDetails)
		user.PUT("/Editprofile", middleware.UserAuth(), controllers.EditUserProfile)
		user.POST("/addaddress", middleware.UserAuth(), controllers.AddAddress)
		user.GET("/searchaddress/:id", middleware.UserAuth(), controllers.ShowAddress)

		user.POST("/add-wishlist", middleware.UserAuth(), controllers.Wishlist)

		user.POST("/forgotpassword", middleware.UserAuth(), controllers.GenrateOtpForForgotPassword)
		user.PUT("/changepassword", middleware.UserAuth(), controllers.ChangePassword)

		user.POST("/addtocart", middleware.UserAuth(), controllers.AddToCart)
		user.GET("/view-cart", middleware.UserAuth(), controllers.ViewCart)
		user.DELETE("/deletecart/:id", middleware.UserAuth(), controllers.DeleteCart)
		user.GET("/checkout", middleware.UserAuth(), controllers.CheckOut)

		user.GET("/payment/razorpay", controllers.Razorpay)
		user.GET("/payment/success", middleware.UserAuth(), controllers.RazorpaySuccess)
		user.GET("/success", middleware.UserAuth(), controllers.Success)

		user.POST("/apply-coupen", middleware.UserAuth(), controllers.ApplyCoupen)

		user.GET("/cod", middleware.UserAuth(), controllers.CashOnDelivery)
		user.GET("/view-order", middleware.UserAuth(), controllers.ShowOrder)
		user.GET("/cancel-order", middleware.UserAuth(), controllers.CancelOrder)
		user.GET("/return-order", middleware.UserAuth(), controllers.ReturnOrder)

	}
}
