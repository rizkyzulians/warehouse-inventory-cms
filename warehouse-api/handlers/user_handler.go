package handlers

import (
	"encoding/json"
	"net/http"
	"warehouse-api/middleware"
	"warehouse-api/models"
	"warehouse-api/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		SendErrorResponse(w, http.StatusUnprocessableEntity, "Username and password are required", "")
		return
	}

	// Find user
	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	if user == nil {
		SendErrorResponse(w, http.StatusUnauthorized, "Invalid username or password", "")
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		SendErrorResponse(w, http.StatusUnauthorized, "Invalid username or password", "")
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  *user,
	}

	SendSuccessResponse(w, http.StatusOK, "Login successful", response, nil)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	user, err := h.userRepo.FindByID(claims.UserID)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	if user == nil {
		SendErrorResponse(w, http.StatusNotFound, "User not found", "")
		return
	}

	SendSuccessResponse(w, http.StatusOK, "Profile retrieved successfully", user, nil)
}
