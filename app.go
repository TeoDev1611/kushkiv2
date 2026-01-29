package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"kushkiv2/internal/db"
	"kushkiv2/internal/service"
	"kushkiv2/pkg/crypto"
	"kushkiv2/pkg/logger"
	"kushkiv2/pkg/util"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	goruntime "runtime"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed internal/mobile/static
var mobileAssets embed.FS

// App struct is the main application controller.
// It handles the lifecycle, dependency injection of services, and exposes methods to the frontend.
type App struct {
	ctx              context.Context
	invoiceService   *service.InvoiceService
	reportService    *service.ReportService
	syncService      *service.SyncService
	cloudService     *service.CloudService
	mailService      *service.MailService
	quotationService *service.QuotationService
	searchService    *service.SearchService
	chartService     *service.ChartService
	taxService       *service.TaxService
	productService   *service.ProductService
	clientService    *service.ClientService

	// Satellite Server
	satelliteToken string
	serverPort     string
}

// NewApp creates a new App application struct
func NewApp() *App {
	for _, arg := range os.Args {
		if arg == "--kushki-debug" {
			logger.DebugMode = true
			break
		}
	}

	return &App{
		invoiceService:   service.NewInvoiceService(),
		reportService:    service.NewReportService(),
		syncService:      service.NewSyncService(),
		cloudService:     service.NewCloudService(),
		mailService:      service.NewMailService(),
		quotationService: service.NewQuotationService(),
		searchService:    service.NewSearchService(),
		chartService:     service.NewChartService(),
		taxService:       service.NewTaxService(),
		productService:   service.NewProductService(),
		clientService:    service.NewClientService(),
		serverPort:       "8085", // Default port
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	
	// Generate Satellite Token (Simple 6-digit PIN for ease of typing on mobile)
	rand.Seed(time.Now().UnixNano())
	a.satelliteToken = fmt.Sprintf("%06d", rand.Intn(1000000))
	logger.Debug("Satellite Token: %s", a.satelliteToken)

	a.startLicenseHeartbeat()
	a.syncService.StartWorker()
	
	// Start Local API Server
	go a.startLocalServer()
}

// --- LOCAL SERVER (SATELLITE) ---

func (a *App) startLocalServer() {
	e := echo.New()
	e.HideBanner = true
	
	// Middleware
	if logger.DebugMode {
		e.Use(middleware.Logger()) // Log requests only in debug mode
	}
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	// e.Use(middleware.Gzip()) // Disabled temporarily to debug latency

	// API Group (Authenticated)
	api := e.Group("/api")
	api.Use(a.authMiddlewareEcho)
	
	api.GET("/inventory", a.handleGetInventoryEcho)
	api.POST("/stock", a.handleUpdateStockEcho)
	api.POST("/pos/scan", a.handlePOSScan)
	api.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Static Assets (Mobile App)
	staticFS, err := fs.Sub(mobileAssets, "internal/mobile/static")
	if err != nil {
		logger.Error("Error loading embedded mobile assets: %v", err)
	} else {
		// Serve index.html for root, and other files normally
		// Using StaticFS with http.FS adapter
		e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(staticFS))))
	}

	logger.Debug("Starting Local Server (Echo) on 0.0.0.0:%s...", a.serverPort)
	logger.Debug("IMPORTANTE: Si no conecta desde el celular, verifique su Firewall (sudo ufw allow %s)", a.serverPort)
	
	if err := e.Start("0.0.0.0:" + a.serverPort); err != nil {
		logger.Error("Error starting local server: %v", err)
	}
}

func (a *App) authMiddlewareEcho(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("X-Kushki-Token")
		// Allow status check without auth (optional, but good for connectivity test)
		if c.Path() == "/api/status" {
			return next(c)
		}
		
		if token != a.satelliteToken {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}
		return next(c)
	}
}

func (a *App) handleGetInventoryEcho(c echo.Context) error {
	products := a.GetProducts()
	return c.JSON(http.StatusOK, products)
}

type StockUpdateRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
	Type     string `json:"type"` 
}

