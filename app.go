package main

import (
	"context"
	"fmt"
	"kushkiv2/internal/db"
	"kushkiv2/internal/service"
	"kushkiv2/pkg/crypto"
	"kushkiv2/pkg/util"
	"os"
	"os/exec"
	"path/filepath"
	runtime "runtime"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	invoiceService *service.InvoiceService
	reportService  *service.ReportService
	syncService    *service.SyncService
	cloudService   *service.CloudService
	// mailService ELIMINADO: Se usa CloudService
}

// DashboardStats contiene las métricas para el frontend
type DashboardStats struct {
	TotalVentas   float64     `json:"totalVentas"`
	TotalFacturas int64       `json:"totalFacturas"`
	Pendientes    int64       `json:"pendientes"`
	SRIOnline     bool        `json:"sriOnline"`
	SalesTrend    []DailySale `json:"salesTrend"`
}

type DailySale struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		invoiceService: service.NewInvoiceService(),
		reportService:  service.NewReportService(),
		syncService:    service.NewSyncService(),
		cloudService:   service.NewCloudService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Iniciar Workers
	a.syncService.StartWorker()
}

// NotifyFrontend envía una señal de toast al frontend desde Go.
func (a *App) NotifyFrontend(tipo, mensaje string) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "toast-notification", map[string]string{
			"type":    tipo, // success, error, info
			"message": mensaje,
		})
	}
}

// --- LICENCIAMIENTO ---

// CheckLicense verifica si el sistema tiene una licencia activa.
func (a *App) CheckLicense() bool {
	var config db.EmisorConfig
	// Verificación básica: Existencia en DB
	result := db.GetDB().First(&config)
	if result.Error != nil {
		return false
	}
	// TODO: En el futuro, validar firma del Token JWT almacenado para mayor seguridad
	return config.LicenseKey != "" && config.LicenseToken != ""
}

// ActivateLicense intenta activar una licencia con el backend.
func (a *App) ActivateLicense(key string) string {
	if key == "" {
		return "Error: La clave de licencia no puede estar vacía"
	}

	// 1. Llamar al Backend
	resp, err := a.cloudService.ActivateLicense(key)
	if err != nil {
		return fmt.Sprintf("Error de activación: %v", err)
	}

	// 2. Guardar Licencia y Token en DB
	var config db.EmisorConfig
	res := db.GetDB().First(&config)
	
	config.LicenseKey = key
	config.LicenseToken = resp.Token // Guardamos el token para autenticación futura
	
	// Asegurar campos mínimos si es la primera vez
	if config.RUC == "" {
		config.RUC = "9999999999999" // Temporal
		config.RazonSocial = "Nuevo Usuario"
	}

	if res.Error != nil {
		db.GetDB().Create(&config)
	} else {
		db.GetDB().Save(&config)
	}

	return fmt.Sprintf("Éxito: %s", resp.Message)
}

// --- REPORTERÍA ---

// ExportSalesExcel permite al usuario guardar el reporte de ventas en Excel.
func (a *App) ExportSalesExcel(startStr, endStr string) string {
	start, _ := time.Parse("2006-01-02", startStr)
	end, _ := time.Parse("2006-01-02", endStr)
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	data, err := a.reportService.GenerateSalesExcel(start, end)
	if err != nil {
		return fmt.Sprintf("Error generando Excel: %v", err)
	}

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: fmt.Sprintf("Reporte_Ventas_%s.xlsx", time.Now().Format("20060102")),
		Title:           "Guardar Reporte de Ventas",
		Filters: []runtime.FileFilter{
			{DisplayName: "Archivos Excel", Pattern: "*.xlsx"},
		},
	})

	if err != nil || path == "" {
		return "Cancelado"
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Sprintf("Error guardando archivo: %v", err)
	}

	return "Reporte exportado exitosamente"
}

// GetTopProducts devuelve los productos más vendidos para gráficos.
func (a *App) GetTopProducts() []service.TopProduct {
	products, _ := a.reportService.GetTopProducts(5)
	return products
}

