package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"kushkiv2/internal/db"
	"strconv"
	"strings"
	"time"
)

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

// ImportProductsFromCSV lee un CSV e inserta/actualiza productos en la base de datos.
// Formato esperado: SKU, Nombre, Precio, Stock, CodigoImpuesto, PorcentajeIVA, Barcode, AuxiliaryCode, MinStock, ExpiryDate, Location
func (s *ProductService) ImportProductsFromCSV(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)
	// Saltar cabecera si existe (asumimos que la primera fila es cabecera si contiene "sku" o "nombre")
	firstRow, err := csvReader.Read()
	if err != nil {
		return 0, fmt.Errorf("error leyendo CSV: %v", err)
	}

	isHeader := false
	for _, cell := range firstRow {
		lower := strings.ToLower(cell)
		if lower == "sku" || lower == "nombre" || lower == "name" {
			isHeader = true
			break
		}
	}

	rowsToProcess := [][]string{}
	if !isHeader {
		rowsToProcess = append(rowsToProcess, firstRow)
	}

	importedCount := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return importedCount, fmt.Errorf("error leyendo fila: %v", err)
		}
		rowsToProcess = append(rowsToProcess, record)
	}

	for _, record := range rowsToProcess {
		if len(record) < 3 {
			continue // Fila inválida
		}

		sku := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])
		price, _ := strconv.ParseFloat(strings.ReplaceAll(record[2], ",", "."), 64)
		
		stock := 0
		if len(record) > 3 {
			stock, _ = strconv.Atoi(record[3])
		}

		taxCode := 2 // Default IVA
		if len(record) > 4 {
			taxCode, _ = strconv.Atoi(record[4])
		}

		taxPercentage := 15 // Default 15%
		if len(record) > 5 {
			taxPercentage, _ = strconv.Atoi(record[5])
		}

		var barcode, auxCode, location string
		var minStock int
		var expiryDate *time.Time

		if len(record) > 6 {
			barcode = strings.TrimSpace(record[6])
		}
		if len(record) > 7 {
			auxCode = strings.TrimSpace(record[7])
		}
		if len(record) > 8 {
			minStock, _ = strconv.Atoi(record[8])
		}
		if len(record) > 9 {
			expiryStr := strings.TrimSpace(record[9])
			if expiryStr != "" {
				parsed, err := time.Parse("2006-01-02", expiryStr)
				if err == nil {
					expiryDate = &parsed
				}
			}
		}
		if len(record) > 10 {
			location = strings.TrimSpace(record[10])
		}

		if sku == "" || name == "" {
			continue
		}

		product := db.Product{
			SKU:           sku,
			Name:          name,
			Price:         price,
			Stock:         stock,
			TaxCode:       taxCode,
			TaxPercentage: taxPercentage,
			Barcode:       barcode,
			AuxiliaryCode: auxCode,
			MinStock:      minStock,
			ExpiryDate:    expiryDate,
			Location:      location,
		}

		// Upsert logic
		var existing db.Product
		if err := db.GetDB().Where("sku = ?", sku).First(&existing).Error; err == nil {
			db.GetDB().Model(&existing).Updates(product)
		} else {
			db.GetDB().Create(&product)
		}
		importedCount++
	}

	return importedCount, nil
}
