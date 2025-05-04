package controller

import (
	"event-ticketing/dto"
	"event-ticketing/entity"
	"event-ticketing/service"
	"event-ticketing/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entity.User true "User registration info"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/register [post]
func (ctrl *authController) Register(c *gin.Context) {
	var log = utils.Log
	var user entity.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Errorf("Failed to bind JSON: %v", err)
		utils.BadRequestResponse(c, "Invalid request body", err.Error())
		return
	}

	if err := utils.ValidateStruct(user); err != nil {
		log.Errorf("Validation error: %v", err)
		utils.BadRequestResponse(c, "Validation error", err.Error())
		return
	}

	if err := ctrl.authService.Register(&user); err != nil {
		log.Errorf("Registration failed: %v", err)
		utils.ConflictResponse(c, "Registration failed", err.Error())
		return
	}

	user.Password = ""

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.UserRequestDto true "User registration info"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login [post]
func (ctrl *authController) Login(c *gin.Context) {
	var log = utils.Log
	var request dto.UserRequestDto

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Errorf("Failed to bind JSON: %v", err)
		utils.BadRequestResponse(c, "Invalid request body", err.Error())
		return
	}

	token, err := ctrl.authService.Login(request.Email, request.Password)
	if err != nil {
		log.Errorf("Login failed: %v", err)
		utils.UnauthorizedResponse(c, "Invalid credentials")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{"token": token})
}

// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/profile [get]
func (ctrl *authController) GetProfile(c *gin.Context) {
	var log = utils.Log

	userID, exists := c.Get("userID")
	if !exists {
		log.Error("User not authenticated")
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	user, err := ctrl.authService.GetUserByID(userID.(string))
	if err != nil {
		log.Errorf("Failed to get user profile: %v", err)
		utils.InternalServerErrorResponse(c, "Failed to get user profile", err.Error())
		return
	}

	user.Password = ""

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}