// GetDashboardStats calcula los KPIs.
func (a *App) GetDashboardStats() DashboardStats {
	var stats DashboardStats
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
	sevenDaysAgo := now.AddDate(0, 0, -6)

	var wg sync.WaitGroup
	wg.Add(4)

	// 1. Total Vendido
	go func() {
		defer wg.Done()
		db.GetDB().Model(&db.Factura{}).
			Where("estado_sri = ? AND fecha_emision BETWEEN ? AND ?", "AUTORIZADO", startOfMonth, endOfMonth).
			Select("COALESCE(SUM(total), 0)").
			Scan(&stats.TotalVentas)
	}()

	// 2. Conteo Total
	go func() {
		defer wg.Done()
		db.GetDB().Model(&db.Factura{}).
			Where("fecha_emision BETWEEN ? AND ?", startOfMonth, endOfMonth).
			Count(&stats.TotalFacturas)
	}()

	// 3. Pendientes
	go func() {
		defer wg.Done()
		db.GetDB().Model(&db.Factura{}).
			Where("estado_sri IN ?", []string{"PENDIENTE", "DEVUELTA", "RECHAZADA", "PENDIENTE_ENVIO"}).
			Count(&stats.Pendientes)
	}()

	// 4. Ping Simulado
	stats.SRIOnline = true

	// 5. Tendencia
	trendMap := make(map[string]float64)
	var trendErr error
	
go func() {
		defer wg.Done()
	
rows, err := db.GetDB().Table("facturas").
			Select("date(fecha_emision) as date, SUM(total) as total").
			Where("estado_sri = ? AND fecha_emision >= ?", "AUTORIZADO", sevenDaysAgo).
			Group("date").
			Order("date ASC").
			Rows()
		
		if err != nil {
			trendErr = err
			return
		}
		defer rows.Close()
		
		for rows.Next() {
			var d string
			var t float64
		
rows.Scan(&d, &t)
			if len(d) >= 10 {
				trendMap[d[:10]] = t
			}
		}
	}()

	wg.Wait()

	if trendErr == nil {
		for i := 0; i < 7; i++ {
			day := sevenDaysAgo.AddDate(0, 0, i).Format("2006-01-02")
			val := 0.0
			if v, ok := trendMap[day]; ok {
				val = v
			}
			stats.SalesTrend = append(stats.SalesTrend, DailySale{Date: day, Total: val})
		}
	}

	return stats
}

// --- SINCRONIZACIÓN ---

func (a *App) GetSyncLogs() []service.SyncLog {
	return a.syncService.GetLogs()
}

func (a *App) TriggerSyncManual() string {
	return a.syncService.TriggerSync()
}

// SelectStoragePath abre el diálogo nativo para elegir carpeta.
func (a *App) SelectStoragePath() string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Seleccionar Carpeta para Guardar Facturas",
	})
	if err != nil {
		return ""
	}
	return selection
}

// saveDocument organiza y guarda archivos físicamente.
func (a *App) saveDocument(secuencial string, fecha time.Time, fileType string, content []byte) error {
	config := a.GetEmisorConfig()
	if config == nil || config.StoragePath == "" {
		return fmt.Errorf("no hay ruta de almacenamiento configurada")
	}

	year := fmt.Sprintf("%d", fecha.Year())
	month := fmt.Sprintf("%02d", fecha.Month())
	fullPath := filepath.Join(config.StoragePath, year, month)

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("error creando directorios: %v", err)
	}

	fileName := fmt.Sprintf("FACTURA-%s.%s", secuencial, fileType)
	finalPath := filepath.Join(fullPath, fileName)

	return os.WriteFile(finalPath, content, 0644)
}

