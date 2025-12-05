-- Migration: Create warehouse database schema
-- Description: Complete database schema for warehouse inventory management system

-- Drop tables if exist (for clean migration)
DROP TABLE IF EXISTS jual_detail CASCADE;
DROP TABLE IF EXISTS jual_header CASCADE;
DROP TABLE IF EXISTS beli_detail CASCADE;
DROP TABLE IF EXISTS beli_header CASCADE;
DROP TABLE IF EXISTS history_stok CASCADE;
DROP TABLE IF EXISTS mstok CASCADE;
DROP TABLE IF EXISTS master_barang CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nama VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'staff')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Master Barang table
CREATE TABLE master_barang (
    id SERIAL PRIMARY KEY,
    kode_barang VARCHAR(50) UNIQUE NOT NULL,
    nama_barang VARCHAR(200) NOT NULL,
    kategori VARCHAR(100),
    satuan VARCHAR(20) NOT NULL,
    harga_beli DECIMAL(15, 2) NOT NULL DEFAULT 0,
    harga_jual DECIMAL(15, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Stock table
CREATE TABLE mstok (
    id SERIAL PRIMARY KEY,
    barang_id INT NOT NULL REFERENCES master_barang(id) ON DELETE CASCADE,
    stok_awal INT NOT NULL DEFAULT 0,
    stok_masuk INT NOT NULL DEFAULT 0,
    stok_keluar INT NOT NULL DEFAULT 0,
    stok_akhir INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(barang_id)
);

-- Stock History table
CREATE TABLE history_stok (
    id SERIAL PRIMARY KEY,
    barang_id INT NOT NULL REFERENCES master_barang(id) ON DELETE CASCADE,
    jenis_transaksi VARCHAR(20) NOT NULL CHECK (jenis_transaksi IN ('masuk', 'keluar')),
    qty INT NOT NULL,
    stok_sebelum INT NOT NULL,
    stok_sesudah INT NOT NULL,
    keterangan TEXT,
    referensi_id INT,
    referensi_tipe VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Purchase Header table
CREATE TABLE beli_header (
    id SERIAL PRIMARY KEY,
    no_faktur VARCHAR(50) UNIQUE NOT NULL,
    tanggal DATE NOT NULL,
    supplier VARCHAR(200) NOT NULL,
    total DECIMAL(15, 2) NOT NULL DEFAULT 0,
    keterangan TEXT,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Purchase Detail table
CREATE TABLE beli_detail (
    id SERIAL PRIMARY KEY,
    beli_header_id INT NOT NULL REFERENCES beli_header(id) ON DELETE CASCADE,
    barang_id INT NOT NULL REFERENCES master_barang(id),
    qty INT NOT NULL,
    harga DECIMAL(15, 2) NOT NULL,
    subtotal DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sales Header table
CREATE TABLE jual_header (
    id SERIAL PRIMARY KEY,
    no_faktur VARCHAR(50) UNIQUE NOT NULL,
    tanggal DATE NOT NULL,
    customer VARCHAR(200) NOT NULL,
    total DECIMAL(15, 2) NOT NULL DEFAULT 0,
    keterangan TEXT,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sales Detail table
CREATE TABLE jual_detail (
    id SERIAL PRIMARY KEY,
    jual_header_id INT NOT NULL REFERENCES jual_header(id) ON DELETE CASCADE,
    barang_id INT NOT NULL REFERENCES master_barang(id),
    qty INT NOT NULL,
    harga DECIMAL(15, 2) NOT NULL,
    subtotal DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_master_barang_kode ON master_barang(kode_barang);
CREATE INDEX idx_mstok_barang_id ON mstok(barang_id);
CREATE INDEX idx_history_stok_barang_id ON history_stok(barang_id);
CREATE INDEX idx_history_stok_jenis ON history_stok(jenis_transaksi);
CREATE INDEX idx_beli_header_no_faktur ON beli_header(no_faktur);
CREATE INDEX idx_beli_detail_header_id ON beli_detail(beli_header_id);
CREATE INDEX idx_jual_header_no_faktur ON jual_header(no_faktur);
CREATE INDEX idx_jual_detail_header_id ON jual_detail(jual_header_id);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_master_barang_updated_at BEFORE UPDATE ON master_barang
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_mstok_updated_at BEFORE UPDATE ON mstok
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_beli_header_updated_at BEFORE UPDATE ON beli_header
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_jual_header_updated_at BEFORE UPDATE ON jual_header
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
