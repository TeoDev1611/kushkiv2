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

// GenerateMasterReportExcel genera un reporte consolidado de todo el sistema.
func (s *ReportService) GenerateMasterReportExcel() ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// 0. HOJA DE RESUMEN
	sheetRes := "Resumen General"
	f.SetSheetName("Sheet1", sheetRes)
	
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"34D399"}, Pattern: 1},
	})

	f.SetCellValue(sheetRes, "A1", "Métrica")
	f.SetCellValue(sheetRes, "B1", "Valor")
	f.SetCellStyle(sheetRes, "A1", "B1", headerStyle)

	var totalVentas float64
	var countFacturas int64
	var countClientes int64
	var countProductos int64

	db.GetDB().Model(&db.Factura{}).Select("SUM(total)").Row().Scan(&totalVentas)
	db.GetDB().Model(&db.Factura{}).Count(&countFacturas)
	db.GetDB().Model(&db.Client{}).Count(&countClientes)
	db.GetDB().Model(&db.Product{}).Count(&countProductos)

	f.SetCellValue(sheetRes, "A2", "Total Ventas Histórico")
	f.SetCellValue(sheetRes, "B2", totalVentas)
	f.SetCellValue(sheetRes, "A3", "Total Facturas Emitidas")
	f.SetCellValue(sheetRes, "B3", countFacturas)
	f.SetCellValue(sheetRes, "A4", "Total Clientes")
	f.SetCellValue(sheetRes, "B4", countClientes)
	f.SetCellValue(sheetRes, "A5", "Total Productos en Inventario")
	f.SetCellValue(sheetRes, "B5", countProductos)
	f.SetCellValue(sheetRes, "A6", "Fecha de Generación")
	f.SetCellValue(sheetRes, "B6", time.Now().Format("02/01/2006 15:04"))

	// 1. HOJA DE VENTAS (HISTORIAL)
	sheetVentas := "Historial Ventas"
	f.NewSheet(sheetVentas)

	headersV := []string{"Fecha", "Secuencial", "Clave Acceso", "Cliente ID", "Total", "Estado"}
	for i, h := range headersV {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetVentas, cell, h)
		f.SetCellStyle(sheetVentas, cell, cell, headerStyle)
	}

	var facturas []db.Factura
	db.GetDB().Order("fecha_emision desc").Find(&facturas)
	for i, fact := range facturas {
		row := i + 2
		f.SetCellValue(sheetVentas, fmt.Sprintf("A%d", row), fact.FechaEmision.Format("02/01/2006 15:04"))
		f.SetCellValue(sheetVentas, fmt.Sprintf("B%d", row), fact.Secuencial)
		f.SetCellValue(sheetVentas, fmt.Sprintf("C%d", row), fact.ClaveAcceso)
		f.SetCellValue(sheetVentas, fmt.Sprintf("D%d", row), fact.ClienteID)
		f.SetCellValue(sheetVentas, fmt.Sprintf("E%d", row), fact.Total)
		f.SetCellValue(sheetVentas, fmt.Sprintf("F%d", row), fact.EstadoSRI)
	}

	// 2. HOJA DE CLIENTES
	sheetClientes := "Clientes"
	f.NewSheet(sheetClientes)
	headersC := []string{"Identificación", "Nombre", "Email", "Teléfono", "Dirección"}
	for i, h := range headersC {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetClientes, cell, h)
		f.SetCellStyle(sheetClientes, cell, cell, headerStyle)
	}

	var clients []db.Client
	db.GetDB().Find(&clients)
	for i, c := range clients {
		row := i + 2
		f.SetCellValue(sheetClientes, fmt.Sprintf("A%d", row), c.ID)
		f.SetCellValue(sheetClientes, fmt.Sprintf("B%d", row), c.Nombre)
		f.SetCellValue(sheetClientes, fmt.Sprintf("C%d", row), c.Email)
		f.SetCellValue(sheetClientes, fmt.Sprintf("D%d", row), c.Telefono)
		f.SetCellValue(sheetClientes, fmt.Sprintf("E%d", row), c.Direccion)
	}

	// 3. HOJA DE PRODUCTOS (INVENTARIO)
	sheetProductos := "Inventario"
	f.NewSheet(sheetProductos)
	headersP := []string{"SKU", "Nombre", "Precio", "Stock", "% IVA"}
	for i, h := range headersP {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetProductos, cell, h)
		f.SetCellStyle(sheetProductos, cell, cell, headerStyle)
	}

	var products []db.Product
	db.GetDB().Find(&products)
	for i, p := range products {
		row := i + 2
		f.SetCellValue(sheetProductos, fmt.Sprintf("A%d", row), p.SKU)
		f.SetCellValue(sheetProductos, fmt.Sprintf("B%d", row), p.Name)
		f.SetCellValue(sheetProductos, fmt.Sprintf("C%d", row), p.Price)
		f.SetCellValue(sheetProductos, fmt.Sprintf("D%d", row), p.Stock)
		f.SetCellValue(sheetProductos, fmt.Sprintf("E%d", row), p.TaxPercentage)
	}

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
