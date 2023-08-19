package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
)

func ListAllCategory(c *gin.Context) {
	var categorys models.Category

	if categorysearch := c.Query("categorysearch"); categorysearch != "" {
		category := database.InitDB().Raw("SElECT * FROM categories WHERE Category=?", categorysearch).Scan(&categorys)
		if category.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": category.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{

		"available categories": categorys,
	})
}

type EditCategoryData struct {
	Category string
}

func EditCategory(c *gin.Context) {
	param := c.Param("id")

	var EditCategory EditCategoryData

	if err := c.ShouldBindJSON(&EditCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		c.Abort()
		return
	}

	var category models.Category
	record := database.InitDB().Model(category).Where("category_id = ?", param).Updates(models.Category{Category: EditCategory.Category})

	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "updated_successfully"})
}

func DeleteCategory(c *gin.Context) {
	param := c.Param("id")

	var category models.Category

	var count uint

	database.InitDB().Raw("select count(category_id) from categories where category_id=?", param).Scan(&count)
	if count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "product does not exist",
		})
		c.Abort()
		return
	}

	record := database.InitDB().Raw("delete from categories where category_id=?", param).Scan(&category)
	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "deleted successfully"})
}

func ProductAdding(c *gin.Context) {

	prodname := c.Request.FormValue("productname")
	price := c.Request.FormValue("price")
	Price, _ := strconv.Atoi(price)
	description := c.Request.FormValue("description")
	color := c.Request.FormValue("color")
	brand := c.Request.FormValue("brandID")
	brands, _ := strconv.Atoi(brand)
	stock := c.Request.FormValue("stock")
	Stock, _ := strconv.Atoi(stock)
	catogory := c.Request.FormValue("categoryID")
	catogoryy, _ := strconv.Atoi(catogory)

	imagepath, _ := c.FormFile("image")
	extension := filepath.Ext(imagepath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagepath, "./public/images"+image)

	var count uint
	database.Db.Raw("select count(*) from products where product_name=?", prodname).Scan(&count)
	fmt.Println(count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "A product with same name already exists",
		})
		c.Abort()
		return
	}

	product := models.Product{
		ProductName: prodname,
		Price:       uint(Price),
		Color:       color,
		Description: description,

		BrandId:    uint(brands),
		CatogeryId: uint(catogoryy),
		Image:      image,
		Stock:      uint(Stock),
	}

	if Stock < 0 {
		c.JSON(404, gin.H{
			"msg": " Stoke value is a negative value",
		})
		return
	} else if Stock >= 0 {
		record := database.Db.Create(&product)
		if record.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "product already exists",
			})
			c.Abort()
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Products added_successfully",
	})
}

var Products []struct {
	ProductID   uint
	ProductName string
	Price       string
	Description string
	Color       string
	Brands      string
	Stock       uint
	Category    string
	Image       string
}

func ProductView(c *gin.Context) {
	Pname := c.Query("Product_Name")
	Sort := c.Query("Sort")
	var count int64
	DB := database.InitDB()

	record := DB.Table("products").
		Joins("JOIN brands ON products.brand_id = brands.id").
		Joins("JOIN categories ON products.catogery_id = categories.id").
		Order("price " + Sort).
		Select("product_id,product_name,price,color,description,stock,image,brands.brands,categories.category").
		Scan(&Products)

	if record.Error != nil {
		c.JSON(401, gin.H{
			"msg": record.Error.Error(),
		})
		return
	}
	var flag bool
	if Sort == "" && Pname == "" {
		count, err1 := strconv.Atoi(c.Query("count"))
		PageN, err2 := strconv.Atoi(c.Query("page"))
		err3 := errors.Join(err1, err2)

		if err3 != nil {
			c.JSON(400, gin.H{
				"Error": "Error in string conversion",
			})
			return
		}

		limit := count
		offset := (PageN - 1) * limit
		normal := DB.Table("products").
			Joins("JOIN brands ON products.brand_id = brands.id").
			Joins("JOIN categories ON products.catogery_id = categories.id").
			Order("price ASC ").
			Select("product_id,product_name,price,color,description,stock,image,brands.brands,categories.category").
			Limit(limit).
			Offset(offset).
			Scan(&Products)

		if normal.Error != nil {
			c.JSON(401, gin.H{
				"msg": record.Error.Error(),
			})
			return
		}

		flag = true
	}

	if Sort == "" {
		record := DB.Table("products").
			Joins("JOIN brands ON products.brand_id = brands.id").
			Joins("JOIN categories ON products.catogery_id = categories.id").
			Where("product_name = ?", Pname).
			Select("product_id,product_name,price,color,description,stock,image,brands.brands,categories.category").
			Scan(&Products).Count(&count)

		if record.Error != nil {
			c.JSON(401, gin.H{
				"msg": "An error occurred while fetching product",
			})
			return
		}

		if count == 0 && !flag {
			c.JSON(http.StatusOK, gin.H{
				"msg": "No product found with the given name",
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"products": Products,
	})
}

func GetProductByID(c *gin.Context) {
	params := c.Param("id")
	// var product models.Product
	record := database.InitDB().Raw("SELECT product_id,product_name,price,color,stock,brands.brands FROM products join brands on products.brand_id = brands.id where product_id=?", params).Scan(&Products)
	if record.Error != nil {
		c.JSON(404, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"product": Products})
}

type EditProductsData struct {
	ProductName string `json:"productName"`
	Price       uint   `json:"price"`
	Brand       string `json:"brand"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

func EditProducts(c *gin.Context) {
	param := c.Param("id")

	var EditProducts EditProductsData
	if err := c.ShouldBindJSON(&EditProducts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		c.Abort()
		return
	}

	var products models.Product
	record := database.InitDB().Model(products).Where("product_id=?", param).Updates(models.Product{ProductName: EditProducts.ProductName,
		Price: EditProducts.Price, Color: EditProducts.Color, Description: EditProducts.Description})
	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "updated_successfully"})
}

func DeletePrdouct(c *gin.Context) {
	Param := c.Param("id")
	var products models.Product
	var count uint
	database.InitDB().Raw("select count(product_id) from products where product_id=?", Param).Scan(&count)
	if count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "product does not exist",
		})
		c.Abort()
		return
	}
	record := database.InitDB().Raw("delete from products where product_id=?", Param).Scan(&products)
	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "deleted successfully"})
}