func (a *App) handleUpdateStockEcho(c echo.Context) error {
	var req StockUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	var product db.Product
	if err := db.GetDB().First(&product, "sku = ?", req.SKU).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	if req.Type == "set" {
		product.Stock = req.Quantity
	} else {
		product.Stock += req.Quantity
	}

	// Use UpdateColumn to only update the stock field and avoid Unique Constraint violations 
	// on other fields (like Barcode="") if they haven't changed but trigger validation.
	if err := db.GetDB().Model(&product).Update("stock", product.Stock).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Notify Frontend
	runtime.EventsEmit(a.ctx, "inventory-updated", product)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"newStock": product.Stock,
	})
}

type POSScanRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"` // Optional, defaults to 1 if 0
}

func (a *App) handlePOSScan(c echo.Context) error {
	var req POSScanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	if req.Quantity == 0 {
		req.Quantity = 1
	}

	// Buscar el producto para enviar info completa (opcional, pero útil)
	var product db.Product
	if err := db.GetDB().First(&product, "sku = ? OR barcode = ?", req.SKU, req.SKU).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
	}

	// Enviar evento a Desktop
	// El frontend escuchará "pos-scan-event" y lo añadirá al carrito
	runtime.EventsEmit(a.ctx, "pos-scan-event", map[string]interface{}{
		"sku":      product.SKU,
		"quantity": req.Quantity,
		"product":  product,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Sent to POS",
		"product": product.Name,
	})
}

type SatelliteConnectionDTO struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Token string `json:"token"`
	URL   string `json:"url"`
}

func (a *App) GetSatelliteConnectionInfo() SatelliteConnectionDTO {
	ip := a.getLocalIP()
	return SatelliteConnectionDTO{
		IP:    ip,
		Port:  a.serverPort,
		Token: a.satelliteToken,
		URL:   fmt.Sprintf("http://%s:%s/?token=%s", ip, a.serverPort, a.satelliteToken),
	}
}

func (a *App) getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	
	var candidates []string

	for _, address := range addrs {
		// Verificar que sea IPv4 y no Loopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				// Prioridad 1: 192.168.x.x (Redes domésticas comunes)
				if strings.HasPrefix(ip, "192.168.") {
					return ip
				}
				// Prioridad 2: 10.x.x.x (Redes corporativas comunes)
				if strings.HasPrefix(ip, "10.") {
					candidates = append(candidates, ip)
					continue
				}
				// Prioridad 3: 172.x.x.x
				if strings.HasPrefix(ip, "172.") {
					candidates = append(candidates, ip)
					continue
				}
				// Otros (ej: IPs públicas o raras)
				candidates = append(candidates, ip)
			}
		}
	}

	// Si no encontramos una 192.168, devolvemos la primera de las candidatas (si hay)
	if len(candidates) > 0 {
		return candidates[0]
	}

	return "127.0.0.1"
}

// startLicenseHeartbeat verifica la validez de la licencia cada 4 horas.
func (a *App) startLicenseHeartbeat() {
	go func() {
		for {
			time.Sleep(4 * time.Hour) // Verificar cada 4 horas

			logger.Debug("Iniciando verificación de licencia (Heartbeat)...")

			var config db.EmisorConfig
			if err := db.GetDB().First(&config).Error; err != nil || config.LicenseKey == "" {
				logger.Debug("Heartbeat saltado: No hay licencia configurada.")
				continue
			}

			// Re-validar licencia con el backend
			resp, err := a.cloudService.ActivateLicense(config.LicenseKey)
			if err != nil {
				logger.Debug("Heartbeat Warning: No se pudo verificar la licencia: %v", err)
				// fmt.Printf("Heartbeat Warning: No se pudo verificar la licencia: %v\n", err)
				// Si el error es crítico (ej: 403), podríamos notificar al frontend para bloquear
				// a.NotifyFrontend("error", "Error verificando licencia. Reinicie la aplicación.")
				continue
			}

			// Actualizar token si es válido
			if resp.Token != "" {
				config.LicenseToken = resp.Token
				db.GetDB().Save(&config)
				logger.Debug("Heartbeat: Licencia verificada y token renovado correctamente.")
			}
		}
	}()
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

	// Validar que exista el token
	if config.LicenseToken == "" {
		return false
	}

	// Validar firma del Token JWT almacenado
	if err := crypto.VerifyLicenseToken(config.LicenseToken); err != nil {
		fmt.Printf("Error validando licencia: %v\n", err)
		return false
	}

	return config.LicenseKey != ""
}

