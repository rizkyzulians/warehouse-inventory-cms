package models

import "time"

type HistoryStok struct {
	ID             int       `json:"id"`
	BarangID       int       `json:"barang_id"`
	JenisTransaksi string    `json:"jenis_transaksi"`
	Qty            int       `json:"qty"`
	StokSebelum    int       `json:"stok_sebelum"`
	StokSesudah    int       `json:"stok_sesudah"`
	Keterangan     string    `json:"keterangan"`
	ReferensiID    *int      `json:"referensi_id"`
	ReferensiTipe  string    `json:"referensi_tipe"`
	CreatedAt      time.Time `json:"created_at"`
}

type HistoryStokWithBarang struct {
	HistoryStok
	KodeBarang string `json:"kode_barang"`
	NamaBarang string `json:"nama_barang"`
}
