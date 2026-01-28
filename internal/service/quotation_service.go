package service

import (
	"fmt"
	"time"

	"kushkiv2/internal/db"
	"kushkiv2/pkg/logger"
	"kushkiv2/pkg/pdf"
	"kushkiv2/pkg/util"
)

type QuotationService struct{}

func NewQuotationService() *QuotationService {
	return &QuotationService{}
}

// GetNextSecuencial obtiene el siguiente número de cotización
func (s *QuotationService) GetNextSecuencial() (string, error) {
	var lastQuotes []db.Quotation
	// Usamos Find con Limit 1 para evitar error "record not found"
	db.GetDB().Order("created_at desc").Limit(1).Find(&lastQuotes)

	if len(lastQuotes) == 0 {
		return "000000001", nil
	}
	lastQuote := lastQuotes[0]

	var currentSec int
	fmt.Sscanf(lastQuote.Secuencial, "%d", &currentSec)
	return fmt.Sprintf("%09d", currentSec+1), nil
}

// CreateQuotation crea una nueva cotización y genera su PDF
func (s *QuotationService) CreateQuotation(dto *db.QuotationDTO) error {
	// 1. Validaciones básicas
	if len(dto.Items) == 0 {
		return fmt.Errorf("la cotización debe tener al menos un ítem")
	}

	// 2. Preparar Datos y Cálculos
	var subtotal15, subtotal0, totalIVA float64
	var itemsDB []db.QuotationItem

	for _, item := range dto.Items {
		base := util.Round(item.Cantidad*item.Precio, 2)
		impuesto := util.Round(base*(item.PorcentajeIVA/100), 2)
		
		if item.PorcentajeIVA > 0 {
			subtotal15 += base
			totalIVA += impuesto
		} else {
			subtotal0 += base
		}

		itemsDB = append(itemsDB, db.QuotationItem{
			ProductoSKU:    item.Codigo,
			Nombre:         item.Nombre,
			Cantidad:       item.Cantidad,
			PrecioUnitario: item.Precio,
			Subtotal:       base,
			PorcentajeIVA:  item.PorcentajeIVA,
		})
	}

	total := subtotal15 + subtotal0 + totalIVA

	// 3. Crear Registro en DB
	quotationDB := &db.Quotation{
		Secuencial:       dto.Secuencial,
		FechaEmision:     time.Now(),
		ClienteID:        dto.ClienteID,
		ClienteNombre:    dto.ClienteNombre,
		ClienteDireccion: dto.ClienteDireccion,
		ClienteEmail:     dto.ClienteEmail,
		ClienteTelefono:  dto.ClienteTelefono,
		Observacion:      dto.Observacion,
		Total:            util.Round(total, 2),
		Subtotal15:       util.Round(subtotal15, 2),
		Subtotal0:        util.Round(subtotal0, 2),
		IVA:              util.Round(totalIVA, 2),
		Estado:           "GENERADA",
	}

	// 4. Generar PDF
	var config db.EmisorConfig
	if err := db.GetDB().First(&config).Error; err != nil {
		logger.Error("Error obteniendo configuración para PDF cotización: %v", err)
	} else {
		pdfBytes, errPDF := pdf.GenerarCotizacionPDF(*quotationDB, itemsDB, config)
		if errPDF != nil {
			logger.Error("Error generando PDF cotización: %v", errPDF)
		} else {
			quotationDB.PDFBytes = pdfBytes
		}
	}

	// 5. Transacción DB
	tx := db.GetDB().Begin()
	if err := tx.Create(quotationDB).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error guardando cotización: %v", err)
	}

	for i := range itemsDB {
		itemsDB[i].QuotationID = quotationDB.ID
		if err := tx.Create(&itemsDB[i]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error guardando item cotización: %v", err)
		}
	}
	
	// Auto-guardar cliente (Upsert)
	var cliente db.Client
	if err := tx.First(&cliente, "id = ?", dto.ClienteID).Error; err != nil {
		// Crear nuevo
		newClient := db.Client{
			ID:        dto.ClienteID,
			Nombre:    dto.ClienteNombre,
			Direccion: dto.ClienteDireccion,
			Email:     dto.ClienteEmail,
			Telefono:  dto.ClienteTelefono,
			TipoID:    "05", // Default
		}
		tx.Create(&newClient)
	} else {
		// Actualizar existente
		cliente.Nombre = dto.ClienteNombre
		cliente.Direccion = dto.ClienteDireccion
		cliente.Email = dto.ClienteEmail
		cliente.Telefono = dto.ClienteTelefono
		tx.Save(&cliente)
	}

	return tx.Commit().Error
}

// GetQuotations obtiene el historial paginado
func (s *QuotationService) GetQuotations(page, pageSize int) ([]db.QuotationDTO, int64) {
	var quotations []db.Quotation
	var total int64
	offset := (page - 1) * pageSize

	db.GetDB().Model(&db.Quotation{}).Count(&total)
	db.GetDB().Order("created_at desc").Limit(pageSize).Offset(offset).Find(&quotations)

	var dtos []db.QuotationDTO
	for _, q := range quotations {
		dtos = append(dtos, db.QuotationDTO{
			ID:            q.ID,
			Secuencial:    q.Secuencial,
			FechaEmision:  q.FechaEmision.Format("02/01/2006"),
			ClienteNombre: q.ClienteNombre,
			Total:         q.Total,
			Estado:        q.Estado,
		})
	}
	return dtos, total
}

func (s *QuotationService) GetPDF(id uint) ([]byte, error) {
	var q db.Quotation
	if err := db.GetDB().First(&q, id).Error; err != nil {
		return nil, err
	}
	if len(q.PDFBytes) == 0 {
		return nil, fmt.Errorf("no existe PDF para esta cotización")
	}
	return q.PDFBytes, nil
}

func (s *QuotationService) ConvertToInvoice(id uint) (*db.FacturaDTO, error) {
	var q db.Quotation
	if err := db.GetDB().First(&q, id).Error; err != nil {
		return nil, err
	}
	
	// Obtener items
	var items []db.QuotationItem
	db.GetDB().Where("quotation_id = ?", q.ID).Find(&items)
	
	var invoiceItems []db.InvoiceItem
	for _, item := range items {
		invoiceItems = append(invoiceItems, db.InvoiceItem{
			Codigo:        item.ProductoSKU,
			Nombre:        item.Nombre,
			Cantidad:      item.Cantidad,
			Precio:        item.PrecioUnitario,
			PorcentajeIVA: item.PorcentajeIVA,
			// Inferir código IVA simple
			CodigoIVA:     "2", // Default simple, idealmente guardar el código exacto
		})
	}
	
	dto := &db.FacturaDTO{
		ClienteID:        q.ClienteID,
		ClienteNombre:    q.ClienteNombre,
		ClienteDireccion: q.ClienteDireccion,
		ClienteEmail:     q.ClienteEmail,
		ClienteTelefono:  q.ClienteTelefono,
		Observacion:      "Basado en Cotización " + q.Secuencial,
		Items:            invoiceItems,
	}
	
	return dto, nil
}