// ActivateLicense intenta activar una licencia con el backend.
func (a *App) ActivateLicense(key string) string {
	if key == "" {
		return "Error: La clave de licencia no puede estar vacía"
	}

	// Validar formato KSH-XXXX-XXXX-XXXX
	re := regexp.MustCompile(`^KSH-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}$`)
	if !re.MatchString(key) {
		return "Error: Formato de licencia inválido. Debe ser KSH-XXXX-XXXX-XXXX"
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

// ExportMasterReport genera un reporte Excel completo de toda la base de datos.
func (a *App) ExportMasterReport() string {
	data, err := a.reportService.GenerateMasterReportExcel()
	if err != nil {
		return fmt.Sprintf("Error generando reporte: %v", err)
	}

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: fmt.Sprintf("Resumen_General_%s.xlsx", time.Now().Format("20060102")),
		Title:           "Guardar Reporte Maestro",
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

	return "Reporte maestro exportado exitosamente"
}

// GetTopProducts devuelve los productos más vendidos para gráficos.
func (a *App) GetTopProducts() []service.TopProduct {
	products, err := a.reportService.GetTopProducts(5)
	if err != nil {
		logger.Error("Error obteniendo top productos: %v", err)
		return []service.TopProduct{}
	}
	return products
}

// GetDashboardStats calcula los KPIs para un rango de fechas específico.
// Utiliza goroutines para realizar consultas a la base de datos en paralelo y mejorar la respuesta.
func (a *App) GetDashboardStats(startStr, endStr string) DashboardStats {
	var stats DashboardStats
	start, _ := time.Parse("2006-01-02", startStr)
	end, _ := time.Parse("2006-01-02", endStr)
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	var facturas []db.Factura
	// Cargamos TODAS las facturas del periodo para procesarlas manualmente (Más seguro que SQL puro en SQLite)
	err := db.GetDB().Where("fecha_emision >= ? AND fecha_emision <= ?", start, end).Find(&facturas).Error
	if err != nil {
		logger.Error("Error cargando facturas para stats: %v", err)
		return stats
	}

	stats.TotalFacturas = int64(len(facturas))
	trendMap := make(map[string]float64)

		for _, f := range facturas {
			// MOTOR DE SUMA ULTRA-FLEXIBLE: 
			// Si tiene un valor positivo, lo contamos como venta. 
			// Esto arregla el problema de las facturas con estado vacío.
			if f.Total > 0 {
				stats.TotalVentas += f.Total
				
				// También lo sumamos a la tendencia
				day := f.FechaEmision.Format("2006-01-02")
				trendMap[day] += f.Total
			}
	
			// Para el conteo de pendientes, solo si el estado es algo fallido explícito
			estado := strings.ToUpper(strings.TrimSpace(f.EstadoSRI))
			if estado != "AUTORIZADO" && estado != "ANULADO" && estado != "" {
				stats.Pendientes++
			}
		}
	stats.SRIOnline = true

	// Generar puntos de tendencia
	days := int(end.Sub(start).Hours() / 24)
	if days > 31 {
		days = 31
	}
	for i := 0; i <= days; i++ {
		day := start.AddDate(0, 0, i).Format("2006-01-02")
		stats.SalesTrend = append(stats.SalesTrend, DailySale{Date: day, Total: trendMap[day]})
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

	// 4. ENVIAR CORREO (SOLO SMTP LOCAL)
	if data.ClienteEmail != "" && len(factura.PDFRIDE) > 0 {
		config := a.GetEmisorConfig()

		go func(email, sec string, pdf []byte, conf *db.EmisorConfigDTO) {
			// Validar configuración SMTP
			if conf == nil || conf.SMTPHost == "" || conf.SMTPPort == 0 || conf.SMTPUser == "" {
				a.NotifyFrontend("error", fmt.Sprintf("No se envió el correo a %s: Servidor SMTP no configurado.", email))
				return
			}

			smtpConfig := service.SMTPConfig{
				Host:     conf.SMTPHost,
				Port:     conf.SMTPPort,
				User:     conf.SMTPUser,
				Password: conf.SMTPPassword,
			}

			err := a.mailService.SendInvoiceEmail(smtpConfig, email, conf.RazonSocial, pdf, sec)

			logEntry := db.MailLog{
				FacturaClave: factura.ClaveAcceso,
				Email:        email,
				Fecha:        time.Now(),
			}

			if err != nil {
				fmt.Printf("Error SMTP local: %v\n", err)
				a.NotifyFrontend("error", fmt.Sprintf("Error enviando correo a %s: %v", email, err))
				logEntry.Estado = "FAILED"
				logEntry.Mensaje = err.Error()
			} else {
				a.NotifyFrontend("success", fmt.Sprintf("Correo enviado a %s exitosamente.", email))
				logEntry.Estado = "SUCCESS"
				logEntry.Mensaje = "Enviado correctamente"
			}
			db.GetDB().Create(&logEntry)
		}(data.ClienteEmail, factura.Secuencial, factura.PDFRIDE, config)
	}

	return fmt.Sprintf("Éxito: Factura %s emitida con clave %s", data.Secuencial, data.ClaveAcceso)
}

// ResendInvoiceEmail reenvía una factura usando SMTP local.
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
	config := a.GetEmisorConfig()

	go func(email, sec string, pdf []byte, conf *db.EmisorConfigDTO) {
		// Validar configuración SMTP
		if conf == nil || conf.SMTPHost == "" || conf.SMTPPort == 0 || conf.SMTPUser == "" {
			a.NotifyFrontend("error", fmt.Sprintf("No se reenvió a %s: Servidor SMTP no configurado.", email))
			return
		}

		smtpConfig := service.SMTPConfig{
			Host:     conf.SMTPHost,
			Port:     conf.SMTPPort,
			User:     conf.SMTPUser,
			Password: conf.SMTPPassword,
		}

		err := a.mailService.SendInvoiceEmail(smtpConfig, email, conf.RazonSocial, pdf, sec)

		logEntry := db.MailLog{
			FacturaClave: claveAcceso,
			Email:        email,
			Fecha:        time.Now(),
		}

		if err != nil {
			a.NotifyFrontend("error", fmt.Sprintf("Error SMTP: %v", err))
			logEntry.Estado = "FAILED"
			logEntry.Mensaje = err.Error()
		} else {
			a.NotifyFrontend("success", fmt.Sprintf("Factura reenviada a %s", email))
			logEntry.Estado = "SUCCESS"
			logEntry.Mensaje = "Reenviado correctamente"
		}
		db.GetDB().Create(&logEntry)
	}(cliente.Email, factura.Secuencial, factura.PDFRIDE, config)

	return "Procesando envío..."
}

// GetMailLogs devuelve el historial de correos.
func (a *App) GetMailLogs() []db.MailLogDTO {
	var logs []db.MailLog
	db.GetDB().Order("fecha desc").Limit(50).Find(&logs)

	var dtos []db.MailLogDTO
	for _, l := range logs {
		dtos = append(dtos, db.MailLogDTO{
			ID:           l.ID,
			FacturaClave: l.FacturaClave,
			Email:        l.Email,
			Estado:       l.Estado,
			Mensaje:      l.Mensaje,
			Fecha:        l.Fecha.Format("02/01/2006 15:04:05"),
		})
	}
	return dtos
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

type BackupDTO struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Date string `json:"date"`
	Path string `json:"path"`
}

type DailySale struct {
	Date  string  `json:"date"`
	Total float64 `json:"total"`
}

type DashboardStats struct {
	TotalFacturas int64       `json:"totalFacturas"`
	TotalVentas   float64     `json:"totalVentas"`
	Pendientes    int         `json:"pendientes"`
	SRIOnline     bool        `json:"sriOnline"`
	SalesTrend    []DailySale `json:"salesTrend"`
}

type FacturasResponse struct {
	Total int64                 `json:"total"`
	Data  []db.FacturaResumenDTO `json:"data"`
}

type QuotationListResponse struct {
	Total int64             `json:"total"`
	Data  []db.QuotationDTO `json:"data"`
}

// GetBackups lista los archivos .zip en la carpeta de respaldos.
func (a *App) GetBackups() []BackupDTO {
	config := a.GetEmisorConfig()
	if config == nil {
		return []BackupDTO{}
	}

	backupDir := config.StoragePath
	if backupDir == "" {
		backupDir, _ = os.UserHomeDir()
	} else {
		backupDir = filepath.Join(backupDir, "Backups")
	}

	files, err := os.ReadDir(backupDir)
	if err != nil {
		return []BackupDTO{}
	}

	var backups []BackupDTO
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".zip" {
			info, _ := f.Info()
			sizeMB := float64(info.Size()) / 1024 / 1024

			backups = append(backups, BackupDTO{
				Name: f.Name(),
				Size: fmt.Sprintf("%.2f MB", sizeMB),
				Date: info.ModTime().Format("02/01/2006 15:04"),
				Path: filepath.Join(backupDir, f.Name()),
			})
		}
	}
	// Ordenar por fecha (más reciente primero) es ideal, pero lo haremos en frontend o aquí invertimos el loop si fuera slice simple.
	// Por simplicidad, retornamos tal cual y el frontend puede ordenar o mostramos los últimos.
	return backups
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
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
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
	decryptedSMTP, _ := crypto.Decrypt(config.SMTPPassword)

	return &db.EmisorConfigDTO{
		RUC:                config.RUC,
		RazonSocial:        config.RazonSocial,
		P12Path:            config.P12Path,
		P12Password:        decryptedPass,
		Ambiente:           config.Ambiente,
		Estab:              config.Estab,
		PtoEmi:             config.PtoEmi,
		Obligado:           config.Obligado,
		ContribuyenteRimpe: config.ContribuyenteRimpe,
		AgenteRetencion:    config.AgenteRetencion,
		StoragePath:        config.StoragePath,
		LogoPath:           config.LogoPath,
		PDFTheme:           config.PDFTheme,
		SMTPHost:           config.SMTPHost,
		SMTPPort:           config.SMTPPort,
		SMTPUser:           config.SMTPUser,
		SMTPPassword:       decryptedSMTP,
	}
}

// SaveEmisorConfig guarda la configuración.
func (a *App) SaveEmisorConfig(dto db.EmisorConfigDTO) string {
	if dto.StoragePath != "" {
		if _, err := os.Stat(dto.StoragePath); os.IsNotExist(err) {
			return "Error: La ruta de almacenamiento no existe"
		}
	}

	// Validación RUC (SRI requiere 13 dígitos)
	if len(dto.RUC) != 13 {
		return "Error: El RUC debe tener exactamente 13 dígitos"
	}

	// Validación SMTP
	if strings.Contains(dto.SMTPHost, "smpt") {
		return "Error: Posible error de escritura en servidor SMTP. ¿Quiso decir 'smtp'?"
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

	encryptedSMTP := ""
	if dto.SMTPPassword != "" {
		encryptedSMTP, err = crypto.Encrypt(dto.SMTPPassword)
		if err != nil {
			return fmt.Sprintf("Error al cifrar contraseña SMTP: %v", err)
		}
	}

	existing.RUC = dto.RUC
	existing.RazonSocial = dto.RazonSocial
	existing.P12Path = dto.P12Path
	existing.P12Password = encryptedPass
	existing.Ambiente = dto.Ambiente
	existing.Estab = dto.Estab
	existing.PtoEmi = dto.PtoEmi
	existing.Obligado = dto.Obligado
	existing.ContribuyenteRimpe = dto.ContribuyenteRimpe
	existing.AgenteRetencion = dto.AgenteRetencion
	existing.StoragePath = dto.StoragePath
	existing.LogoPath = dto.LogoPath
	existing.PDFTheme = dto.PDFTheme

	existing.SMTPHost = dto.SMTPHost
	existing.SMTPPort = dto.SMTPPort
	existing.SMTPUser = dto.SMTPUser
	existing.SMTPPassword = encryptedSMTP

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

// SelectAndSaveLogo abre diálogo para imagen y la procesa.
func (a *App) SelectAndSaveLogo() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Seleccionar Logo de Empresa",
		Filters: []runtime.FileFilter{
			{DisplayName: "Imágenes", Pattern: "*.png;*.jpg;*.jpeg"},
		},
	})
	if err != nil || selection == "" {
		return ""
	}

	// Definir dónde guardar (Carpeta de usuario/.kushki/assets o local)
	// Usaremos la StoragePath definida o un default
	config := a.GetEmisorConfig()
	targetDir := "assets"
	if config != nil && config.StoragePath != "" {
		targetDir = filepath.Join(config.StoragePath, "assets")
	} else {
		// Fallback local
		cwd, _ := os.Getwd()
		targetDir = filepath.Join(cwd, "assets")
	}

	finalPath, err := util.ProcessAndSaveLogo(selection, targetDir)
	if err != nil {
		logger.Error("Error procesando logo: %v", err)
		return "Error procesando imagen"
	}

	return finalPath
}

