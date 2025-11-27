package routes

import (
	"net/http"

	"example.com/rest-api/dto"
	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create new user
// @Description This endpoint registers a user
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UserSignupDTO true "User Data"
// @Success 200 {object} dto.UserSignupResponseDTO
// @Router /signup [post]
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
	// address rec should also be created
	var newAddress models.Address

	newAddress.Name = ""
	newAddress.Address = ""
	newAddress.UserID = newUser.ID

	err = newAddress.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	context.JSON(http.StatusCreated, dto.UserSignupResponseDTO{
		ID:    newUser.ID,
		Email: newUser.Email,
	})
}

// LoginUser godoc
// @Summary Login user
// @Description This endpoint logs in a user
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UserLoginDTO true "Login Data"
// @Success 200 {object} dto.UserLoginResponseDTO
// @Router /login [post]
func login(context *gin.Context) {
	var loginData models.User
	err := context.ShouldBindJSON(&loginData)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_role, err := loginData.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateJWT(loginData.Email, loginData.ID, *user_role)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description This endpoint retrieves the profile of the logged-in user
// @Tags User
// @Produce json
// @Success 200 {object} dto.GetUserProfileResponseDTO
// @Router /user/profile [get]
// @Security BearerAuth
func getUserProfile(context *gin.Context) {
	userID := context.GetString("user_id")
	user, err := models.GetUserByID(userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	address, err := models.GetAddressByUserID(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user address"})
		return
	}

	profileResponse := dto.GetUserProfileResponseDTO{
		Email:   user.Email,
		Address: address.Address,
		Name:    address.Name,
		Role:    user.Role,
	}
	context.JSON(http.StatusOK, profileResponse)

}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description This endpoint updates the profile of the logged-in user
// @Tags User
// @Accept json
// @Produce json
// @Param user body dto.UpdateUserProfileDTO true "Update User Data"
// @Success 200 {object} dto.GetUserProfileResponseDTO
// @Router /user/profile [patch]
// @Security BearerAuth
func UpdateUserProfile(context *gin.Context) {
	userID := context.GetString("user_id")
	//userRole, _ := context.Get("user_role")
	var updateData dto.UpdateUserProfileDTO
	if err := context.ShouldBindJSON(&updateData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := models.GetUserByID(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user"})
		return
	}

	if updateData.Email != nil {
		user.Email = *updateData.Email
	}
	err = user.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}
	user_address, err := models.GetAddressByUserID(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user address"})
		return
	}
	if user_address != nil {
		if updateData.Address != nil {
			user_address.Address = *updateData.Address
		}
		if updateData.Name != nil {
			user_address.Name = *updateData.Name
		}
		err = user_address.Update()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update address"})
			return
		}
	} else {
		// Create new address if it doesn't exist
		var newAddress models.Address
		if updateData.Name != nil {
			newAddress.Name = *updateData.Name
		} else {
			newAddress.Name = ""
		}
		if updateData.Address != nil {
			newAddress.Address = *updateData.Address
		} else {
			newAddress.Address = ""
		}
		newAddress.UserID = userID
		err = newAddress.Save()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create address"})
			return
		}
	}

	profileResponse := dto.GetUserProfileResponseDTO{
		Email:   user.Email,
		Address: user_address.Address,
		Name:    user_address.Name,
		Role:    user.Role,
	}
	context.JSON(http.StatusOK, profileResponse)
}
