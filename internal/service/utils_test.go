package service

import (
	"kushkiv2/internal/db"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB configura una base de datos en memoria para testing.
// Se comparte entre todos los tests del paquete service.
func setupTestDB() *gorm.DB {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Inyectar la DB de prueba
	db.SetDB(database)

	// Correr migraciones
	db.Migrate(database)

	return database
}
