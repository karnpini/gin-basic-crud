package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := setupRouter()
	r.Run(":3000")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	h := CustomerHandler{}
	h.Initialize()

	r.GET("/customers", h.GetAllCustomer)
	r.GET("/customers/:id", h.GetCustomer)
	r.POST("/customers", h.SaveCustomer)
	r.PUT("/customers/:id", h.UpdateCustomer)
	r.DELETE("/customers/:id", h.DeleteCustomer)

	return r
}

type CustomerHandler struct {
	DB *gorm.DB
}

func (h *CustomerHandler) Initialize() {
	db, err := gorm.Open(sqlite.Open("customer.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	db.AutoMigrate(&Customer{})
	h.DB = db
}

type Customer struct {
	Id        uint   `gorm:"primary_key" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

func (h *CustomerHandler) GetAllCustomer(c *gin.Context) {
	customer := []Customer{}

	h.DB.Find(&customer)

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) SaveCustomer(c *gin.Context) {
	customer := Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}
	println(id)
	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DB.Save(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := Customer{}

	if err := h.DB.Find(&customer, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := h.DB.Delete(&customer).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
