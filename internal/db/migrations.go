package db

import (
	"log"

	"gorm.io/gorm"
)

// RunMigrations ejecuta el AutoMigrate de GORM para asegurar que las tablas existan.
func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&EmisorConfig{},
		&Factura{},
		&Product{},
		&Client{},
		&EmailQueue{},
		&FacturaItem{},
		&MailLog{},
		&Quotation{},
		&QuotationItem{},
	)
	
	// OPTIMIZACIÓN: Índices manuales para el Dashboard y Buscador
	// GORM a veces no crea índices compuestos óptimos automáticamente
	if !db.Migrator().HasIndex(&Factura{}, "idx_factura_fecha_estado") {
		db.Exec("CREATE INDEX idx_factura_fecha_estado ON facturas(fecha_emision, estado_sri)")
	}
	if !db.Migrator().HasIndex(&Factura{}, "idx_factura_created") {
		db.Exec("CREATE INDEX idx_factura_created ON facturas(created_at DESC)")
	}
	if !db.Migrator().HasIndex(&Factura{}, "idx_factura_cliente_search") {
		db.Exec("CREATE INDEX idx_factura_cliente_search ON facturas(cliente_id, secuencial)")
	}

	seedEmisor(db)
}

func seedEmisor(db *gorm.DB) {
	var count int64
	db.Model(&EmisorConfig{}).Count(&count)
	if count == 0 {
		defaultEmisor := EmisorConfig{
			RUC:         "1790011223001",
			RazonSocial: "EMISOR DE PRUEBA S.A.",
			Ambiente:    1, // Pruebas
			Estab:       "001",
			PtoEmi:      "001",
			Obligado:    true,
			P12Path:     "/ruta/ficticia/certificado.p12", // Placeholder
			P12Password: "password123",                    // Placeholder
		}
		db.Create(&defaultEmisor)
		log.Println("Se ha creado un Emisor de prueba por defecto.")
	}
}
