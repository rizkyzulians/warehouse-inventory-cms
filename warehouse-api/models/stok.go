package models

import "time"

type Stok struct {
	ID         int       `json:"id"`
	BarangID   int       `json:"barang_id"`
	StokAwal   int       `json:"stok_awal"`
	StokMasuk  int       `json:"stok_masuk"`
	StokKeluar int       `json:"stok_keluar"`
	StokAkhir  int       `json:"stok_akhir"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type StokWithBarang struct {
	Stok
	KodeBarang string `json:"kode_barang"`
	NamaBarang string `json:"nama_barang"`
	Satuan     string `json:"satuan"`
}
