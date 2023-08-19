package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/controllers"
	"github.com/sreerag_v/Ecom/middleware"
)

func AdminRoutes(ctx *gin.Engine) {
	admin := ctx.Group("/admin")
	{
		admin.POST("/signup", controllers.AdminSignup)
		admin.POST("/login", controllers.AdminLogin)
		admin.GET("/home", middleware.AdminAuth(), controllers.AdminHome)

		admin.POST("/add-brand", middleware.AdminAuth(), controllers.AddBrand)
		admin.GET("/view-brand", middleware.AdminAuth(), controllers.ViewBrand)
		admin.GET("/edit-brand", middleware.AdminAuth(), controllers.EditBrand)

		admin.GET("/userdata", middleware.AdminAuth(), controllers.UserData)
		admin.PUT("/userdata/block/:id", middleware.AdminAuth(), controllers.BlockUser)
		admin.PUT("/userdata/unblock/:id", middleware.AdminAuth(), controllers.UnBlockUser)

		admin.POST("/add-category", middleware.AdminAuth(), controllers.AddCategories)
		admin.GET("/getcategory", middleware.AdminAuth(), controllers.ListAllCategory)
		admin.PUT("/editcategory/:id", middleware.AdminAuth(), controllers.EditCategory)
		admin.DELETE("/deletecategory/:id", middleware.AdminAuth(), controllers.DeleteCategory)

		admin.POST("/addproduct", middleware.AdminAuth(), controllers.ProductAdding)
		admin.GET("/view-products", middleware.AdminAuth(), controllers.ProductView)
		admin.GET("/view-products/:id", middleware.AdminAuth(), controllers.GetProductByID)
		admin.PUT("/editproducts/:id", middleware.AdminAuth(), controllers.EditProducts)
		admin.DELETE("/deleteproducts/:id", middleware.AdminAuth(), controllers.DeletePrdouct)

		admin.POST("/coupon/add", middleware.AdminAuth(), controllers.AddCoupon)
		admin.POST("/coupon/checkcoupon", middleware.AdminAuth(), controllers.CheckCoupon)

		admin.GET("user-orders", middleware.AdminAuth(), controllers.ShowAllOrders)
		admin.GET("user-orders/:id", middleware.AdminAuth(), controllers.ShowOrderById)

		// sales Report
		admin.GET("/order/salesreport", controllers.SalesReport)
		admin.GET("/order/salesreport/download/excel", controllers.DownloadExel)
		admin.GET("/order/salesreport/download/pdf", controllers.Downloadpdf)
	}
}
