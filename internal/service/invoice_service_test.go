package service

import (
	"kushkiv2/internal/db"
	"testing"
	"time"
)

func TestGetNextSecuencial(t *testing.T) {
	database := setupTestDB()
	svc := NewInvoiceService()

	// 1. Caso Base: No existen facturas
	sec, err := svc.GetNextSecuencial()
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if sec != "000000001" {
		t.Errorf("Esperado 000000001, obtenido %s", sec)
	}

	// 2. Caso con historial: Insertar factura '000000005'
	database.Create(&db.Factura{
		Secuencial:   "000000005",
		ClaveAcceso:  "dummy_key",
		FechaEmision: time.Now(),
		Total:        100.00,
	})

	sec, err = svc.GetNextSecuencial()
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if sec != "000000006" {
		t.Errorf("Esperado 000000006, obtenido %s", sec)
	}
}

func TestEmitirFactura_Validations(t *testing.T) {
	setupTestDB()
	svc := NewInvoiceService()

	// Caso 1: Item sin código (Regla 2025)
	dto := &db.FacturaDTO{
		Items: []db.InvoiceItem{
			{Codigo: "", Nombre: "Producto Test", Cantidad: 1, Precio: 10},
		},
	}
	err := svc.EmitirFactura(dto)
	if err == nil {
		t.Error("Se esperaba error por falta de código SKU, pero no ocurrió")
	}

	// Caso 2: Monto > $1000 con Pago 01 (Sin sistema financiero)
	dto2 := &db.FacturaDTO{
		FormaPago: "01",
		Items: []db.InvoiceItem{
			{Codigo: "SKU1", Nombre: "Caro", Cantidad: 1, Precio: 1200, PorcentajeIVA: 15, CodigoIVA: "4"},
		},
	}
	// Necesitamos configurar emisor para que no falle en el primer paso
	// (setupTestDB ya crea un emisor default con Migrate->seedEmisor)
	
	err = svc.EmitirFactura(dto2)
	if err == nil {
		t.Error("Se esperaba error por monto > $1000 sin sistema financiero")
	}
}
