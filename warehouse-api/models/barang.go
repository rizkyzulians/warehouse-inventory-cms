package models

import "time"

type Barang struct {
	ID         int       `json:"id"`
	KodeBarang string    `json:"kode_barang"`
	NamaBarang string    `json:"nama_barang"`
	Kategori   string    `json:"kategori"`
	Satuan     string    `json:"satuan"`
	HargaBeli  float64   `json:"harga_beli"`
	HargaJual  float64   `json:"harga_jual"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BarangWithStok struct {
	Barang
	QtyMasuk  int `json:"qty_masuk"`
	QtyKeluar int `json:"qty_keluar"`
	StokAkhir int `json:"qty_akhir"`
}

type CreateBarangRequest struct {
	KodeBarang string  `json:"kode_barang"`
	NamaBarang string  `json:"nama_barang"`
	Kategori   string  `json:"kategori"`
	Satuan     string  `json:"satuan"`
	HargaBeli  float64 `json:"harga_beli"`
	HargaJual  float64 `json:"harga_jual"`
}

type UpdateBarangRequest struct {
	NamaBarang string  `json:"nama_barang"`
	Kategori   string  `json:"kategori"`
	Satuan     string  `json:"satuan"`
	HargaBeli  float64 `json:"harga_beli"`
	HargaJual  float64 `json:"harga_jual"`
}
