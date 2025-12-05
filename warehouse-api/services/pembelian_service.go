package services

import (
	"database/sql"
	"fmt"
	"warehouse-api/models"
	"warehouse-api/repositories"
)

type PembelianService interface {
	CreatePembelian(req *models.CreatePembelianRequest, userID int) (*models.BeliHeaderWithDetail, error)
	GetAllPembelian(limit, offset int) ([]models.BeliHeader, int, error)
	GetPembelianByID(id int) (*models.BeliHeaderWithDetail, error)
}

type pembelianService struct {
	db            *sql.DB
	pembelianRepo repositories.PembelianRepository
	barangRepo    repositories.BarangRepository
	stokRepo      repositories.StokRepository
}

func NewPembelianService(db *sql.DB, pembelianRepo repositories.PembelianRepository,
	barangRepo repositories.BarangRepository, stokRepo repositories.StokRepository) PembelianService {
	return &pembelianService{
		db:            db,
		pembelianRepo: pembelianRepo,
		barangRepo:    barangRepo,
		stokRepo:      stokRepo,
	}
}

func (s *pembelianService) CreatePembelian(req *models.CreatePembelianRequest, userID int) (*models.BeliHeaderWithDetail, error) {
	// Validate details
	if len(req.Details) == 0 {
		return nil, fmt.Errorf("details cannot be empty")
	}

	// Calculate total
	var total float64
	for _, detail := range req.Details {
		// Validate barang exists
		_, err := s.barangRepo.FindByID(detail.BarangID)
		if err != nil {
			return nil, fmt.Errorf("barang with id %d not found", detail.BarangID)
		}

		subtotal := float64(detail.Qty) * detail.Harga
		total += subtotal
	}

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create header
	header := &models.BeliHeader{
		NoFaktur:   req.NoFaktur,
		Tanggal:    req.Tanggal,
		Supplier:   req.Supplier,
		Total:      total,
		Keterangan: req.Keterangan,
		CreatedBy:  userID,
	}

	if err := s.pembelianRepo.CreateHeader(tx, header); err != nil {
		return nil, err
	}

	// Process each detail
	var details []models.BeliDetailWithBarang
	for _, detailReq := range req.Details {
		// Create detail
		detail := &models.BeliDetail{
			BeliHeaderID: header.ID,
			BarangID:     detailReq.BarangID,
			Qty:          detailReq.Qty,
			Harga:        detailReq.Harga,
			Subtotal:     float64(detailReq.Qty) * detailReq.Harga,
		}

		if err := s.pembelianRepo.CreateDetail(tx, detail); err != nil {
			return nil, err
		}

		// Get current stock
		currentStok, err := s.stokRepo.FindByBarangID(detailReq.BarangID)
		if err != nil {
			return nil, err
		}

		// If stock doesn't exist, create it
		if currentStok == nil {
			if err := s.stokRepo.CreateStok(tx, detailReq.BarangID); err != nil {
				return nil, err
			}
			currentStok = &models.Stok{
				BarangID:  detailReq.BarangID,
				StokAkhir: 0,
			}
		}

		// Update stock (add qty)
		if err := s.stokRepo.UpdateStok(tx, detailReq.BarangID, detailReq.Qty, 0); err != nil {
			return nil, err
		}

		// Insert history
		history := &models.HistoryStok{
			BarangID:       detailReq.BarangID,
			JenisTransaksi: "masuk",
			Qty:            detailReq.Qty,
			StokSebelum:    currentStok.StokAkhir,
			StokSesudah:    currentStok.StokAkhir + detailReq.Qty,
			Keterangan:     fmt.Sprintf("Pembelian - %s", req.NoFaktur),
			ReferensiID:    &header.ID,
			ReferensiTipe:  "pembelian",
		}

		if err := s.stokRepo.InsertHistory(tx, history); err != nil {
			return nil, err
		}

		// Get barang info for response
		barang, _ := s.barangRepo.FindByID(detailReq.BarangID)
		detailWithBarang := models.BeliDetailWithBarang{
			BeliDetail: *detail,
			KodeBarang: barang.KodeBarang,
			NamaBarang: barang.NamaBarang,
			Satuan:     barang.Satuan,
		}
		details = append(details, detailWithBarang)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	result := &models.BeliHeaderWithDetail{
		BeliHeader: *header,
		Details:    details,
	}

	return result, nil
}

func (s *pembelianService) GetAllPembelian(limit, offset int) ([]models.BeliHeader, int, error) {
	return s.pembelianRepo.FindAll(limit, offset)
}

func (s *pembelianService) GetPembelianByID(id int) (*models.BeliHeaderWithDetail, error) {
	return s.pembelianRepo.FindByID(id)
}
