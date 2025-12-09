package repositories

import (
	"database/sql"
	"fmt"
	"warehouse-api/models"
)

type PenjualanRepository interface {
	CreateHeader(tx *sql.Tx, header *models.JualHeader) error
	CreateDetail(tx *sql.Tx, detail *models.JualDetail) error
	FindAll(limit, offset int) ([]models.JualHeader, int, error)
	FindByID(id int) (*models.JualHeaderWithDetail, error)
	GenerateNoFaktur(tanggal string) (string, error)
}

type penjualanRepository struct {
	db *sql.DB
}

func NewPenjualanRepository(db *sql.DB) PenjualanRepository {
	return &penjualanRepository{db: db}
}

func (r *penjualanRepository) CreateHeader(tx *sql.Tx, header *models.JualHeader) error {
	query := `INSERT INTO jual_header (no_faktur, tanggal, customer, total, keterangan, created_by)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	return tx.QueryRow(query, header.NoFaktur, header.Tanggal, header.Customer,
		header.Total, header.Keterangan, header.CreatedBy).Scan(
		&header.ID, &header.CreatedAt, &header.UpdatedAt,
	)
}

func (r *penjualanRepository) CreateDetail(tx *sql.Tx, detail *models.JualDetail) error {
	query := `INSERT INTO jual_detail (jual_header_id, barang_id, qty, harga, subtotal)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	return tx.QueryRow(query, detail.JualHeaderID, detail.BarangID, detail.Qty,
		detail.Harga, detail.Subtotal).Scan(&detail.ID, &detail.CreatedAt)
}

func (r *penjualanRepository) FindAll(limit, offset int) ([]models.JualHeader, int, error) {
	var headers []models.JualHeader
	var total int

	// Count total
	countQuery := `SELECT COUNT(*) FROM jual_header`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := `SELECT id, no_faktur, tanggal, customer, total, keterangan,
	          created_by, created_at, updated_at
	          FROM jual_header ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.JualHeader
		err := rows.Scan(&h.ID, &h.NoFaktur, &h.Tanggal, &h.Customer, &h.Total,
			&h.Keterangan, &h.CreatedBy, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		headers = append(headers, h)
	}

	return headers, total, nil
}

func (r *penjualanRepository) FindByID(id int) (*models.JualHeaderWithDetail, error) {
	// Get header
	header := &models.JualHeaderWithDetail{}
	queryHeader := `SELECT id, no_faktur, tanggal, customer, total, keterangan,
	                created_by, created_at, updated_at
	                FROM jual_header WHERE id = $1`

	err := r.db.QueryRow(queryHeader, id).Scan(
		&header.ID, &header.NoFaktur, &header.Tanggal, &header.Customer,
		&header.Total, &header.Keterangan, &header.CreatedBy,
		&header.CreatedAt, &header.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("penjualan not found")
	}
	if err != nil {
		return nil, err
	}

	// Get details
	queryDetail := `SELECT d.id, d.jual_header_id, d.barang_id, d.qty, d.harga,
	                d.subtotal, d.created_at, b.kode_barang, b.nama_barang, b.satuan
	                FROM jual_detail d
	                JOIN master_barang b ON d.barang_id = b.id
	                WHERE d.jual_header_id = $1`

	rows, err := r.db.Query(queryDetail, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []models.JualDetailWithBarang
	for rows.Next() {
		var d models.JualDetailWithBarang
		err := rows.Scan(&d.ID, &d.JualHeaderID, &d.BarangID, &d.Qty, &d.Harga,
			&d.Subtotal, &d.CreatedAt, &d.KodeBarang, &d.NamaBarang, &d.Satuan)
		if err != nil {
			return nil, err
		}
		details = append(details, d)
	}

	header.Details = details
	return header, nil
}

func (r *penjualanRepository) GenerateNoFaktur(tanggal string) (string, error) {
	// Format: JL/YYYYMMDD/001
	// Extract date from tanggal (format: YYYY-MM-DD)
	datePrefix := tanggal[0:4] + tanggal[5:7] + tanggal[8:10]

	var lastNumber int
	query := `SELECT COALESCE(MAX(CAST(SUBSTRING(no_faktur FROM LENGTH(no_faktur) - 2) AS INTEGER)), 0)
	          FROM jual_header 
	          WHERE no_faktur LIKE $1`

	pattern := fmt.Sprintf("JL/%s/%%", datePrefix)
	err := r.db.QueryRow(query, pattern).Scan(&lastNumber)
	if err != nil {
		return "", err
	}

	nextNumber := lastNumber + 1
	return fmt.Sprintf("JL/%s/%03d", datePrefix, nextNumber), nil
}
