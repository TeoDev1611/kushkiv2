package pdf_test

import (
	"kushkiv2/pkg/pdf"
	srixml "kushkiv2/pkg/xml"
	"testing"
)

func TestGenerarRIDE(t *testing.T) {
	// 1. Setup Data Mock
	factura := srixml.FacturaXML{
		InfoTributaria: srixml.InfoTributaria{
			RazonSocial: "Empresa de Prueba S.A.",
			Ruc:         "1790000000001",
			ClaveAcceso: "2801202601179000000000120010010000000011234567813",
			Secuencial:  "000000001",
			Estab:       "001",
			PtoEmi:      "001",
			DirMatriz:   "Av. Amazonas y Naciones Unidas",
			Ambiente:    "1",
		},
		InfoFactura: srixml.InfoFactura{
			FechaEmision:            "28/01/2026",
			DirEstablecimiento:      "Av. Amazonas y Naciones Unidas",
			RazonSocialComprador:    "Cliente Final",
			IdentificacionComprador: "9999999999999",
			DireccionComprador:      "Quito",
			ImporteTotal:            115.00,
			TotalDescuento:          0.00,
			TotalConImpuestos: []srixml.TotalImpuesto{
				{Codigo: "2", CodigoPorcentaje: "2", BaseImponible: 100.00, Valor: 15.00},
			},
		},
		Detalles: []srixml.Detalle{
			{
				CodigoPrincipal:        "PROD001",
				Descripcion:            "Licencia de Software",
				Cantidad:               1,
				PrecioUnitario:         100.00,
				PrecioTotalSinImpuesto: 100.00,
			},
		},
		InfoAdicional: []srixml.CampoAdicional{
			{Nombre: "Email", Value: "cliente@test.com"},
		},
	}

	themes := []string{"modern", "minimal", "corporate"}

	for _, theme := range themes {
		t.Run("Theme_"+theme, func(t *testing.T) {
			// Test generation without logo
			bytes, err := pdf.GenerarRIDE(factura, "", theme)
			if err != nil {
				t.Fatalf("Error generando PDF con tema %s: %v", theme, err)
			}
			if len(bytes) == 0 {
				t.Errorf("PDF generado vacío para tema %s", theme)
			}
			
			// Simple magic bytes check for PDF
			if string(bytes[0:4]) != "%PDF" {
				t.Errorf("El archivo generado no parece ser un PDF válido (header incorrecto)")
			}
		})
	}

	t.Run("UnknownTheme_DefaultToModern", func(t *testing.T) {
		bytes, err := pdf.GenerarRIDE(factura, "", "unknown_theme_xyz")
		if err != nil {
			t.Fatalf("Error con tema desconocido: %v", err)
		}
		if len(bytes) == 0 {
			t.Error("Debería generar PDF (fallback a modern)")
		}
	})
}
