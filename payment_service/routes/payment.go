package routes

import (
	"net/http"

	"example.com/rest-api/dto"
	"example.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

// ListUserPayments godoc
// @Summary      Get all Payments for user
// @Description  get Payments
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.PaymentListResponseDTO
// @Router       /payments [get]
// @Security BearerAuth
func listUserPayments(c *gin.Context) {
	user_id := c.GetString("user_id")
	Payments, error := models.GetPaymentsForUserId(user_id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
	}
	var response dto.PaymentListResponseDTO
	var payments []dto.PaymentResponseDTO
	for _, payment := range Payments {
		payments = append(payments, dto.PaymentResponseDTO{
			OrderId:       payment.OrderId,
			UserId:        payment.UserId,
			CartId:        payment.CartId,
			PaymentStatus: payment.PaymentStatus,
			CreatedAt:     payment.CreatedAt,
			UpdatedAt:     payment.UpdatedAt,
		})
	}
	response.Payments = payments
	c.JSON(http.StatusOK, response)
}
