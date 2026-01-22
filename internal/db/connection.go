package db

import (
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// SetDB permite inyectar una instancia de base de datos (útil para tests).
func SetDB(database *gorm.DB) {
	db = database
	// Marcar once como hecho para evitar sobrescritura si se llama a GetDB después
	once.Do(func() {})
}

// GetDB devuelve la instancia única de la base de datos (Singleton).
func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		// Activamos _journal_mode=WAL para permitir concurrencia real (Lecturas y Escrituras simultáneas)
		// y _busy_timeout para esperar si hay bloqueos en lugar de fallar inmediatamente.
		db, err = gorm.Open(sqlite.Open("kushki.db?_journal_mode=WAL&_busy_timeout=5000"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Error al conectar con la base de datos: %v", err)
		}

		// Configurar Pool de Conexiones
		sqlDB, err := db.DB()
		if err == nil {
			// Aumentamos conexiones para soportar las goroutines del Dashboard y Sync
			sqlDB.SetMaxOpenConns(25) 
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetConnMaxLifetime(0)
		}
	})
	return db
}

// CloseDB cierra la conexión SQL subyacente.
func CloseDB() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
