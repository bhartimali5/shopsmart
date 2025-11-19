package routes

import (
	"net/http"

	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
)

func signUp(context *gin.Context) {
	var newUser models.User
	if err := context.ShouldBindJSON(&newUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := newUser.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "User created succesfully!", "user": newUser})
}

func login(context *gin.Context) {
	var loginData models.User
	err := context.ShouldBindJSON(&loginData)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = loginData.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	token, err := utils.GenerateJWT(loginData.Email, int64(loginData.ID))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func getUserProfile(context *gin.Context) {
	userID := context.GetInt64("user_id")
	user, err := models.GetUserByID(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user profile"})
		return
	}
	// Fetch addresses associated with the user
	address, err := models.GetAddressByUserID(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user profile"})
		return
	}
	if user == nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	profile := gin.H{
		"id":       user.ID,
		"username": user.UserName,
		"email":    user.Email,
		"address":  address,
	}
	context.JSON(http.StatusOK, profile)
}