// CreateInvoice expone la funcionalidad de emisión de facturas al frontend.
func (a *App) CreateInvoice(data db.FacturaDTO) string {
	// 1. Emitir
	err := a.invoiceService.EmitirFactura(&data)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	// 2. Recuperar la factura
	var factura db.Factura
	if err := db.GetDB().First(&factura, "secuencial = ?", data.Secuencial).Error; err != nil {
		return "Advertencia: Factura emitida pero no se pudo recuperar para guardar archivos."
	}

	// 3. Guardar Archivos Locales
	if errSave := a.saveDocument(factura.Secuencial, factura.FechaEmision, "xml", factura.XMLFirmado); errSave != nil {
		fmt.Printf("Error guardando XML local: %v\n", errSave)
	}
	if len(factura.PDFRIDE) > 0 {
		if errSave := a.saveDocument(factura.Secuencial, factura.FechaEmision, "pdf", factura.PDFRIDE); errSave != nil {
			fmt.Printf("Error guardando PDF local: %v\n", errSave)
		}
	}

	// 4. ENVIAR CORREO VIA CLOUD API (Asíncrono)
	// Ya no usamos SMTP local. Se delega al CloudService.
	if data.ClienteEmail != "" && len(factura.PDFRIDE) > 0 {
		go func(email, sec string, pdf []byte) {
			filename := fmt.Sprintf("FACTURA-%s.pdf", sec)
			err := a.cloudService.SendPDFReport(email, pdf, filename)
			if err != nil {
				fmt.Printf("Error enviando email cloud: %v\n", err)
				a.NotifyFrontend("error", fmt.Sprintf("Error enviando correo: %v", err))
			} else {
				a.NotifyFrontend("success", fmt.Sprintf("Correo enviado a %s", email))
			}
		}(data.ClienteEmail, factura.Secuencial, factura.PDFRIDE)
	}

	return fmt.Sprintf("Éxito: Factura %s emitida con clave %s", data.Secuencial, data.ClaveAcceso)
}

// ResendInvoiceEmail reenvía una factura usando la API Cloud.
func (a *App) ResendInvoiceEmail(claveAcceso string) string {
	var factura db.Factura
	if err := db.GetDB().First(&factura, "clave_acceso = ?", claveAcceso).Error; err != nil {
		return "Error: Factura no encontrada"
	}

	var cliente db.Client
	if err := db.GetDB().First(&cliente, "id = ?", factura.ClienteID).Error; err != nil || cliente.Email == "" {
		return "Error: El cliente no tiene un correo electrónico configurado"
	}

	// Envío Asíncrono
	go func(email, sec string, pdf []byte) {
		filename := fmt.Sprintf("FACTURA-%s.pdf", sec)
		err := a.cloudService.SendPDFReport(email, pdf, filename)
		if err != nil {
			a.NotifyFrontend("error", fmt.Sprintf("Fallo al reenviar: %v", err))
		} else {
			a.NotifyFrontend("success", fmt.Sprintf("Factura reenviada a %s", email))
		}
	}(cliente.Email, factura.Secuencial, factura.PDFRIDE)

	return "Procesando envío..."
}

// SelectBackupPath abre diálogo para seleccionar carpeta.
func (a *App) SelectBackupPath() string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Seleccionar Carpeta para Respaldos Automáticos",
	})
	if err != nil {
		return ""
	}
	return selection
}

// CreateBackup ejecuta el respaldo.
func (a *App) CreateBackup() error {
	config := a.GetEmisorConfig()
	if config == nil {
		return fmt.Errorf("no hay configuración")
	}

	backupDir := config.StoragePath
	if backupDir == "" {
		backupDir, _ = os.UserHomeDir()
	} else {
		backupDir = filepath.Join(backupDir, "Backups")
		_ = os.MkdirAll(backupDir, 0755)
	}

	filename := fmt.Sprintf("Backup_Kushki_%s.zip", time.Now().Format("20060102_150405"))
	destPath := filepath.Join(backupDir, filename)

	sources := make(map[string]string)
	cwd, _ := os.Getwd()
	dbPath := filepath.Join(cwd, "kushki.db")
	sources[dbPath] = "DB"

	if config.StoragePath != "" {
		sources[config.StoragePath] = "Docs"
	}

	return util.CreateBackupZip(destPath, sources)
}

// GetNextSecuencial devuelve el siguiente número disponible.
func (a *App) GetNextSecuencial() string {
	sec, err := a.invoiceService.GetNextSecuencial()
	if err != nil {
		return "000000001"
	}
	return sec
}

