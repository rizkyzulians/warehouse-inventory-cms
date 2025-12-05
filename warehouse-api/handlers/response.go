package handlers

import (
	"encoding/json"
	"net/http"
	"warehouse-api/models"
)

// SendSuccessResponse sends a standardized success response
func SendSuccessResponse(w http.ResponseWriter, status int, message string, data interface{}, meta *models.Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := models.Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}

	json.NewEncoder(w).Encode(response)
}

// SendErrorResponse sends a standardized error response
func SendErrorResponse(w http.ResponseWriter, status int, message string, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := models.ErrorResponse{
		Success: false,
		Message: message,
		Error:   error,
	}

	json.NewEncoder(w).Encode(response)
}

// SendErrorResponseWithCode sends a standardized error response with error code
func SendErrorResponseWithCode(w http.ResponseWriter, status int, message string, error string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := models.ErrorResponse{
		Success: false,
		Message: message,
		Error:   error,
		Code:    code,
	}

	json.NewEncoder(w).Encode(response)
}
