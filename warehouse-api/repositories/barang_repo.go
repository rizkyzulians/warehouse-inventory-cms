package repositories

import (
	"database/sql"
	"fmt"
	"warehouse-api/models"
)

type BarangRepository interface {
	FindAll(search string, limit, offset int) ([]models.Barang, int, error)
	FindByID(id int) (*models.Barang, error)
	FindAllWithStok(search string, limit, offset int) ([]models.BarangWithStok, int, error)
	Create(barang *models.Barang) error
	Update(barang *models.Barang) error
	Delete(id int) error
	GenerateKodeBarang() (string, error)
}

type barangRepository struct {
	db *sql.DB
}

func NewBarangRepository(db *sql.DB) BarangRepository {
	return &barangRepository{db: db}
}

func (r *barangRepository) FindAll(search string, limit, offset int) ([]models.Barang, int, error) {
	var barangs []models.Barang
	var total int

	// Count total
	countQuery := `SELECT COUNT(*) FROM master_barang WHERE 
	               nama_barang ILIKE $1 OR kode_barang ILIKE $1`
	searchPattern := "%" + search + "%"
	err := r.db.QueryRow(countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	query := `SELECT id, kode_barang, nama_barang, kategori, satuan, 
	          harga_beli, harga_jual, created_at, updated_at 
	          FROM master_barang 
	          WHERE nama_barang ILIKE $1 OR kode_barang ILIKE $1
	          ORDER BY id DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var b models.Barang
		err := rows.Scan(&b.ID, &b.KodeBarang, &b.NamaBarang, &b.Kategori,
			&b.Satuan, &b.HargaBeli, &b.HargaJual, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		barangs = append(barangs, b)
	}

	return barangs, total, nil
}

func (r *barangRepository) FindByID(id int) (*models.Barang, error) {
	barang := &models.Barang{}
	query := `SELECT id, kode_barang, nama_barang, kategori, satuan, 
	          harga_beli, harga_jual, created_at, updated_at 
	          FROM master_barang WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&barang.ID, &barang.KodeBarang, &barang.NamaBarang, &barang.Kategori,
		&barang.Satuan, &barang.HargaBeli, &barang.HargaJual,
		&barang.CreatedAt, &barang.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("barang not found")
	}
	if err != nil {
		return nil, err
	}

	return barang, nil
}

func (r *barangRepository) FindAllWithStok(search string, limit, offset int) ([]models.BarangWithStok, int, error) {
	var barangs []models.BarangWithStok
	var total int

	// Count total
	countQuery := `SELECT COUNT(*) FROM master_barang b
	               LEFT JOIN mstok s ON b.id = s.barang_id
	               WHERE b.nama_barang ILIKE $1 OR b.kode_barang ILIKE $1`
	searchPattern := "%" + search + "%"
	err := r.db.QueryRow(countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	query := `SELECT b.id, b.kode_barang, b.nama_barang, b.kategori, b.satuan,
	          b.harga_beli, b.harga_jual, b.created_at, b.updated_at,
	          COALESCE((SELECT SUM(h.qty) FROM history_stok h WHERE h.barang_id = b.id AND h.jenis_transaksi = 'masuk'), 0) as qty_masuk,
	          COALESCE((SELECT SUM(h.qty) FROM history_stok h WHERE h.barang_id = b.id AND h.jenis_transaksi = 'keluar'), 0) as qty_keluar,
	          COALESCE(s.stok_akhir, 0) as stok_akhir
	          FROM master_barang b
	          LEFT JOIN mstok s ON b.id = s.barang_id
	          WHERE b.nama_barang ILIKE $1 OR b.kode_barang ILIKE $1
	          ORDER BY b.id DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var b models.BarangWithStok
		err := rows.Scan(&b.ID, &b.KodeBarang, &b.NamaBarang, &b.Kategori,
			&b.Satuan, &b.HargaBeli, &b.HargaJual, &b.CreatedAt, &b.UpdatedAt,
			&b.QtyMasuk, &b.QtyKeluar, &b.StokAkhir)
		if err != nil {
			return nil, 0, err
		}
		barangs = append(barangs, b)
	}

	return barangs, total, nil
}

func (r *barangRepository) GenerateKodeBarang() (string, error) {
	var lastNumber int
	query := `SELECT COALESCE(MAX(CAST(SUBSTRING(kode_barang FROM 4) AS INTEGER)), 0) 
	          FROM master_barang WHERE kode_barang ~ '^BRG[0-9]+$'`

	err := r.db.QueryRow(query).Scan(&lastNumber)
	if err != nil {
		return "", err
	}

	nextNumber := lastNumber + 1
	return fmt.Sprintf("BRG%03d", nextNumber), nil
}

func (r *barangRepository) Create(barang *models.Barang) error {
	// Auto-generate kode barang if empty
	if barang.KodeBarang == "" {
		kode, err := r.GenerateKodeBarang()
		if err != nil {
			return err
		}
		barang.KodeBarang = kode
	}

	query := `INSERT INTO master_barang (kode_barang, nama_barang, kategori, satuan, harga_beli, harga_jual)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, barang.KodeBarang, barang.NamaBarang, barang.Kategori,
		barang.Satuan, barang.HargaBeli, barang.HargaJual).Scan(
		&barang.ID, &barang.CreatedAt, &barang.UpdatedAt,
	)
}

func (r *barangRepository) Update(barang *models.Barang) error {
	query := `UPDATE master_barang SET nama_barang = $1, kategori = $2, satuan = $3,
	          harga_beli = $4, harga_jual = $5 WHERE id = $6`

	result, err := r.db.Exec(query, barang.NamaBarang, barang.Kategori, barang.Satuan,
		barang.HargaBeli, barang.HargaJual, barang.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("barang not found")
	}

	return nil
}

func (r *barangRepository) Delete(id int) error {
	query := `DELETE FROM master_barang WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("barang not found")
	}

	return nil
}