// GetFacturasPaginated devuelve el historial.
func (a *App) GetFacturasPaginated(page int, pageSize int) FacturasResponse {
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 10 }
	offset := (page - 1) * pageSize

	var facturas []db.Factura
	var total int64

	db.GetDB().Model(&db.Factura{}).Count(&total)
	db.GetDB().Order("created_at desc").Limit(pageSize).Offset(offset).Find(&facturas)

	clientIDs := make([]string, 0)
	uniqueIDs := make(map[string]bool)
	for _, f := range facturas {
		if !uniqueIDs[f.ClienteID] {
			uniqueIDs[f.ClienteID] = true
			clientIDs = append(clientIDs, f.ClienteID)
		}
	}

	var clients []db.Client
	if len(clientIDs) > 0 {
		db.GetDB().Select("id, nombre").Where("id IN ?", clientIDs).Find(&clients)
	}

	clientMap := make(map[string]string)
	for _, c := range clients {
		clientMap[c.ID] = c.Nombre
	}
	
	var resumenes []db.FacturaResumenDTO
	for _, f := range facturas {
		nombreCliente := f.ClienteID
		if name, ok := clientMap[f.ClienteID]; ok && name != "" {
			nombreCliente = name
		}

		resumenes = append(resumenes, db.FacturaResumenDTO{
			ClaveAcceso: f.ClaveAcceso,
			Secuencial:  f.Secuencial,
			Fecha:       f.FechaEmision.Format("02/01/2006 15:04"),
			Cliente:     nombreCliente,
			Total:       f.Total,
			Estado:      f.EstadoSRI,
			TienePDF:    len(f.PDFRIDE) > 0,
		})
	}

	if resumenes == nil {
		resumenes = []db.FacturaResumenDTO{}
	}

	return FacturasResponse{
		Total: total,
		Data:  resumenes,
	}
}

// OpenFacturaPDF abre el PDF con el visor del sistema.
func (a *App) OpenFacturaPDF(claveAcceso string) string {
	var factura db.Factura
	if err := db.GetDB().First(&factura, "clave_acceso = ?", claveAcceso).Error; err != nil {
		return "Error: Factura no encontrada"
	}

	if len(factura.PDFRIDE) == 0 {
		return "Error: Esta factura no tiene RIDE generado"
	}

	tmpDir := os.TempDir()
	fileName := fmt.Sprintf("RIDE-%s.pdf", factura.Secuencial)
	filePath := filepath.Join(tmpDir, fileName)

	if err := os.WriteFile(filePath, factura.PDFRIDE, 0644); err != nil {
		return fmt.Sprintf("Error escribiendo archivo temporal: %v", err)
	}

	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filePath)
	default: // linux
		cmd = exec.Command("xdg-open", filePath)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("Error abriendo visor PDF: %v", err)
	}

	return "Abriendo PDF..."
}

// OpenInvoiceFolder abre la carpeta.
func (a *App) OpenInvoiceFolder(claveAcceso string) string {
	var factura db.Factura
	if err := db.GetDB().First(&factura, "clave_acceso = ?", claveAcceso).Error; err != nil {
		return "Error: Factura no encontrada"
	}

	config := a.GetEmisorConfig()
	if config == nil || config.StoragePath == "" {
		return "Error: No se ha configurado una ruta de almacenamiento"
	}

	_ = a.saveDocument(factura.Secuencial, factura.FechaEmision, "xml", factura.XMLFirmado)
	if len(factura.PDFRIDE) > 0 {
		_ = a.saveDocument(factura.Secuencial, factura.FechaEmision, "pdf", factura.PDFRIDE)
	}

	year := fmt.Sprintf("%d", factura.FechaEmision.Year())
	month := fmt.Sprintf("%02d", factura.FechaEmision.Month())
	fullPath := filepath.Join(config.StoragePath, year, month)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Sprintf("Error: La carpeta %s no pudo ser creada", fullPath)
	}

	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "darwin":
		cmd = exec.Command("open", fullPath)
	case "windows":
		cmd = exec.Command("explorer", fullPath)
	default: // linux
		cmd = exec.Command("xdg-open", fullPath)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("Error abriendo carpeta: %v", err)
	}

	return "Abriendo carpeta..."
}

