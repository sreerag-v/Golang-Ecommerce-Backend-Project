package controllers

import (
	"fmt"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/models"
	"github.com/tealeg/xlsx"
)

func SalesReport(c *gin.Context) {

	// want to fetch the dates from the url
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// converting the date string to time.time
	fromtime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid start Date",
		})
		return
	}

	totime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid end Date",
		})
	}

	// wnat to fetch the dates from oder_details table
	var orderDetails []models.OderDetails
	DB := database.InitDB()

	result := DB.Preload("Product").Preload("Payment").
		Where("created_at BETWEEN ? AND ?", fromtime, totime).
		Find(&orderDetails)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	// create exel file
	ex := excelize.NewFile()

	// create a newsheet
	SheetName := "Sheet1"
	index := ex.NewSheet(SheetName)

	// set values of headers
	ex.SetCellValue(SheetName, "A1", "Order Date")
	ex.SetCellValue(SheetName, "B1", "Order ID")
	ex.SetCellValue(SheetName, "C1", "Product name")
	ex.SetCellValue(SheetName, "D1", "Price")
	ex.SetCellValue(SheetName, "E1", "Total Amount")
	ex.SetCellValue(SheetName, "F1", "Payment method")
	ex.SetCellValue(SheetName, "G1", "Payment Status")

	// Set the value of cell
	for i, report := range orderDetails {
		row := i + 2
		ex.SetCellValue(SheetName, fmt.Sprintf("A%d", row), report.CreatedAt.Format("01/02/2006"))
		ex.SetCellValue(SheetName, fmt.Sprintf("B%d", row), report.Oderid)
		ex.SetCellValue(SheetName, fmt.Sprintf("C%d", row), report.Product.ProductName)
		ex.SetCellValue(SheetName, fmt.Sprintf("D%d", row), report.Product.Price)
		ex.SetCellValue(SheetName, fmt.Sprintf("E%d", row), report.Payment.Totalamount)
		ex.SetCellValue(SheetName, fmt.Sprintf("F%d", row), report.Payment.PaymentMethod)
		ex.SetCellValue(SheetName, fmt.Sprintf("G%d", row), report.Payment.Status)
	}

	// Set active sheet of the workbook.
	ex.SetActiveSheet(index)

	if err := ex.SaveAs("./public/SalesReport.xlsx"); err != nil {
		fmt.Println(err)
	}

	CovertingExelToPdf(c)
	c.HTML(200, "salesReport.html", gin.H{})
}

func CovertingExelToPdf(c *gin.Context) {
	// Open the Excel file
	xlFile, err := xlsx.OpenFile("./public/SalesReport.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 14)

	// Convertig each cell in the Excel file to a PDF cell
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row.Cells {
				//if there is any empty cell values skiping that
				if cell.Value == "" {
					continue
				}

				pdf.Cell(40, 10, cell.Value)
			}
			pdf.Ln(-1)
		}
	}

	// Save the PDF document
	err = pdf.OutputFileAndClose("./public/SalesReport.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("PDF saved successfully.")
}

func DownloadExel(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=SalesReport.xlsx")
	c.Header("Content-Type", "application/xlsx")
	c.File("./public/SalesReport.xlsx")
}

func Downloadpdf(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=SalesReport.pdf")
	c.Header("Content-Type", "application/pdf")
	c.File("./public/SalesReport.pdf")
}
