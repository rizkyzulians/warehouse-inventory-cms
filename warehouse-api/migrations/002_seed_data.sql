-- Seed data for warehouse inventory system
-- Password for all users: "password123" (hashed with bcrypt)

-- Insert users
-- Password for all users: "password123"
INSERT INTO users (username, password, nama, role) VALUES
('admin', '$2a$10$NTtChK8Jva6eu5QDk52wbevhlKnuxBLWr3qHJvwGfQGO4K09eKSAq', 'Administrator', 'admin'),
('staff1', '$2a$10$NTtChK8Jva6eu5QDk52wbevhlKnuxBLWr3qHJvwGfQGO4K09eKSAq', 'Staff Satu', 'staff'),
('staff2', '$2a$10$NTtChK8Jva6eu5QDk52wbevhlKnuxBLWr3qHJvwGfQGO4K09eKSAq', 'Staff Dua', 'staff');

-- Insert master barang
INSERT INTO master_barang (kode_barang, nama_barang, kategori, satuan, harga_beli, harga_jual) VALUES
('BRG001', 'Laptop Dell Latitude 5420', 'Elektronik', 'Unit', 8500000.00, 11000000.00),
('BRG002', 'Mouse Logitech M170', 'Aksesoris', 'Unit', 50000.00, 75000.00),
('BRG003', 'Keyboard Mechanical RGB', 'Aksesoris', 'Unit', 350000.00, 500000.00),
('BRG004', 'Monitor LG 24 Inch', 'Elektronik', 'Unit', 1500000.00, 2000000.00),
('BRG005', 'Printer Canon IP2770', 'Elektronik', 'Unit', 650000.00, 900000.00),
('BRG006', 'Hard Disk External 1TB', 'Storage', 'Unit', 700000.00, 950000.00),
('BRG007', 'RAM DDR4 8GB', 'Komponen', 'Unit', 450000.00, 600000.00),
('BRG008', 'SSD Samsung 256GB', 'Storage', 'Unit', 500000.00, 700000.00),
('BRG009', 'Webcam Logitech C270', 'Aksesoris', 'Unit', 350000.00, 500000.00),
('BRG010', 'Headset Gaming RGB', 'Aksesoris', 'Unit', 250000.00, 400000.00);

-- Insert initial stock (stok_awal = stok_akhir)
INSERT INTO mstok (barang_id, stok_awal, stok_masuk, stok_keluar, stok_akhir) VALUES
(1, 10, 0, 0, 10),
(2, 50, 0, 0, 50),
(3, 30, 0, 0, 30),
(4, 15, 0, 0, 15),
(5, 20, 0, 0, 20),
(6, 25, 0, 0, 25),
(7, 40, 0, 0, 40),
(8, 35, 0, 0, 35),
(9, 18, 0, 0, 18),
(10, 22, 0, 0, 22);

-- Insert sample purchase transactions
INSERT INTO beli_header (no_faktur, tanggal, supplier, total, keterangan, created_by) VALUES
('PO-2025-001', '2025-01-15', 'PT Supplier Elektronik', 25000000.00, 'Pembelian rutin bulan Januari', 1),
('PO-2025-002', '2025-02-10', 'CV Aksesori Komputer', 5000000.00, 'Stok aksesoris', 1);

-- Insert purchase details
INSERT INTO beli_detail (beli_header_id, barang_id, qty, harga, subtotal) VALUES
(1, 1, 2, 8500000.00, 17000000.00),
(1, 4, 4, 1500000.00, 6000000.00),
(1, 5, 2, 650000.00, 1300000.00),
(2, 2, 20, 50000.00, 1000000.00),
(2, 3, 10, 350000.00, 3500000.00),
(2, 10, 15, 250000.00, 3750000.00);

-- Insert sample sales transactions
INSERT INTO jual_header (no_faktur, tanggal, customer, total, keterangan, created_by) VALUES
('SO-2025-001', '2025-01-20', 'PT Client Teknologi', 22000000.00, 'Penjualan ke corporate client', 2),
('SO-2025-002', '2025-02-15', 'Toko Komputer Jaya', 3450000.00, 'Penjualan retail', 2);

-- Insert sales details
INSERT INTO jual_detail (jual_header_id, barang_id, qty, harga, subtotal) VALUES
(1, 1, 2, 11000000.00, 22000000.00),
(2, 2, 10, 75000.00, 750000.00),
(2, 3, 5, 500000.00, 2500000.00),
(2, 9, 2, 500000.00, 1000000.00);

-- Insert stock history for initial stock
INSERT INTO history_stok (barang_id, jenis_transaksi, qty, stok_sebelum, stok_sesudah, keterangan, referensi_tipe) VALUES
(1, 'masuk', 10, 0, 10, 'Stok awal', 'INITIAL'),
(2, 'masuk', 50, 0, 50, 'Stok awal', 'INITIAL'),
(3, 'masuk', 30, 0, 30, 'Stok awal', 'INITIAL'),
(4, 'masuk', 15, 0, 15, 'Stok awal', 'INITIAL'),
(5, 'masuk', 20, 0, 20, 'Stok awal', 'INITIAL'),
(6, 'masuk', 25, 0, 25, 'Stok awal', 'INITIAL'),
(7, 'masuk', 40, 0, 40, 'Stok awal', 'INITIAL'),
(8, 'masuk', 35, 0, 35, 'Stok awal', 'INITIAL'),
(9, 'masuk', 18, 0, 18, 'Stok awal', 'INITIAL'),
(10, 'masuk', 22, 0, 22, 'Stok awal', 'INITIAL');

-- Note: In production, purchase and sales transactions should update stock through application logic
-- The seed data above shows initial state only
