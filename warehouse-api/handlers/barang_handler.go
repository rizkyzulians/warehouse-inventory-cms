package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"warehouse-api/models"
	"warehouse-api/repositories"

	"github.com/gorilla/mux"
)

type BarangHandler struct {
	barangRepo repositories.BarangRepository
}

func NewBarangHandler(barangRepo repositories.BarangRepository) *BarangHandler {
	return &BarangHandler{barangRepo: barangRepo}
}

func (h *BarangHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	barangs, total, err := h.barangRepo.FindAll(search, limit, offset)
	if err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get barang", err.Error())
		return
	}

	meta := &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	SendSuccessResponse(w, http.StatusOK, "Barang retrieved successfully", barangs, meta)
}

func (h *BarangHandler) GetAllWithStok(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	barangs, total, err := h.barangRepo.FindAllWithStok(search, limit, offset)
	if err != nil {
		log.Printf("Error in GetAllWithStok: %v", err)
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get barang", err.Error())
		return
	}

	meta := &models.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	SendSuccessResponse(w, http.StatusOK, "Barang with stock retrieved successfully", barangs, meta)
}

func (h *BarangHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	barang, err := h.barangRepo.FindByID(id)
	if err != nil {
		if err.Error() == "barang not found" {
			SendErrorResponse(w, http.StatusNotFound, "Barang not found", "")
			return
		}
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get barang", err.Error())
		return
	}

	SendSuccessResponse(w, http.StatusOK, "Barang retrieved successfully", barang, nil)
}

func (h *BarangHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBarangRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate input - kode barang sudah auto-generate
	if req.NamaBarang == "" || req.Satuan == "" {
		SendErrorResponse(w, http.StatusUnprocessableEntity, "Nama barang and satuan are required", "")
		return
	}

	barang := &models.Barang{
		KodeBarang: req.KodeBarang,
		NamaBarang: req.NamaBarang,
		Kategori:   req.Kategori,
		Satuan:     req.Satuan,
		HargaBeli:  req.HargaBeli,
		HargaJual:  req.HargaJual,
	}

	if err := h.barangRepo.Create(barang); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to create barang", err.Error())
		return
	}

	SendSuccessResponse(w, http.StatusCreated, "Barang created successfully", barang, nil)
}

func (h *BarangHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	var req models.UpdateBarangRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate input
	if req.NamaBarang == "" || req.Satuan == "" {
		SendErrorResponse(w, http.StatusUnprocessableEntity, "Nama barang and satuan are required", "")
		return
	}

	// Check if barang exists
	existing, err := h.barangRepo.FindByID(id)
	if err != nil {
		if err.Error() == "barang not found" {
			SendErrorResponse(w, http.StatusNotFound, "Barang not found", "")
			return
		}
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get barang", err.Error())
		return
	}

	barang := &models.Barang{
		ID:         id,
		KodeBarang: existing.KodeBarang,
		NamaBarang: req.NamaBarang,
		Kategori:   req.Kategori,
		Satuan:     req.Satuan,
		HargaBeli:  req.HargaBeli,
		HargaJual:  req.HargaJual,
	}

	if err := h.barangRepo.Update(barang); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to update barang", err.Error())
		return
	}

	SendSuccessResponse(w, http.StatusOK, "Barang updated successfully", barang, nil)
}

func (h *BarangHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid ID", err.Error())
		return
	}

	// Check if barang exists
	_, err = h.barangRepo.FindByID(id)
	if err != nil {
		if err.Error() == "barang not found" {
			SendErrorResponse(w, http.StatusNotFound, "Barang not found", "")
			return
		}
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to get barang", err.Error())
		return
	}

	if err := h.barangRepo.Delete(id); err != nil {
		SendErrorResponse(w, http.StatusInternalServerError, "Failed to delete barang", err.Error())
		return
	}

	SendSuccessResponse(w, http.StatusOK, "Barang deleted successfully", nil, nil)
}
