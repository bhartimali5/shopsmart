package routes

import (
	"fmt"
	"net/http"

	"example.com/rest-api/dto"
	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
)

// godoc
// @Summary Add Item to Cart
// @Description Adds an item to the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param item body dto.AddItemRequest true "Item to add"
// @Router /cart/items [post]
// @Security BearerAuth
func AddItemToCart(c *gin.Context) {
	userId := c.GetString("user_id")

	var newItem models.CartItem
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	product_id := newItem.ProductId
	// Fetch product price from product service
	productPrice, err := utils.FetchProductPrice(newItem.ProductId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch product price"})
		return
	}

	// Get or create active cart for user
	current_user_cart, _ := models.GetActiveCartByUserId(userId)

	if current_user_cart == (models.Cart{}) {
		new_cart, err := models.CreateCart(models.Cart{UserId: userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create cart"})
			return
		}
		current_user_cart = new_cart
		newItem.CartId = current_user_cart.ID
		newItem.Price = float64(newItem.Quantity) * productPrice
		addedItem, err := models.AddItemToCart(newItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item to cart"})
			return
		}
		response := dto.AddItemResponse{
			Message: "Item has been added to cart!",
			Item: dto.CartItem{
				ID:        addedItem.ID,
				ProductId: addedItem.ProductId,
				Quantity:  addedItem.Quantity,
				Price:     addedItem.Price,
			},
		}

		c.JSON(http.StatusOK, response)
		return
	}

	// If cart exists, check if item already in cart
	newItem.CartId = current_user_cart.ID
	cart_items := models.GetCartItemsByCartId(current_user_cart.ID)

	for _, item := range cart_items {
		if item.ProductId == product_id {
			item.Quantity += newItem.Quantity
			item.Price = float64(item.Quantity) * productPrice

			// Update the item quantity and price in the database
			err := models.UpdateCartItem(item.ID, item.Quantity, item.Price)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, dto.AddItemResponse{
				Message: "Item quantity has been updated in cart!",
				Item: dto.CartItem{
					ID:        item.ID,
					ProductId: item.ProductId,
					Quantity:  item.Quantity,
					Price:     item.Price,
				},
			})
			return
		}
	}

	//Total price for product for all quantity
	newItem.Price = float64(newItem.Quantity) * productPrice
	addedItem, err := models.AddItemToCart(newItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item to cart"})
		return
	}
	response := dto.AddItemResponse{
		Message: "Item has been added to cart!",
		Item: dto.CartItem{
			ID:        addedItem.ID,
			ProductId: addedItem.ProductId,
			Quantity:  addedItem.Quantity,
			Price:     addedItem.Price,
		},
	}

	c.JSON(http.StatusOK, response)
}

// godoc
// @Summary Remove Item from Cart
// @Description Removes an item from the user's cart
// @Tags Cart
// @Param itemId path string true "Item ID"
// @Router /cart/items/{itemId} [delete]
// @Security BearerAuth
func RemoveItemFromCart(c *gin.Context) {
	userId := c.GetString("user_id")
	current_user_cart, err := models.GetActiveCartByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cart is empty!"})
		return
	}
	if current_user_cart == (models.Cart{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	itemId := c.Param("itemId")
	err = models.RemoveItemFromCart(itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove item from cart"})
		return
	}
	c.JSON(http.StatusOK, dto.RemoveItemResponse{
		Message: "Item has been removed from cart!",
		ItemId:  itemId,
	})
}

// godoc
// @Summary View Cart
// @Description Retrieves the user's cart and its items
// @Tags Cart
// @Accept       json
// @Produce      json
// @Success      200  {array} dto.CartItem
// @Router /cart/items [get]
// @Security BearerAuth
func ViewCart(c *gin.Context) {
	userId := c.GetString("user_id")
	type CartData struct {
		UserId     string  `json:"user_id"`
		TotalPrice float64 `json:"total_price"`
	}
	current_user_cart, err := models.GetActiveCartByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"Message": "Cart is empty!"})
		return
	}
	if current_user_cart == (models.Cart{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	cartItems := models.GetCartItemsByCartId(current_user_cart.ID)
	var totalAmount float64
	for items := range cartItems {
		totalAmount += cartItems[items].Price
	}
	fmt.Println(current_user_cart, totalAmount)
	c.JSON(http.StatusOK, gin.H{"cart": CartData{
		UserId:     userId,
		TotalPrice: totalAmount,
	}, "items": cartItems})
}

// godoc
// @Summary Clear Cart
// @Description Removes all items from the user's cart
// @Tags Cart
// @Router /cart/clear [delete]
// @Security BearerAuth
func ClearCart(c *gin.Context) {
	userId := c.GetString("user_id")
	current_user_cart, err := models.GetActiveCartByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cart is empty!"})
		return
	}
	if current_user_cart == (models.Cart{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	err = models.RemoveAllItemsFromCart(current_user_cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not clear cart"})
		return
	}
	err = models.DeleteUserCart(current_user_cart.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not clear cart"})
		return
	}
	c.JSON(http.StatusOK, dto.ClearCartResponse{
		Message: "Cart has been cleared!",
	})
}

// godoc
// @Summary Update Cart Item Quantity
// @Description Updates the quantity of a specific item in the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param itemId path string true "Item ID"
// @Param item body dto.UpdateItemRequest true "Updated item quantity"
// @Router /cart/items/{itemId} [patch]
// @Security BearerAuth
func UpdateCartItemQuantity(c *gin.Context) {
	userId := c.GetString("user_id")
	cartItemId := c.Param("itemId")

	// Verify that the cart item belongs to the user's cart
	current_user_cart, err := models.GetActiveCartByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cart is empty!"})
		return
	}
	if current_user_cart == (models.Cart{}) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}
	var updateRequest dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	cartItem, err := models.GetCartItemById(cartItemId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}
	if cartItem.CartId != current_user_cart.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to modify this item"})
		return
	}

	// Fetch product price from product service
	productPrice, err := utils.FetchProductPrice(cartItem.ProductId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch product price"})
		return
	}

	newPrice := float64(updateRequest.Quantity) * productPrice

	err = models.UpdateCartItem(cartItemId, updateRequest.Quantity, newPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update item quantity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item quantity has been updated in cart!"})
}