// TestSMTPConnection verifica si las credenciales de correo funcionan.
func (a *App) TestSMTPConnection(dto db.EmisorConfigDTO) string {
	smtpConfig := service.SMTPConfig{
		Host:     dto.SMTPHost,
		Port:     dto.SMTPPort,
		User:     dto.SMTPUser,
		Password: dto.SMTPPassword,
	}

	// Intentar enviar un correo de prueba al mismo usuario
	err := a.mailService.SendTestEmail(smtpConfig, dto.SMTPUser)
	if err != nil {
		return fmt.Sprintf("Error de conexión: %v", err)
	}

	return "Éxito: Correo de prueba enviado a " + dto.SMTPUser
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
	// OPTIMIZACIÓN: Limitamos a 50 resultados para evitar congelar la UI con grandes volúmenes de datos.
	db.GetDB().Where("nombre LIKE ? OR id LIKE ?", likeQuery, likeQuery).Limit(50).Find(&clients)
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
		expiryStr := ""
		if p.ExpiryDate != nil {
			expiryStr = p.ExpiryDate.Format("2006-01-02")
		}
		dtos = append(dtos, db.ProductDTO{
			SKU:           p.SKU,
			Name:          p.Name,
			Price:         p.Price,
			Stock:         p.Stock,
			TaxCode:       strconv.Itoa(p.TaxCode),
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

func (a *App) SearchProducts(query string) []db.ProductDTO {
	res, err := a.searchService.FuzzySearchProducts(query)
	if err != nil {
		logger.Error("Error fuzzy product: %v", err)
		return []db.ProductDTO{}
	}
	return res
}

func (a *App) SaveProduct(dto db.ProductDTO) string {
	taxCodeInt, err := strconv.Atoi(dto.TaxCode)
	if err != nil {
		return fmt.Sprintf("Error: Código de impuesto inválido: %v", err)
	}

	var expiryDate *time.Time
	if dto.ExpiryDate != "" {
		parsed, err := time.Parse("2006-01-02", dto.ExpiryDate)
		if err == nil {
			expiryDate = &parsed
		}
	}

	var existing db.Product
	result := db.GetDB().First(&existing, "sku = ?", dto.SKU)
	if result.Error == nil {
		existing.Name = dto.Name
		existing.Price = dto.Price
		existing.Stock = dto.Stock
		existing.TaxCode = taxCodeInt
		existing.TaxPercentage = dto.TaxPercentage
		existing.Barcode = dto.Barcode
		existing.AuxiliaryCode = dto.AuxiliaryCode
		existing.MinStock = dto.MinStock
		existing.ExpiryDate = expiryDate
		existing.Location = dto.Location

		if err := db.GetDB().Save(&existing).Error; err != nil {
			return fmt.Sprintf("Error actualizando producto: %v", err)
		}
	} else {
		newProd := db.Product{
			SKU:           dto.SKU,
			Name:          dto.Name,
			Price:         dto.Price,
			Stock:         dto.Stock,
			TaxCode:       taxCodeInt,
			TaxPercentage: dto.TaxPercentage,
			Barcode:       dto.Barcode,
			AuxiliaryCode: dto.AuxiliaryCode,
			MinStock:      dto.MinStock,
			ExpiryDate:    expiryDate,
			Location:      dto.Location,
		}
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

// ImportProductsCSV permite al usuario seleccionar un archivo CSV e importar productos masivamente.
func (a *App) ImportProductsCSV() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Seleccionar Archivo CSV de Productos",
		Filters: []runtime.FileFilter{
			{DisplayName: "Archivos CSV", Pattern: "*.csv"},
		},
	})
	if err != nil || selection == "" {
		return "Cancelado"
	}

	file, err := os.Open(selection)
	if err != nil {
		return fmt.Sprintf("Error abriendo archivo: %v", err)
	}
	defer file.Close()

	count, err := a.productService.ImportProductsFromCSV(file)
	if err != nil {
		return fmt.Sprintf("Error importando productos: %v", err)
	}

	return fmt.Sprintf("Éxito: Se importaron/actualizaron %d productos", count)
}

// ImportClientsCSV permite al usuario seleccionar un archivo CSV e importar clientes masivamente.
func (a *App) ImportClientsCSV() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Seleccionar Archivo CSV de Clientes",
		Filters: []runtime.FileFilter{
			{DisplayName: "Archivos CSV", Pattern: "*.csv"},
		},
	})
	if err != nil || selection == "" {
		return "Cancelado"
	}

	file, err := os.Open(selection)
	if err != nil {
		return fmt.Sprintf("Error abriendo archivo: %v", err)
	}
	defer file.Close()

	count, err := a.clientService.ImportClientsFromCSV(file)
	if err != nil {
		return fmt.Sprintf("Error importando clientes: %v", err)
	}

	return fmt.Sprintf("Éxito: Se importaron/actualizaron %d clientes", count)
}

