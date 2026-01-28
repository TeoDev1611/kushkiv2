package service

import (
	"kushkiv2/internal/db"
	"kushkiv2/pkg/util"
	"strings"
	"time"
)

type TaxSummary struct {
	Ventas15          float64 `json:"ventas15"`          // Campo 401
	Ventas0           float64 `json:"ventas0"`           // Campo 403
	IvaGenerado       float64 `json:"ivaGenerado"`       // Campo 411
	RetencionesIva    float64 `json:"retencionesIva"`    // Campo 609
	FactorProporcion  float64 `json:"factorProporcion"`  // Campo 702
	ImpuestoSugerido  float64 `json:"impuestoSugerido"`
}

type TaxService struct{}

func NewTaxService() *TaxService {
	return &TaxService{}
}

func (s *TaxService) GetVATSummary(startDate, endDate time.Time) TaxSummary {
	var summary TaxSummary
	var facturas []db.Factura

	// 1. Cargar facturas autorizadas del periodo
	db.GetDB().Where("fecha_emision >= ? AND fecha_emision <= ?", startDate, endDate).Find(&facturas)

	for _, f := range facturas {
		estado := strings.ToUpper(strings.TrimSpace(f.EstadoSRI))
		// Solo contamos las que tienen valor y no están anuladas (basado en lógica previa)
		if f.Total > 0 && estado != "ANULADO" {
			summary.Ventas15 += f.Subtotal15
			summary.Ventas0 += f.Subtotal0
			summary.IvaGenerado += f.IVA
		}
	}

	// 2. Calcular Factor de Proporcionalidad
	// Formula: Ventas con derecho a crédito (15%) / Ventas Totales
	totalVentas := summary.Ventas15 + summary.Ventas0
	if totalVentas > 0 {
		summary.FactorProporcion = util.Round(summary.Ventas15 / totalVentas, 4)
	} else {
		summary.FactorProporcion = 1.0
	}

	// 3. Obtener Retenciones Recibidas
	db.GetDB().Model(&db.RetencionRecibida{}).
		Where("fecha_emision >= ? AND fecha_emision <= ? AND tipo = ?", startDate, endDate, "IVA").
		Select("COALESCE(SUM(valor_retenido), 0)").Scan(&summary.RetencionesIva)

	// 4. Calcular Impuesto sugerido (IVA Cobrado - Retenciones)
	// Nota: No tenemos el IVA de compras (casilleros 500) aún, así que restamos retenciones del generado.
	res := summary.IvaGenerado - summary.RetencionesIva
	if res < 0 { res = 0 }
	summary.ImpuestoSugerido = util.Round(res, 2)

	return summary
}
