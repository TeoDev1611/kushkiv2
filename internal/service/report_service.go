package service

import (
	"fmt"
	"kushkiv2/internal/db"
	"time"

	"github.com/xuri/excelize/v2"
)

type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

// GenerateSalesExcel genera un archivo Excel con las ventas en un rango de fechas.
func (s *ReportService) GenerateSalesExcel(startDate, endDate time.Time) ([]byte, error) {
	var facturas []db.Factura
	// Buscar facturas autorizadas en el rango
	err := db.GetDB().Where("fecha_emision BETWEEN ? AND ?", startDate, endDate).
		Order("fecha_emision asc").Find(&facturas).Error
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := "Ventas"
	index, _ := f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// Estilos
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"34D399"}, Pattern: 1},
	})

	// Encabezados
	headers := []string{"Fecha", "Secuencial", "Clave de Acceso", "Cliente ID", "Subtotal 15%", "Subtotal 0%", "IVA", "Total", "Estado"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Datos
	for i, fact := range facturas {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), fact.FechaEmision.Format("02/01/2006 15:04"))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fact.Secuencial)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), fact.ClaveAcceso)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), fact.ClienteID)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), fact.Subtotal15)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), fact.Subtotal0)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), fact.IVA)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), fact.Total)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), fact.EstadoSRI)
	}

	f.SetActiveSheet(index)
	
	// Guardar a buffer de memoria
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type TopProduct struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Total    float64 `json:"total"`
}

// GetTopProducts obtiene los productos más vendidos (Basado en facturas con valor).
func (s *ReportService) GetTopProducts(limit int) ([]TopProduct, error) {
	var results []TopProduct
	
	// Usamos alias explícitos para que coincidan con los tags JSON: sku, name, quantity, total
	err := db.GetDB().Table("factura_items").
		Select("factura_items.producto_sku as sku, factura_items.nombre as name, SUM(factura_items.cantidad) as quantity, SUM(factura_items.subtotal) as total").
		Joins("JOIN facturas ON facturas.clave_acceso = factura_items.factura_clave").
		Where("facturas.total > 0").
		Group("factura_items.producto_sku, factura_items.nombre").
		Order("quantity DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}
