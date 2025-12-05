package repositories

import (
	"database/sql"
	"fmt"
	"warehouse-api/models"
)

type StokRepository interface {
	FindAll() ([]models.StokWithBarang, error)
	FindByBarangID(barangID int) (*models.Stok, error)
	UpdateStok(tx *sql.Tx, barangID int, stokMasuk int, stokKeluar int) error
	CreateStok(tx *sql.Tx, barangID int) error
	InsertHistory(tx *sql.Tx, history *models.HistoryStok) error
	GetHistoryAll(limit, offset int) ([]models.HistoryStokWithBarang, int, error)
	GetHistoryByBarangID(barangID int, limit, offset int) ([]models.HistoryStok, int, error)
}

type stokRepository struct {
	db *sql.DB
}

func NewStokRepository(db *sql.DB) StokRepository {
	return &stokRepository{db: db}
}

func (r *stokRepository) FindAll() ([]models.StokWithBarang, error) {
	var stoks []models.StokWithBarang

	query := `SELECT s.id, s.barang_id, s.stok_awal, s.stok_masuk, s.stok_keluar, 
	          s.stok_akhir, s.created_at, s.updated_at,
	          b.kode_barang, b.nama_barang, b.satuan
	          FROM mstok s
	          JOIN master_barang b ON s.barang_id = b.id
	          ORDER BY s.id DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.StokWithBarang
		err := rows.Scan(&s.ID, &s.BarangID, &s.StokAwal, &s.StokMasuk,
			&s.StokKeluar, &s.StokAkhir, &s.CreatedAt, &s.UpdatedAt,
			&s.KodeBarang, &s.NamaBarang, &s.Satuan)
		if err != nil {
			return nil, err
		}
		stoks = append(stoks, s)
	}

	return stoks, nil
}

func (r *stokRepository) FindByBarangID(barangID int) (*models.Stok, error) {
	stok := &models.Stok{}
	query := `SELECT id, barang_id, stok_awal, stok_masuk, stok_keluar, stok_akhir, 
	          created_at, updated_at FROM mstok WHERE barang_id = $1`

	err := r.db.QueryRow(query, barangID).Scan(
		&stok.ID, &stok.BarangID, &stok.StokAwal, &stok.StokMasuk,
		&stok.StokKeluar, &stok.StokAkhir, &stok.CreatedAt, &stok.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return stok, nil
}

func (r *stokRepository) UpdateStok(tx *sql.Tx, barangID int, stokMasuk int, stokKeluar int) error {
	query := `UPDATE mstok SET 
	          stok_masuk = stok_masuk + $1,
	          stok_keluar = stok_keluar + $2,
	          stok_akhir = stok_akhir + $1 - $2
	          WHERE barang_id = $3`

	result, err := tx.Exec(query, stokMasuk, stokKeluar, barangID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stok not found for barang_id: %d", barangID)
	}

	return nil
}

func (r *stokRepository) CreateStok(tx *sql.Tx, barangID int) error {
	query := `INSERT INTO mstok (barang_id, stok_awal, stok_masuk, stok_keluar, stok_akhir)
	          VALUES ($1, 0, 0, 0, 0)`

	_, err := tx.Exec(query, barangID)
	return err
}

func (r *stokRepository) InsertHistory(tx *sql.Tx, history *models.HistoryStok) error {
	query := `INSERT INTO history_stok (barang_id, jenis_transaksi, qty, stok_sebelum, 
	          stok_sesudah, keterangan, referensi_id, referensi_tipe)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`

	return tx.QueryRow(query, history.BarangID, history.JenisTransaksi, history.Qty,
		history.StokSebelum, history.StokSesudah, history.Keterangan,
		history.ReferensiID, history.ReferensiTipe).Scan(&history.ID, &history.CreatedAt)
}

func (r *stokRepository) GetHistoryAll(limit, offset int) ([]models.HistoryStokWithBarang, int, error) {
	var histories []models.HistoryStokWithBarang
	var total int

	// Count total
	countQuery := `SELECT COUNT(*) FROM history_stok`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := `SELECT h.id, h.barang_id, h.jenis_transaksi, h.qty, h.stok_sebelum, 
	          h.stok_sesudah, h.keterangan, h.referensi_id, h.referensi_tipe, h.created_at,
	          b.kode_barang, b.nama_barang
	          FROM history_stok h
	          JOIN master_barang b ON h.barang_id = b.id
	          ORDER BY h.created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.HistoryStokWithBarang
		err := rows.Scan(&h.ID, &h.BarangID, &h.JenisTransaksi, &h.Qty,
			&h.StokSebelum, &h.StokSesudah, &h.Keterangan, &h.ReferensiID,
			&h.ReferensiTipe, &h.CreatedAt, &h.KodeBarang, &h.NamaBarang)
		if err != nil {
			return nil, 0, err
		}
		histories = append(histories, h)
	}

	return histories, total, nil
}

func (r *stokRepository) GetHistoryByBarangID(barangID int, limit, offset int) ([]models.HistoryStok, int, error) {
	var histories []models.HistoryStok
	var total int

	// Count total
	countQuery := `SELECT COUNT(*) FROM history_stok WHERE barang_id = $1`
	err := r.db.QueryRow(countQuery, barangID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get data
	query := `SELECT id, barang_id, jenis_transaksi, qty, stok_sebelum, stok_sesudah,
	          keterangan, referensi_id, referensi_tipe, created_at
	          FROM history_stok WHERE barang_id = $1
	          ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, barangID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.HistoryStok
		err := rows.Scan(&h.ID, &h.BarangID, &h.JenisTransaksi, &h.Qty,
			&h.StokSebelum, &h.StokSesudah, &h.Keterangan, &h.ReferensiID,
			&h.ReferensiTipe, &h.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		histories = append(histories, h)
	}

	return histories, total, nil
}
