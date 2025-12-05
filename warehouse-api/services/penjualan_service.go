package services

import (
	"database/sql"
	"fmt"
	"warehouse-api/models"
	"warehouse-api/repositories"
)

type PenjualanService interface {
	CreatePenjualan(req *models.CreatePenjualanRequest, userID int) (*models.JualHeaderWithDetail, error)
	GetAllPenjualan(limit, offset int) ([]models.JualHeader, int, error)
	GetPenjualanByID(id int) (*models.JualHeaderWithDetail, error)
}

type penjualanService struct {
	db            *sql.DB
	penjualanRepo repositories.PenjualanRepository
	barangRepo    repositories.BarangRepository
	stokRepo      repositories.StokRepository
}

func NewPenjualanService(db *sql.DB, penjualanRepo repositories.PenjualanRepository,
	barangRepo repositories.BarangRepository, stokRepo repositories.StokRepository) PenjualanService {
	return &penjualanService{
		db:            db,
		penjualanRepo: penjualanRepo,
		barangRepo:    barangRepo,
		stokRepo:      stokRepo,
	}
}

func (s *penjualanService) CreatePenjualan(req *models.CreatePenjualanRequest, userID int) (*models.JualHeaderWithDetail, error) {
	// Validate details
	if len(req.Details) == 0 {
		return nil, fmt.Errorf("details cannot be empty")
	}

	// Validate stock availability and calculate total
	var total float64
	for _, detail := range req.Details {
		// Validate barang exists
		_, err := s.barangRepo.FindByID(detail.BarangID)
		if err != nil {
			return nil, fmt.Errorf("barang with id %d not found", detail.BarangID)
		}

		// Check stock
		currentStok, err := s.stokRepo.FindByBarangID(detail.BarangID)
		if err != nil {
			return nil, err
		}

		if currentStok == nil || currentStok.StokAkhir < detail.Qty {
			return nil, &InsufficientStockError{
				BarangID:     detail.BarangID,
				RequestedQty: detail.Qty,
				AvailableQty: 0,
			}
		}

		if currentStok.StokAkhir < detail.Qty {
			return nil, &InsufficientStockError{
				BarangID:     detail.BarangID,
				RequestedQty: detail.Qty,
				AvailableQty: currentStok.StokAkhir,
			}
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
	header := &models.JualHeader{
		NoFaktur:   req.NoFaktur,
		Tanggal:    req.Tanggal,
		Customer:   req.Customer,
		Total:      total,
		Keterangan: req.Keterangan,
		CreatedBy:  userID,
	}

	if err := s.penjualanRepo.CreateHeader(tx, header); err != nil {
		return nil, err
	}

	// Process each detail
	var details []models.JualDetailWithBarang
	for _, detailReq := range req.Details {
		// Create detail
		detail := &models.JualDetail{
			JualHeaderID: header.ID,
			BarangID:     detailReq.BarangID,
			Qty:          detailReq.Qty,
			Harga:        detailReq.Harga,
			Subtotal:     float64(detailReq.Qty) * detailReq.Harga,
		}

		if err := s.penjualanRepo.CreateDetail(tx, detail); err != nil {
			return nil, err
		}

		// Get current stock
		currentStok, err := s.stokRepo.FindByBarangID(detailReq.BarangID)
		if err != nil {
			return nil, err
		}

		// Update stock (reduce qty)
		if err := s.stokRepo.UpdateStok(tx, detailReq.BarangID, 0, detailReq.Qty); err != nil {
			return nil, err
		}

		// Insert history
		history := &models.HistoryStok{
			BarangID:       detailReq.BarangID,
			JenisTransaksi: "keluar",
			Qty:            detailReq.Qty,
			StokSebelum:    currentStok.StokAkhir,
			StokSesudah:    currentStok.StokAkhir - detailReq.Qty,
			Keterangan:     fmt.Sprintf("Penjualan - %s", req.NoFaktur),
			ReferensiID:    &header.ID,
			ReferensiTipe:  "penjualan",
		}

		if err := s.stokRepo.InsertHistory(tx, history); err != nil {
			return nil, err
		}

		// Get barang info for response
		barang, _ := s.barangRepo.FindByID(detailReq.BarangID)
		detailWithBarang := models.JualDetailWithBarang{
			JualDetail: *detail,
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

	result := &models.JualHeaderWithDetail{
		JualHeader: *header,
		Details:    details,
	}

	return result, nil
}

func (s *penjualanService) GetAllPenjualan(limit, offset int) ([]models.JualHeader, int, error) {
	return s.penjualanRepo.FindAll(limit, offset)
}

func (s *penjualanService) GetPenjualanByID(id int) (*models.JualHeaderWithDetail, error) {
	return s.penjualanRepo.FindByID(id)
}

// Custom error for insufficient stock
type InsufficientStockError struct {
	BarangID     int
	RequestedQty int
	AvailableQty int
}

func (e *InsufficientStockError) Error() string {
	return fmt.Sprintf("insufficient stock for barang_id %d: requested %d, available %d",
		e.BarangID, e.RequestedQty, e.AvailableQty)
}
