package service

import (
	"kushkiv2/internal/db"
	"testing"
)

func TestIntegration_QuotationToInvoiceFlow(t *testing.T) {
	database := setupTestDB()
	qSvc := NewQuotationService()
	
	// 1. Configurar Emisor (Necesario para InvoiceService)
	database.Create(&db.EmisorConfig{
		RUC: "1712345678001",
		RazonSocial: "Empresa de Pruebas",
		Estab: "001",
		PtoEmi: "001",
		Ambiente: 1,
	})

	// 2. Crear una Cotización
	qDTO := &db.QuotationDTO{
		Secuencial: "000000001",
		ClienteID: "1700000001",
		ClienteNombre: "Comprador de Prueba",
		Items: []db.QuotationItemDTO{
			{Codigo: "P1", Nombre: "Prod 1", Cantidad: 2, Precio: 10, PorcentajeIVA: 15},
		},
	}
	err := qSvc.CreateQuotation(qDTO)
	if err != nil {
		t.Fatalf("Error creando cotización: %v", err)
	}

	var q db.Quotation
	database.First(&q, "secuencial = ?", "000000001")

	// 3. Convertir a FacturaDTO
	fDTO, err := qSvc.ConvertToInvoice(q.ID)
	if err != nil {
		t.Fatalf("Error convirtiendo: %v", err)
	}

	if fDTO.ClienteNombre != "Comprador de Prueba" {
		t.Errorf("Nombre de cliente no coincide")
	}

	// 4. Emitir Factura
	// Nota: EmitirFactura requiere firma electrónica válida si se intenta firmar.
	// Sin embargo, InvoiceService podría tener mocks o fallar en el paso de firma.
	// Vamos a ver si InvoiceService permite emitir sin firma en modo test o si falla graciosamente.
	// REVISIÓN: EmitirFactura llama a SRI. 
	// Para propósitos de este test, solo verificamos que la conversión de datos fue correcta.
	
	if len(fDTO.Items) != 1 {
		t.Errorf("Cantidad de items incorrecta")
	}
	if fDTO.Items[0].Codigo != "P1" {
		t.Errorf("Código de item incorrecto")
	}
}
