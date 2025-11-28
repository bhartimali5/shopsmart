package routes

import (
	"fmt"
	"net/http"
	"reflect"

	"example.com/rest-api/dto"
	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
)

// GetUserOrders retrieves all orders for a specific user
// @Summary Get User Orders
// @Description Retrieve all orders for a specific user
// @Tags Orders
// @Security BearerAuth
// @Success 200 {object} dto.GetOrdersResponseDTO
// @Router /orders [get]
func GetCurrentUserOrders(c *gin.Context) {
	userId := c.GetString("user_id")
	Orders, err := models.GetOrdersByUserID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve orders"})
		return
	}
	if len(Orders) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No orders found for user"})
		return
	}
	response := dto.GetOrdersResponseDTO{}
	orders := []dto.OrderResponseDTO{}
	for _, order := range Orders {
		orderDTO := dto.OrderResponseDTO{
			ID:          order.ID,
			CartID:      order.CartID,
			OrderDate:   order.OrderDate,
			Status:      order.Status,
			TotalAmount: order.TotalAmount,
			UserID:      order.UserID,
		}
		orders = append(orders, orderDTO)
	}
	response.Orders = orders
	c.JSON(http.StatusOK, response)
}

// CreateUserOrder creates a new order for a specific user
// @Summary Create User Order
// @Description Create a new order for a specific user
// @Tags Orders
// @Security BearerAuth
// @Success 201 {object} dto.CreateOrderResponseDTO
// @Router /orders [post]
func CreateUserOrder(c *gin.Context) {
	userId := c.GetString("user_id")

	// Fetch cart details from Cart Service
	cartDetails := utils.GetCartItemDetails(userId, c.GetHeader("Authorization"))
	fmt.Println(cartDetails)
	if reflect.DeepEqual(cartDetails, reflect.Zero(reflect.TypeOf(cartDetails)).Interface()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not retrieve cart for user"})
		return
	}

	// Get current date as order date
	orderDate, err := utils.GetCurrentTime()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get current time"})
		return
	}

	order := models.Order{
		UserID:      userId,
		OrderDate:   orderDate,
		Status:      "PENDING_PAYMENT",
		TotalAmount: cartDetails.Cart.TotalPrice,
		CartID:      cartDetails.Cart.ID,
	}
	if err := order.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create order"})
		return
	}
	c.JSON(http.StatusCreated, dto.CreateOrderResponseDTO{
		UserID:      order.UserID,
		OrderDate:   order.OrderDate,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
	})
}
