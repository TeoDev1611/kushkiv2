package service

import (
	"fmt"
	"kushkiv2/internal/db"

	"github.com/sahilm/fuzzy"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

// InvoiceSearchItem ayuda a fuzzy a saber qué texto buscar y qué ID retornar
type InvoiceSearchItem struct {
	ClaveAcceso   string
	SearchContent string // Ej: "Juan Perez 001-001-000000123"
}

// InvoiceSource implementa la interfaz fuzzy.Source
type InvoiceSource []InvoiceSearchItem

func (s InvoiceSource) String(i int) string { return s[i].SearchContent }
func (s InvoiceSource) Len() int           { return len(s) }

// FuzzySearchInvoices realiza la búsqueda difusa en memoria
func (s *SearchService) FuzzySearchInvoices(query string) ([]db.FacturaResumenDTO, error) {
	var facturas []db.Factura
	// 1. Cargamos las facturas con sus clientes para tener los nombres
	// Usamos un JOIN para eficiencia
	err := db.GetDB().Table("facturas").
		Select("facturas.*, clients.nombre as cliente_nombre_temp"). // Campo temporal para el nombre
		Joins("left join clients on clients.id = facturas.cliente_id").
		Order("facturas.created_at desc").
		Limit(500). // Sensato para búsqueda en memoria en escritorio
		Find(&facturas).Error

	if err != nil {
		return nil, err
	}

	// 2. Preparamos el origen de datos para fuzzy
	source := make(InvoiceSource, len(facturas))
	facturaMap := make(map[string]db.Factura)

	for i, f := range facturas {
		// Construimos un string que combine datos relevantes: Nombre, RUC/CI, Secuencial, Total, Estado
		// Esto permite búsquedas como "Juan 50.00" o "Autorizado 001"
		searchStr := fmt.Sprintf("%s %s %s %.2f %s", f.ClienteID, f.Secuencial, f.ClaveAcceso, f.Total, f.EstadoSRI)
		
		source[i] = InvoiceSearchItem{
			ClaveAcceso:   f.ClaveAcceso,
			SearchContent: searchStr,
		}
		facturaMap[f.ClaveAcceso] = f
	}

	// 3. Si el query es vacío, devolvemos las últimas 20 facturas por defecto
	if query == "" {
		return s.mapToDTO(facturas[:min(20, len(facturas))]), nil
	}

	// 4. Búsqueda difusa
	matches := fuzzy.FindFrom(query, source)
	
	// 5. Mapear resultados ordenados por relevancia (score)
	var finalInvoices []db.Factura
	for _, match := range matches {
		finalInvoices = append(finalInvoices, facturaMap[source[match.Index].ClaveAcceso])
	}

	return s.mapToDTO(finalInvoices), nil
}

func (s *SearchService) mapToDTO(facturas []db.Factura) []db.FacturaResumenDTO {
	dtos := make([]db.FacturaResumenDTO, 0)
	for _, f := range facturas {
		dtos = append(dtos, db.FacturaResumenDTO{
			ClaveAcceso: f.ClaveAcceso,
			Secuencial:  f.Secuencial,
			Fecha:       f.FechaEmision.Format("02/01/2006 15:04"),
			Cliente:     f.ClienteID,
			Total:       f.Total,
			Estado:      f.EstadoSRI,
			TienePDF:    len(f.PDFRIDE) > 0,
		})
	}
	return dtos
}

func min(a, b int) int {
	if a < b { return a }
	return b
}

// --- CLIENTS ---

type ClientSearchItem struct {
	ID            string
	SearchContent string
}

type ClientSource []ClientSearchItem

func (s ClientSource) String(i int) string { return s[i].SearchContent }
func (s ClientSource) Len() int           { return len(s) }

func (s *SearchService) FuzzySearchClients(query string) ([]db.ClientDTO, error) {
	var clients []db.Client
	// Cargar datos ligeros
	err := db.GetDB().Limit(2000).Find(&clients).Error
	if err != nil {
		return nil, err
	}

	source := make(ClientSource, len(clients))
	clientMap := make(map[string]db.Client)

	for i, c := range clients {
		// Búsqueda por Nombre, CI/RUC, Email
		searchStr := fmt.Sprintf("%s %s %s", c.Nombre, c.ID, c.Email)
		source[i] = ClientSearchItem{
			ID:            c.ID,
			SearchContent: searchStr,
		}
		clientMap[c.ID] = c
	}

	if query == "" {
		return s.mapClientsToDTO(clients[:min(50, len(clients))]), nil
	}

	matches := fuzzy.FindFrom(query, source)
	
	var finalClients []db.Client
	for _, match := range matches {
		finalClients = append(finalClients, clientMap[source[match.Index].ID])
	}

	return s.mapClientsToDTO(finalClients), nil
}

func (s *SearchService) mapClientsToDTO(clients []db.Client) []db.ClientDTO {
	dtos := make([]db.ClientDTO, 0)
	for _, c := range clients {
		dtos = append(dtos, db.ClientDTO{
			ID:        c.ID,
			TipoID:    c.TipoID,
			Nombre:    c.Nombre,
			Direccion: c.Direccion,
			Email:     c.Email,
			Telefono:  c.Telefono,
		})
	}
	return dtos
}

// --- PRODUCTS ---

type ProductSearchItem struct {
	SKU           string
	SearchContent string
}

type ProductSource []ProductSearchItem

func (s ProductSource) String(i int) string { return s[i].SearchContent }
func (s ProductSource) Len() int           { return len(s) }

func (s *SearchService) FuzzySearchProducts(query string) ([]db.ProductDTO, error) {
	var products []db.Product
	err := db.GetDB().Limit(2000).Find(&products).Error
	if err != nil {
		return nil, err
	}

	source := make(ProductSource, len(products))
	productMap := make(map[string]db.Product)

	for i, p := range products {
		// Búsqueda por Nombre, SKU, Precio, Código de Barras y Código Auxiliar
		searchStr := fmt.Sprintf("%s %s %s %s %.2f", p.Name, p.SKU, p.Barcode, p.AuxiliaryCode, p.Price)
		source[i] = ProductSearchItem{
			SKU:           p.SKU,
			SearchContent: searchStr,
		}
		productMap[p.SKU] = p
	}

	if query == "" {
		return s.mapProductsToDTO(products[:min(50, len(products))]), nil
	}

	matches := fuzzy.FindFrom(query, source)
	
	var finalProducts []db.Product
	for _, match := range matches {
		finalProducts = append(finalProducts, productMap[source[match.Index].SKU])
	}

	return s.mapProductsToDTO(finalProducts), nil
}

func (s *SearchService) mapProductsToDTO(products []db.Product) []db.ProductDTO {
	dtos := make([]db.ProductDTO, 0)
	for _, p := range products {
		expiryStr := ""
		if p.ExpiryDate != nil {
			expiryStr = p.ExpiryDate.Format("2006-01-02")
		}

		dtos = append(dtos, db.ProductDTO{
			SKU:           p.SKU,
			Name:          p.Name,
			Price:         p.Price,
			Stock:         p.Stock,
			TaxCode:       fmt.Sprintf("%d", p.TaxCode),
			TaxPercentage: p.TaxPercentage,
			Barcode:       p.Barcode,
			AuxiliaryCode: p.AuxiliaryCode,
			MinStock:      p.MinStock,
			ExpiryDate:    expiryStr,
			Location:      p.Location,
		})
	}
	return dtos
}