// OpenInvoiceXML abre el XML.
func (a *App) OpenInvoiceXML(claveAcceso string) string {
	var factura db.Factura
	if err := db.GetDB().First(&factura, "clave_acceso = ?", claveAcceso).Error; err != nil {
		return "Error: Factura no encontrada"
	}

	config := a.GetEmisorConfig()
	if config == nil || config.StoragePath == "" {
		return "Error: No se ha configurado una ruta de almacenamiento"
	}

	if err := a.saveDocument(factura.Secuencial, factura.FechaEmision, "xml", factura.XMLFirmado); err != nil {
		return fmt.Sprintf("Error restaurando archivo XML: %v", err)
	}

	year := fmt.Sprintf("%d", factura.FechaEmision.Year())
	month := fmt.Sprintf("%02d", factura.FechaEmision.Month())
	fileName := fmt.Sprintf("FACTURA-%s.xml", factura.Secuencial)
	filePath := filepath.Join(config.StoragePath, year, month, fileName)

	var cmd *exec.Cmd
	switch goruntime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filePath)
	default: // linux
		cmd = exec.Command("xdg-open", filePath)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("Error abriendo XML: %v", err)
	}

	return "Abriendo XML..."
}

// GetEmisorConfig devuelve la configuración.
func (a *App) GetEmisorConfig() *db.EmisorConfigDTO {
	var config db.EmisorConfig
	result := db.GetDB().First(&config)
	if result.Error != nil {
		return nil
	}
	decryptedPass, _ := crypto.Decrypt(config.P12Password)
	// Eliminado desencriptado SMTP

	return &db.EmisorConfigDTO{
		RUC:         config.RUC,
		RazonSocial: config.RazonSocial,
		P12Path:     config.P12Path,
		P12Password: decryptedPass,
		Ambiente:    config.Ambiente,
		Estab:       config.Estab,
		PtoEmi:      config.PtoEmi,
		Obligado:    config.Obligado,
		StoragePath: config.StoragePath,
	}
}

// SaveEmisorConfig guarda la configuración.
func (a *App) SaveEmisorConfig(dto db.EmisorConfigDTO) string {
	if dto.StoragePath != "" {
		if _, err := os.Stat(dto.StoragePath); os.IsNotExist(err) {
			return "Error: La ruta de almacenamiento no existe"
		}
	}

	if err := crypto.ValidateCert(dto.P12Path, dto.P12Password); err != nil {
		return fmt.Sprintf("Error de Validación: %v", err)
	}

	var existing db.EmisorConfig
	result := db.GetDB().First(&existing)
	
	encryptedPass, err := crypto.Encrypt(dto.P12Password)
	if err != nil {
		return fmt.Sprintf("Error al cifrar contraseña firma: %v", err)
	}

	existing.RUC = dto.RUC
	existing.RazonSocial = dto.RazonSocial
	existing.P12Path = dto.P12Path
	existing.P12Password = encryptedPass
	existing.Ambiente = dto.Ambiente
	existing.Estab = dto.Estab
	existing.PtoEmi = dto.PtoEmi
	existing.Obligado = dto.Obligado
	existing.StoragePath = dto.StoragePath

	if result.Error == nil {
		if err := db.GetDB().Save(&existing).Error; err != nil {
			return fmt.Sprintf("Error al actualizar: %v", err)
		}
	} else {
		if err := db.GetDB().Create(&existing).Error; err != nil {
			return fmt.Sprintf("Error al guardar: %v", err)
		}
	}
	return "Configuración guardada exitosamente"
}

// SelectCertificate abre diálogo para .p12
func (a *App) SelectCertificate() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Seleccionar Certificado Digital (.p12)",
		Filters: []runtime.FileFilter{
			{DisplayName: "Certificados P12", Pattern: "*.p12;*.pfx"},
		},
	})
	if err != nil {
		return ""
	}
	return selection
}

// --- GESTIÓN DE CLIENTES ---

func (a *App) GetClients() []db.ClientDTO {
	var clients []db.Client
	db.GetDB().Find(&clients)
	var dtos []db.ClientDTO
	for _, c := range clients {
		dtos = append(dtos, db.ClientDTO{ID: c.ID, TipoID: c.TipoID, Nombre: c.Nombre, Direccion: c.Direccion, Email: c.Email, Telefono: c.Telefono})
	}
	return dtos
}

