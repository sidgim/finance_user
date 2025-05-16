package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// carga .env en os.Getenv
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, usando variables de entorno del sistema")
	}

	srv, err := InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize: %v", err)
	}
	log.Printf("🚀 Server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
