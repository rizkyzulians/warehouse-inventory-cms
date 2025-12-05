package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"warehouse-api/middleware"
	"warehouse-api/models"
	"warehouse-api/services"

	"github.com/gorilla/mux"
)

type PenjualanHandler struct {
	penjualanService services.PenjualanService
}

func NewPenjualanHandler(penjualanService services.PenjualanService) *PenjualanHandler {
	return &PenjualanHandler{penjualanService: penjualanService}
}

func (h *PenjualanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePenjualanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate input
	if req.NoFaktur == "" || req.Tanggal == "" || req.Customer == "" {
		SendErrorResponse(w, http.StatusUnprocessableEntity, "No faktur, tanggal, and customer are required", "")
		return
	}

	if len(req.Details) == 0 {
		SendErrorResponse(w, http.StatusUnprocessableEntity, "Details cannot be empty", "")
		return
	}

	// Get user from context
	claims, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	result, err := h.penjualanService.CreatePenjualan(&req, claims.UserID)
	if err != nil {
		// Check if it's an insufficient stock error
		if insufficientErr, ok := err.(*services.InsufficientStockError); ok {
			SendErrorResponseWithCode(w, http.StatusBadRequest, "Insufficient stock", insufficientErr.Error(), "INSUFFICIENT_STOCK")
			return
		}
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to create penjualan", err.Error())
		return
	}

	SendSuccessResponse(w, http.StatusCreated, "Penjualan created successfully", result, nil)
}

func (h *PenjualanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	penjualans, total, err := h.penjualanService.GetAllPenjualan(limit, offset)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get penjualan", err.Error())
		return
	}

	meta := &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	SendSuccessResponse(w, http.StatusOK, "Penjualan retrieved successfully", penjualans, meta)
}

func (h *PenjualanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	penjualan, err := h.penjualanService.GetPenjualanByID(id)
	if err != nil {
		if err.Error() == "penjualan not found" {
			SendErrorResponse(w, http.StatusNotFound, "Penjualan not found", "")
			return
		}
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get penjualan", err.Error())
		return
	}

	SendSuccessResponse(w, http.StatusOK, "Penjualan retrieved successfully", penjualan, nil)
}