// --- GESTIÓN DE COTIZACIONES ---

func (a *App) GetNextQuotationSecuencial() string {
	sec, err := a.quotationService.GetNextSecuencial()
	if err != nil {
		return "000000001"
	}
	return sec
}

func (a *App) CreateQuotation(dto db.QuotationDTO) string {
	err := a.quotationService.CreateQuotation(&dto)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return "Cotización creada exitosamente"
}

func (a *App) GetQuotations(page, pageSize int) QuotationListResponse {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	dtos, total := a.quotationService.GetQuotations(page, pageSize)
	return QuotationListResponse{Total: total, Data: dtos}
}

func (a *App) OpenQuotationPDF(id uint) string {
	pdfBytes, err := a.quotationService.GetPDF(id)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	tmpDir := os.TempDir()
	fileName := fmt.Sprintf("COT-%d.pdf", id)
	filePath := filepath.Join(tmpDir, fileName)

	if err := os.WriteFile(filePath, pdfBytes, 0644); err != nil {
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

func (a *App) ConvertQuotationToInvoice(id uint) *db.FacturaDTO {
	dto, err := a.quotationService.ConvertToInvoice(id)
	if err != nil {
		logger.Error("Error convirtiendo cotización: %v", err)
		return nil
	}
	return dto
}

// --- NUEVOS MÉTODOS (Fuzzy Search & Charts) ---

func (a *App) GetVATSummary(startStr, endStr string) service.TaxSummary {
	start, _ := time.Parse("2006-01-02", startStr)
	end, _ := time.Parse("2006-01-02", endStr)
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return a.taxService.GetVATSummary(start, end)
}

func (a *App) SearchInvoicesSmart(query string) []db.FacturaResumenDTO {
	res, err := a.searchService.FuzzySearchInvoices(query)
	if err != nil {
		logger.Error("Error fuzzy search: %v", err)
		return []db.FacturaResumenDTO{}
	}
	return res
}

type ChartsDTO struct {
	RevenueBar string `json:"revenueBar"`
	ClientsPie string `json:"clientsPie"`
}

func (a *App) GetStatisticsCharts() ChartsDTO {
	bar, _ := a.chartService.GenerateRevenueChart()
	pie, _ := a.chartService.GenerateClientsPie()
	return ChartsDTO{
		RevenueBar: bar,
		ClientsPie: pie,
	}
}
