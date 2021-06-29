package controllers

import (
	"net/http"
	"strconv"

	"project/config"
	m "project/middlewares"
	"project/models"

	"github.com/labstack/echo"
)

/* POST /user --> to add user data
{
    "name": "Hulk",
    "email": "banner@stapp03",
    "password": "BigGr33n",
	"role": "superadmin"
}
*/
func CreateUsersController(c echo.Context) error {
	user := models.Users{}
	c.Bind(&user)
	// Validate
	if err, ok := user.Validate(); !ok {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "ValidationError", Status: err})
	}
	// Check for existing user
	var existing_email string
	if err := config.DB.Debug().Model(&user).Select("email").Where("email = ?", user.Email).Scan(&existing_email).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error Occured while registering user",
			Status:  err,
		})
	}
	// Prevent duplicate email
	if user.Email == existing_email {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "User Already Exists"})
	}
	// Store the credentials data to database
	if err := config.DB.Debug().Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Message: "Error Occured while registering user",
			Status:  err,
		})
	}
	user_response := models.UsersResponse{
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "User registered",
		Data:    user_response,
	})
}

// POST /login --> to login or create user token
func LoginUsersController(c echo.Context) error {
	// Variable to store data from database
	user := &models.Users{}
	// Bind with particular struct to validate
	login_user := &models.UserLoginCreds{}
	if err := c.Bind(login_user); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Cannot Login User", Status: err,
		})
	}
	// Validate
	if err, ok := login_user.Validate(); !ok {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "ValidationError", Status: err,
		})
	}
	// Find user
	if err := config.DB.Debug().Model(&models.Users{}).Where("email = ? AND password = ?", login_user.Email, login_user.Password).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error()})
	}
	// Create token
	accessToken, accessTokenExpiresAt, err := m.CreateToken(login_user.Email, login_user.Role, user.Id, false)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail create JWT Token",
			Data:    err.Error(),
		})
	}
	// Refresh token
	refreshToken, refreshTokenExpiresAt, err := m.CreateToken(login_user.Email, login_user.Role, user.Id, true)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: "Fail create JWT Refresh Token",
			Data:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message:               "Logged in",
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	})
}

// POST /refresh --> to refresh token when expires
// Must include the token in the header when making post request
func RefreshTokenController(c echo.Context) error {
	// Decode and get the email from token
	email, ok := m.GetFieldFromToken("email", c)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Cannot identify accesstoken"})
	}

	// Get user from database based on email
	user := models.Users{}
	result := config.DB.Debug().Model(models.Users{}).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: result.Error,
		})
	}
	// If user not available or bad token
	if result.RowsAffected < 1 {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "Cannot refresh token"})
	}
	// Give the user some access token
	accessToken, accessTokenExpiresAt, err := m.CreateToken(user.Email, user.Role, user.Id, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "Cannot refresh token", Status: err})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message:              "New access token created",
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
	})
}

// GET /users --> to get all users data
// Must have a superadmin level to access
func GetAllUsersController(c echo.Context) error {
	user := models.Users{}
	users := []models.Users{}
	// Check whether users data exist
	counted_users, err := user.CountAllUsers(config.DB)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	if counted_users == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "There's no user data yet in database",
			Status:  err.Error(),
		})
	}
	// With pagination implemented
	// Find all the available users from database
	if err := config.DB.Debug().Scopes(models.Paginate(c)).Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get all the users data successful",
		Data:    users,
	})
}

// GET /user/:id --> to get a user data specified by id
func GetUserController(c echo.Context) error {
	// Check the id parameter
	user_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status: err.Error(),
		})
	}
	// Retreive user object with the primary key (id)
	user := models.Users{}
	if err := config.DB.Debug().First(&user, user_id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Message: "Get a user data succesful",
		Data:    user,
	})
}
