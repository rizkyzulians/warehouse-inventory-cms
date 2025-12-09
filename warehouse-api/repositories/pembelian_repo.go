package repositories

import (
	"database/sql"
	"fmt"
	"warehouse-api/models"
)

type PembelianRepository interface {
	CreateHeader(tx *sql.Tx, header *models.BeliHeader) error
	CreateDetail(tx *sql.Tx, detail *models.BeliDetail) error
	FindAll(limit, offset int) ([]models.BeliHeader, int, error)
	FindByID(id int) (*models.BeliHeaderWithDetail, error)
	GenerateNoFaktur(tanggal string) (string, error)
}

type pembelianRepository struct {
	db *sql.DB
}

func NewPembelianRepository(db *sql.DB) PembelianRepository {
	return &pembelianRepository{db: db}
}

func (r *pembelianRepository) CreateHeader(tx *sql.Tx, header *models.BeliHeader) error {
	query := `INSERT INTO beli_header (no_faktur, tanggal, supplier, total, keterangan, created_by)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	return tx.QueryRow(query, header.NoFaktur, header.Tanggal, header.Supplier,
		header.Total, header.Keterangan, header.CreatedBy).Scan(
		&header.ID, &header.CreatedAt, &header.UpdatedAt,
	)
}

func (r *pembelianRepository) CreateDetail(tx *sql.Tx, detail *models.BeliDetail) error {
	query := `INSERT INTO beli_detail (beli_header_id, barang_id, qty, harga, subtotal)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	return tx.QueryRow(query, detail.BeliHeaderID, detail.BarangID, detail.Qty,
		detail.Harga, detail.Subtotal).Scan(&detail.ID, &detail.CreatedAt)
}

func (r *pembelianRepository) FindAll(limit, offset int) ([]models.BeliHeader, int, error) {
	var headers []models.BeliHeader
	var total int

	// Count total
	countQuery := `SELECT COUNT(*) FROM beli_header`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := `SELECT id, no_faktur, tanggal, supplier, total, keterangan, 
	          created_by, created_at, updated_at
	          FROM beli_header ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.BeliHeader
		err := rows.Scan(&h.ID, &h.NoFaktur, &h.Tanggal, &h.Supplier, &h.Total,
			&h.Keterangan, &h.CreatedBy, &h.CreatedAt, &h.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		headers = append(headers, h)
	}

	return headers, total, nil
}

func (r *pembelianRepository) FindByID(id int) (*models.BeliHeaderWithDetail, error) {
	// Get header
	header := &models.BeliHeaderWithDetail{}
	queryHeader := `SELECT id, no_faktur, tanggal, supplier, total, keterangan,
	                created_by, created_at, updated_at
	                FROM beli_header WHERE id = $1`

	err := r.db.QueryRow(queryHeader, id).Scan(
		&header.ID, &header.NoFaktur, &header.Tanggal, &header.Supplier,
		&header.Total, &header.Keterangan, &header.CreatedBy,
		&header.CreatedAt, &header.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pembelian not found")
	}
	if err != nil {
		return nil, err
	}

	// Get details
	queryDetail := `SELECT d.id, d.beli_header_id, d.barang_id, d.qty, d.harga, 
	                d.subtotal, d.created_at, b.kode_barang, b.nama_barang, b.satuan
	                FROM beli_detail d
	                JOIN master_barang b ON d.barang_id = b.id
	                WHERE d.beli_header_id = $1`

	rows, err := r.db.Query(queryDetail, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []models.BeliDetailWithBarang
	for rows.Next() {
		var d models.BeliDetailWithBarang
		err := rows.Scan(&d.ID, &d.BeliHeaderID, &d.BarangID, &d.Qty, &d.Harga,
			&d.Subtotal, &d.CreatedAt, &d.KodeBarang, &d.NamaBarang, &d.Satuan)
		if err != nil {
			return nil, err
		}
		details = append(details, d)
	}

	header.Details = details
	return header, nil
}

func (r *pembelianRepository) GenerateNoFaktur(tanggal string) (string, error) {
	// Format: BL/YYYYMMDD/001
	// Extract date from tanggal (format: YYYY-MM-DD)
	datePrefix := tanggal[0:4] + tanggal[5:7] + tanggal[8:10]

	var lastNumber int
	query := `SELECT COALESCE(MAX(CAST(SUBSTRING(no_faktur FROM LENGTH(no_faktur) - 2) AS INTEGER)), 0)
	          FROM beli_header 
	          WHERE no_faktur LIKE $1`

	pattern := fmt.Sprintf("BL/%s/%%", datePrefix)
	err := r.db.QueryRow(query, pattern).Scan(&lastNumber)
	if err != nil {
		return "", err
	}

	nextNumber := lastNumber + 1
	return fmt.Sprintf("BL/%s/%03d", datePrefix, nextNumber), nil
}
