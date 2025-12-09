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

type PembelianHandler struct {
	pembelianService services.PembelianService
}

func NewPembelianHandler(pembelianService services.PembelianService) *PembelianHandler {
	return &PembelianHandler{pembelianService: pembelianService}
}

func (h *PembelianHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePembelianRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate input - no faktur sudah auto-generate
	if req.Tanggal == "" || req.Supplier == "" {
		SendErrorResponse(w, http.StatusUnprocessableEntity, "Tanggal and supplier are required", "")
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

	result, err := h.pembelianService.CreatePembelian(&req, claims.UserID)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to create pembelian", err.Error())
		return
	}

	SendSuccessResponse(w, http.StatusCreated, "Pembelian created successfully", result, nil)
}

func (h *PembelianHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	pembelians, total, err := h.pembelianService.GetAllPembelian(limit, offset)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get pembelian", err.Error())
		return
	}

	meta := &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	SendSuccessResponse(w, http.StatusOK, "Pembelian retrieved successfully", pembelians, meta)
}

func (h *PembelianHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	pembelian, err := h.pembelianService.GetPembelianByID(id)
	if err != nil {
		if err.Error() == "pembelian not found" {
			SendErrorResponse(w, http.StatusNotFound, "Pembelian not found", "")
			return
		}
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get pembelian", err.Error())
		return
	}

	SendSuccessResponse(w, http.StatusOK, "Pembelian retrieved successfully", pembelian, nil)
}
