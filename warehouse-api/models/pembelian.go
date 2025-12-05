package models

import "time"

type BeliHeader struct {
	ID         int       `json:"id"`
	NoFaktur   string    `json:"no_faktur"`
	Tanggal    string    `json:"tanggal"`
	Supplier   string    `json:"supplier"`
	Total      float64   `json:"total"`
	Keterangan string    `json:"keterangan"`
	CreatedBy  int       `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BeliDetail struct {
	ID           int       `json:"id"`
	BeliHeaderID int       `json:"beli_header_id"`
	BarangID     int       `json:"barang_id"`
	Qty          int       `json:"qty"`
	Harga        float64   `json:"harga"`
	Subtotal     float64   `json:"subtotal"`
	CreatedAt    time.Time `json:"created_at"`
}

type BeliDetailWithBarang struct {
	BeliDetail
	KodeBarang string `json:"kode_barang"`
	NamaBarang string `json:"nama_barang"`
	Satuan     string `json:"satuan"`
}

type BeliHeaderWithDetail struct {
	BeliHeader
	Details []BeliDetailWithBarang `json:"details"`
}

type CreatePembelianRequest struct {
	NoFaktur   string                  `json:"no_faktur"`
	Tanggal    string                  `json:"tanggal"`
	Supplier   string                  `json:"supplier"`
	Keterangan string                  `json:"keterangan"`
	Details    []CreatePembelianDetail `json:"details"`
}

type CreatePembelianDetail struct {
	BarangID int     `json:"barang_id"`
	Qty      int     `json:"qty"`
	Harga    float64 `json:"harga"`
}
