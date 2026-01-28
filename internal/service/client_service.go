package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"kushkiv2/internal/db"
	"strings"
)

type ClientService struct{}

func NewClientService() *ClientService {
	return &ClientService{}
}

// ImportClientsFromCSV lee un CSV e inserta/actualiza clientes en la base de datos.
// Formato esperado: ID (RUC/Cédula), TipoID, Nombre, Dirección, Email, Teléfono
func (s *ClientService) ImportClientsFromCSV(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)
	// Saltar cabecera si existe
	firstRow, err := csvReader.Read()
	if err != nil {
		return 0, fmt.Errorf("error leyendo CSV: %v", err)
	}

	isHeader := false
	for _, cell := range firstRow {
		lower := strings.ToLower(cell)
		if lower == "id" || lower == "ruc" || lower == "cedula" || lower == "nombre" {
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
			continue // Mínimo ID, TipoID, Nombre
		}

		id := strings.TrimSpace(record[0])
		tipoID := strings.TrimSpace(record[1])
		nombre := strings.TrimSpace(record[2])
		
		direccion := ""
		if len(record) > 3 {
			direccion = strings.TrimSpace(record[3])
		}

		email := ""
		if len(record) > 4 {
			email = strings.TrimSpace(record[4])
		}

		telefono := ""
		if len(record) > 5 {
			telefono = strings.TrimSpace(record[5])
		}

		if id == "" || nombre == "" {
			continue
		}

		if tipoID == "" {
			if len(id) == 13 {
				tipoID = "04" // RUC
			} else if len(id) == 10 {
				tipoID = "05" // Cédula
			} else {
				tipoID = "08" // Pasaporte/Identificación Exterior
			}
		}

		client := db.Client{
			ID:        id,
			TipoID:    tipoID,
			Nombre:    nombre,
			Direccion: direccion,
			Email:     email,
			Telefono:  telefono,
		}

		// Upsert logic
		var existing db.Client
		if err := db.GetDB().Where("id = ?", id).First(&existing).Error; err == nil {
			db.GetDB().Model(&existing).Updates(client)
		} else {
			db.GetDB().Create(&client)
		}
		importedCount++
	}

	return importedCount, nil
}
