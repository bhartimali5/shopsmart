package routes

import (
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
	if reflect.DeepEqual(cartDetails, reflect.Zero(reflect.TypeOf(cartDetails)).Interface()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty or could not be retrieved"})
		return
	}

	// Get current date as order date
	orderDate, err := utils.GetCurrentTime()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get current time"})
		return
	}

	// check if the userid in cartDetails matches the userId from the token
	if cartDetails.Cart.UserID != userId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access to cart"})
		return
	}
	order := models.Order{
		UserID:      userId,
		OrderDate:   orderDate,
		Status:      "PENDING_PAYMENT",
		TotalAmount: cartDetails.Cart.TotalPrice,
		CartID:      cartDetails.Items[0].CartId,
	}
	// Save order using transaction
	tx, err := order.SaveTx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Clear user's cart after order creation
	err = utils.ClearUserCart(userId, c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not clear user cart"})
		// Rollback: the created order if cart clearing fails
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not finalize order creation"})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateOrderResponseDTO{
		UserID:      order.UserID,
		OrderDate:   order.OrderDate,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
	})
}

// UpdateOrderStatus updates the status of an existing order
// @Summary Update Order Status
// @Description Update the status of an existing order
// @Tags Orders
// @Security BearerAuth
// @Param order_id path string true "Order ID"
// @Param status body dto.UpdateOrderStatusRequestDTO true "New Status"
// @Router /orders/{order_id}/status [patch]
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("order_id")
	var statusUpdate dto.UpdateOrderStatusRequestDTO
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	order, err := models.GetOrderByID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	order.Status = *statusUpdate.Status
	err = order.UpdateStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update order status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})

}
