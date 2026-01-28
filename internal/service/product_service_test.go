package service

import (
	"kushkiv2/internal/db"
	"strings"
	"testing"
)

func TestImportProductsFromCSV(t *testing.T) {
	database := setupTestDB()
	svc := NewProductService()

	csvContent := `SKU,Nombre,Precio,Stock,CodigoImpuesto,PorcentajeIVA
P100,Producto Test 1,10.50,100,2,15
P200,Producto Test 2,20.00,50,2,12
`
	reader := strings.NewReader(csvContent)
	count, err := svc.ImportProductsFromCSV(reader)

	if err != nil {
		t.Fatalf("Error importando: %v", err)
	}

	if count != 2 {
		t.Errorf("Esperado 2 productos, obtenido %d", count)
	}

	var p1 db.Product
	if err := database.First(&p1, "sku = ?", "P100").Error; err != nil {
		t.Fatalf("Producto P100 no encontrado")
	}

	if p1.Price != 10.50 {
		t.Errorf("Precio incorrecto: %f", p1.Price)
	}

	if p1.Stock != 100 {
		t.Errorf("Stock incorrecto: %d", p1.Stock)
	}

	// Probar actualización (Upsert)
	csvContentUpdate := `SKU,Nombre,Precio,Stock,CodigoImpuesto,PorcentajeIVA
P100,Producto Test 1 Actualizado,15.00,200,2,15
`
	readerUpdate := strings.NewReader(csvContentUpdate)
	count, _ = svc.ImportProductsFromCSV(readerUpdate)

	if count != 1 {
		t.Errorf("Esperado 1 producto actualizado, obtenido %d", count)
	}

	database.First(&p1, "sku = ?", "P100")
	if p1.Name != "Producto Test 1 Actualizado" || p1.Price != 15.00 {
		t.Errorf("Actualización fallida: %+v", p1)
	}
}
