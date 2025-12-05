package models

import "time"

type JualHeader struct {
	ID         int       `json:"id"`
	NoFaktur   string    `json:"no_faktur"`
	Tanggal    string    `json:"tanggal"`
	Customer   string    `json:"customer"`
	Total      float64   `json:"total"`
	Keterangan string    `json:"keterangan"`
	CreatedBy  int       `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type JualDetail struct {
	ID           int       `json:"id"`
	JualHeaderID int       `json:"jual_header_id"`
	BarangID     int       `json:"barang_id"`
	Qty          int       `json:"qty"`
	Harga        float64   `json:"harga"`
	Subtotal     float64   `json:"subtotal"`
	CreatedAt    time.Time `json:"created_at"`
}

type JualDetailWithBarang struct {
	JualDetail
	KodeBarang string `json:"kode_barang"`
	NamaBarang string `json:"nama_barang"`
	Satuan     string `json:"satuan"`
}

type JualHeaderWithDetail struct {
	JualHeader
	Details []JualDetailWithBarang `json:"details"`
}

type CreatePenjualanRequest struct {
	NoFaktur   string                  `json:"no_faktur"`
	Tanggal    string                  `json:"tanggal"`
	Customer   string                  `json:"customer"`
	Keterangan string                  `json:"keterangan"`
	Details    []CreatePenjualanDetail `json:"details"`
}

type CreatePenjualanDetail struct {
	BarangID int     `json:"barang_id"`
	Qty      int     `json:"qty"`
	Harga    float64 `json:"harga"`
}