func (a *App) SearchClients(query string) []db.ClientDTO {
	var clients []db.Client
	likeQuery := "%" + query + "%"
	db.GetDB().Where("nombre LIKE ? OR id LIKE ?", likeQuery, likeQuery).Find(&clients)
	var dtos []db.ClientDTO
	for _, c := range clients {
		dtos = append(dtos, db.ClientDTO{ID: c.ID, TipoID: c.TipoID, Nombre: c.Nombre, Direccion: c.Direccion, Email: c.Email, Telefono: c.Telefono})
	}
	return dtos
}

func (a *App) SaveClient(dto db.ClientDTO) string {
	var existing db.Client
	result := db.GetDB().First(&existing, "id = ?", dto.ID)
	if result.Error == nil {
		existing.TipoID = dto.TipoID
		existing.Nombre = dto.Nombre
		existing.Direccion = dto.Direccion
		existing.Email = dto.Email
		existing.Telefono = dto.Telefono
		if err := db.GetDB().Save(&existing).Error; err != nil {
			return fmt.Sprintf("Error actualizando cliente: %v", err)
		}
	} else {
		newClient := db.Client{ID: dto.ID, TipoID: dto.TipoID, Nombre: dto.Nombre, Direccion: dto.Direccion, Email: dto.Email, Telefono: dto.Telefono}
		if err := db.GetDB().Create(&newClient).Error; err != nil {
			return fmt.Sprintf("Error creando cliente: %v", err)
		}
	}
	return "Cliente guardado exitosamente"
}

func (a *App) DeleteClient(id string) string {
	if err := db.GetDB().Delete(&db.Client{}, "id = ?", id).Error; err != nil {
		return fmt.Sprintf("Error eliminando cliente: %v", err)
	}
	return "Cliente eliminado"
}

// --- GESTIÓN DE PRODUCTOS ---

func (a *App) GetProducts() []db.ProductDTO {
	var products []db.Product
	db.GetDB().Find(&products)
	var dtos []db.ProductDTO
	for _, p := range products {
		dtos = append(dtos, db.ProductDTO{SKU: p.SKU, Name: p.Name, Price: p.Price, Stock: p.Stock, TaxCode: p.TaxCode, TaxPercentage: p.TaxPercentage})
	}
	return dtos
}

func (a *App) SearchProducts(query string) []db.ProductDTO {
	var products []db.Product
	likeQuery := "%" + query + "%"
	db.GetDB().Where("name LIKE ? OR sku LIKE ?", likeQuery, likeQuery).Find(&products)
	var dtos []db.ProductDTO
	for _, p := range products {
		dtos = append(dtos, db.ProductDTO{SKU: p.SKU, Name: p.Name, Price: p.Price, Stock: p.Stock, TaxCode: p.TaxCode, TaxPercentage: p.TaxPercentage})
	}
	return dtos
}

func (a *App) SaveProduct(dto db.ProductDTO) string {
	var existing db.Product
	result := db.GetDB().First(&existing, "sku = ?", dto.SKU)
	if result.Error == nil {
		existing.Name = dto.Name
		existing.Price = dto.Price
		existing.Stock = dto.Stock
		existing.TaxCode = dto.TaxCode
		existing.TaxPercentage = dto.TaxPercentage
		if err := db.GetDB().Save(&existing).Error; err != nil {
			return fmt.Sprintf("Error actualizando producto: %v", err)
		}
	} else {
		newProd := db.Product{SKU: dto.SKU, Name: dto.Name, Price: dto.Price, Stock: dto.Stock, TaxCode: dto.TaxCode, TaxPercentage: dto.TaxPercentage}
		if err := db.GetDB().Create(&newProd).Error; err != nil {
			return fmt.Sprintf("Error creando producto: %v", err)
		}
	}
	return "Producto guardado exitosamente"
}

func (a *App) DeleteProduct(sku string) string {
	if err := db.GetDB().Delete(&db.Product{}, "sku = ?", sku).Error; err != nil {
		return fmt.Sprintf("Error eliminando producto: %v", err)
	}
	return "Producto eliminado"
}