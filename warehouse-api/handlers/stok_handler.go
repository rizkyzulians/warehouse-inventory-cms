package handlers

import (
	"net/http"
	"strconv"
	"warehouse-api/models"
	"warehouse-api/repositories"

	"github.com/gorilla/mux"
)

type StokHandler struct {
	stokRepo repositories.StokRepository
}

func NewStokHandler(stokRepo repositories.StokRepository) *StokHandler {
	return &StokHandler{stokRepo: stokRepo}
}

func (h *StokHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	stoks, err := h.stokRepo.FindAll()
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get stock", err.Error())
		return
	}

	SendSuccessResponse(w, http.StatusOK, "Stock retrieved successfully", stoks, nil)
}

func (h *StokHandler) GetByBarangID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	barangID, err := strconv.Atoi(vars["barang_id"])
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid barang ID", err.Error())
		return
	}

	stok, err := h.stokRepo.FindByBarangID(barangID)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get stock", err.Error())
		return
	}

	if stok == nil {
		SendErrorResponse(w, http.StatusNotFound, "Stock not found", "")
		return
	}

	SendSuccessResponse(w, http.StatusOK, "Stock retrieved successfully", stok, nil)
}

func (h *StokHandler) GetHistoryAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	histories, total, err := h.stokRepo.GetHistoryAll(limit, offset)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get stock history", err.Error())
		return
	}

	meta := &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	SendSuccessResponse(w, http.StatusOK, "Stock history retrieved successfully", histories, meta)
}

func (h *StokHandler) GetHistoryByBarangID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	barangID, err := strconv.Atoi(vars["barang_id"])
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid barang ID", err.Error())
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	histories, total, err := h.stokRepo.GetHistoryByBarangID(barangID, limit, offset)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get stock history", err.Error())
		return
	}

	meta := &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	SendSuccessResponse(w, http.StatusOK, "Stock history retrieved successfully", histories, meta)
}
