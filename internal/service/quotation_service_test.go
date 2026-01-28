package service

import (
	"kushkiv2/internal/db"
	"testing"
)

func TestQuotation_GetNextSecuencial(t *testing.T) {
	database := setupTestDB()
	svc := NewQuotationService()

	// 1. Caso Base: No existen cotizaciones
	sec, err := svc.GetNextSecuencial()
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if sec != "000000001" {
		t.Errorf("Esperado 000000001, obtenido %s", sec)
	}

	// 2. Insertar cotización
	database.Create(&db.Quotation{
		Secuencial: "000000005",
		Total:      100.00,
	})

	sec, err = svc.GetNextSecuencial()
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if sec != "000000006" {
		t.Errorf("Esperado 000000006, obtenido %s", sec)
	}
}

func TestCreateQuotation(t *testing.T) {
	database := setupTestDB()
	svc := NewQuotationService()
	
	// Crear configuración de emisor necesaria para PDF
	database.Create(&db.EmisorConfig{
		RUC: "1234567890001",
		RazonSocial: "Empresa Test",
	})

	dto := &db.QuotationDTO{
		Secuencial:      "000000001",
		ClienteID:       "1712345678",
		ClienteNombre:   "Juan Perez",
		ClienteEmail:    "juan@test.com",
		Items: []db.QuotationItemDTO{
			{Codigo: "P1", Nombre: "Prod 1", Cantidad: 1, Precio: 100, PorcentajeIVA: 15},
			{Codigo: "P2", Nombre: "Prod 2", Cantidad: 2, Precio: 50, PorcentajeIVA: 0},
		},
	}

	err := svc.CreateQuotation(dto)
	if err != nil {
		t.Fatalf("Error creando cotización: %v", err)
	}

	// Verificar DB
	var q db.Quotation
	if err := database.First(&q, "secuencial = ?", "000000001").Error; err != nil {
		t.Fatalf("Cotización no guardada en DB")
	}

	// Total = (100 * 1.15) + (100) = 115 + 100 = 215
	if q.Total != 215.00 {
		t.Errorf("Total incorrecto: esperado 215.00, obtenido %.2f", q.Total)
	}
	
	// Verificar Cliente Auto-guardado
	var c db.Client
	if err := database.First(&c, "id = ?", dto.ClienteID).Error; err != nil {
		t.Errorf("Cliente no fue guardado automáticamente")
	}
}

func TestConvertToInvoice(t *testing.T) {
	database := setupTestDB()
	svc := NewQuotationService()

	// Crear cotización base
	q := db.Quotation{
		Secuencial:    "000000001",
		ClienteID:     "1712345678",
		ClienteNombre: "Juan Perez",
		Total:         115.00,
	}
	database.Create(&q)
	database.Create(&db.QuotationItem{
		QuotationID:    q.ID,
		ProductoSKU:    "P1",
		Nombre:         "Prod 1",
		Cantidad:       1,
		PrecioUnitario: 100,
		PorcentajeIVA:  15,
	})

	// Convertir
	invoiceDTO, err := svc.ConvertToInvoice(q.ID)
	if err != nil {
		t.Fatalf("Error convirtiendo: %v", err)
	}

	if invoiceDTO.ClienteNombre != q.ClienteNombre {
		t.Errorf("Nombre cliente incorrecto en factura")
	}
	if len(invoiceDTO.Items) != 1 {
		t.Errorf("Items incorrectos: esperado 1, obtenido %d", len(invoiceDTO.Items))
	}
	if invoiceDTO.Items[0].Codigo != "P1" {
		t.Errorf("Código item incorrecto")
	}
}
