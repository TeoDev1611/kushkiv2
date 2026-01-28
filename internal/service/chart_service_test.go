package service

import (
	"kushkiv2/internal/db"
	"strings"
	"testing"
	"time"
)

func TestGenerateRevenueChart(t *testing.T) {
	database := setupTestDB()
	svc := NewChartService()

	// Insertar datos
	now := time.Now()
	// Mes actual
	database.Create(&db.Factura{ClaveAcceso: "1", FechaEmision: now, Total: 100, EstadoSRI: "AUTORIZADO"})
	database.Create(&db.Factura{ClaveAcceso: "2", FechaEmision: now, Total: 50, EstadoSRI: "AUTORIZADO"})
	// Mes pasado
	lastMonth := now.AddDate(0, -1, 0)
	database.Create(&db.Factura{ClaveAcceso: "3", FechaEmision: lastMonth, Total: 200, EstadoSRI: "AUTORIZADO"})
	// Factura no autorizada (no debe sumar)
	database.Create(&db.Factura{ClaveAcceso: "4", FechaEmision: now, Total: 500, EstadoSRI: "ANULADO"})

	html, err := svc.GenerateRevenueChart()
	if err != nil {
		t.Fatalf("Error generando chart: %v", err)
	}

	if !strings.Contains(html, "Evolución de Ingresos") {
		t.Errorf("El HTML no contiene el título esperado")
	}
	// go-echarts genera un HTML grande, verificar contenido exacto es difícil,
	// pero podemos verificar que no esté vacío.
	if len(html) < 100 {
		t.Errorf("HTML generado demasiado corto")
	}
}

func TestGenerateClientsPie(t *testing.T) {
	database := setupTestDB()
	svc := NewChartService()

	database.Create(&db.Client{ID: "1", Nombre: "Cliente Top"})
	database.Create(&db.Client{ID: "2", Nombre: "Cliente Small"})

	database.Create(&db.Factura{ClaveAcceso: "10", ClienteID: "1", Total: 1000, EstadoSRI: "AUTORIZADO"})
	database.Create(&db.Factura{ClaveAcceso: "11", ClienteID: "2", Total: 100, EstadoSRI: "AUTORIZADO"})

	html, err := svc.GenerateClientsPie()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if !strings.Contains(html, "Top 5 Clientes") {
		t.Errorf("HTML título incorrecto")
	}
	if !strings.Contains(html, "Cliente Top") {
		t.Errorf("HTML no contiene datos del cliente")
	}
}
