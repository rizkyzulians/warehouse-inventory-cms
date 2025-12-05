package main

import (
	"log"
	"net/http"
	"warehouse-api/config"
	"warehouse-api/handlers"
	"warehouse-api/middleware"
	"warehouse-api/repositories"
	"warehouse-api/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	barangRepo := repositories.NewBarangRepository(db)
	stokRepo := repositories.NewStokRepository(db)
	pembelianRepo := repositories.NewPembelianRepository(db)
	penjualanRepo := repositories.NewPenjualanRepository(db)

	// Initialize services
	pembelianService := services.NewPembelianService(db, pembelianRepo, barangRepo, stokRepo)
	penjualanService := services.NewPenjualanService(db, penjualanRepo, barangRepo, stokRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)
	barangHandler := handlers.NewBarangHandler(barangRepo)
	stokHandler := handlers.NewStokHandler(stokRepo)
	pembelianHandler := handlers.NewPembelianHandler(pembelianService)
	penjualanHandler := handlers.NewPenjualanHandler(penjualanService)

	// Setup router
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware)

	// API prefix
	api := r.PathPrefix("/api").Subrouter()

	// Public routes (no authentication)
	api.HandleFunc("/login", userHandler.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET", "OPTIONS")

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// User routes
	protected.HandleFunc("/profile", userHandler.GetProfile).Methods("GET", "OPTIONS")

	// Barang routes (all authenticated users)
	protected.HandleFunc("/barang", barangHandler.GetAll).Methods("GET", "OPTIONS")
	protected.HandleFunc("/barang/stok", barangHandler.GetAllWithStok).Methods("GET", "OPTIONS")
	protected.HandleFunc("/barang/{id}", barangHandler.GetByID).Methods("GET", "OPTIONS")

	// Admin only routes for barang create/update
	adminBarang := protected.PathPrefix("").Subrouter()
	adminBarang.Use(middleware.RequireRole("admin"))
	adminBarang.HandleFunc("/barang", barangHandler.Create).Methods("POST", "OPTIONS")
	adminBarang.HandleFunc("/barang/{id}", barangHandler.Update).Methods("PUT", "OPTIONS")
	adminBarang.HandleFunc("/barang/{id}", barangHandler.Delete).Methods("DELETE", "OPTIONS")

	// Stok routes (specific routes BEFORE generic routes)
	protected.HandleFunc("/stok", stokHandler.GetAll).Methods("GET", "OPTIONS")
	protected.HandleFunc("/stok/history", stokHandler.GetHistoryAll).Methods("GET", "OPTIONS")
	protected.HandleFunc("/stok/history/{barang_id}", stokHandler.GetHistoryByBarangID).Methods("GET", "OPTIONS")
	protected.HandleFunc("/stok/{barang_id}", stokHandler.GetByBarangID).Methods("GET", "OPTIONS")

	// Pembelian routes
	protected.HandleFunc("/pembelian", pembelianHandler.GetAll).Methods("GET", "OPTIONS")
	protected.HandleFunc("/pembelian/{id}", pembelianHandler.GetByID).Methods("GET", "OPTIONS")
	protected.HandleFunc("/pembelian", pembelianHandler.Create).Methods("POST", "OPTIONS")

	// Penjualan routes
	protected.HandleFunc("/penjualan", penjualanHandler.GetAll).Methods("GET", "OPTIONS")
	protected.HandleFunc("/penjualan/{id}", penjualanHandler.GetByID).Methods("GET", "OPTIONS")
	protected.HandleFunc("/penjualan", penjualanHandler.Create).Methods("POST", "OPTIONS")

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server starting on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
