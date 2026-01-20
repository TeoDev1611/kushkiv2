package service

import (
	"fmt"
	"kushkiv2/internal/db"
	"kushkiv2/pkg/sri"
	"sync"
	"time"
)

type SyncLog struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`   // "Envío", "Autorización", "Check"
	Status    string `json:"status"`   // "Success", "Error", "Info"
	Detail    string `json:"detail"`   // Resumen corto
	Request   string `json:"request"`  // Payload completo
	Response  string `json:"response"` // Respuesta completa
}

type SyncService struct {
	sriClient *sri.SRIClient
	logs      []SyncLog
	mu        sync.Mutex
}

func NewSyncService() *SyncService {
	return &SyncService{
		sriClient: sri.NewSRIClient(),
		logs:      make([]SyncLog, 0),
	}
}

// AddLog registra un evento en el historial en memoria (Límite 100).
func (s *SyncService) AddLog(action, status, detail, req, resp string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log := SyncLog{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp: time.Now().Format("15:04:05"),
		Action:    action,
		Status:    status,
		Detail:    detail,
		Request:   req,
		Response:  resp,
	}

	// Insertar al inicio (LIFO)
	s.logs = append([]SyncLog{log}, s.logs...)

	// Limitar a 100 logs
	if len(s.logs) > 100 {
		s.logs = s.logs[:100]
	}
}

func (s *SyncService) GetLogs() []SyncLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Retornar copia para evitar race conditions
	copia := make([]SyncLog, len(s.logs))
	copy(copia, s.logs)
	return copia
}

// StartWorker inicia la sincronización en segundo plano.
func (s *SyncService) StartWorker() {
	go func() {
		for {
			time.Sleep(2 * time.Minute) // Verificar cada 2 minutos
			s.SyncPendingInvoices()
		}
	}()
}

// TriggerSync permite la ejecución manual desde el frontend.
func (s *SyncService) TriggerSync() string {
	go s.SyncPendingInvoices()
	return "Sincronización iniciada en segundo plano..."
}

func (s *SyncService) SyncPendingInvoices() {
	// Verificar conectividad básica primero
	if !s.sriClient.CheckConnectivity() {
		s.AddLog("Conectividad", "Error", "No hay conexión con el SRI", "", "")
		return
	}

	var pending []db.Factura
	// Buscar facturas que no pudieron enviarse por red
	db.GetDB().Where("estado_sri = ?", "PENDIENTE_ENVIO").Find(&pending)

	if len(pending) == 0 {
		return
	}

	s.AddLog("Proceso Batch", "Info", fmt.Sprintf("Procesando %d facturas pendientes...", len(pending)), "", "")

	// WORKER POOL: Límite de concurrencia (ej. 3 hilos simultáneos)
	maxConcurrency := 3
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	for _, f := range pending {
		wg.Add(1)
		sem <- struct{}{} // Adquirir token (bloquea si hay 3 ejecutándose)

		go func(factura db.Factura) {
			defer wg.Done()
			defer func() { <-sem }() // Liberar token

			s.processSingleInvoice(&factura)
		}(f)
	}

	wg.Wait()
	s.AddLog("Proceso Batch", "Success", "Sincronización finalizada", "", "")
}

func (s *SyncService) processSingleInvoice(f *db.Factura) {
	reqLog := fmt.Sprintf("Factura: %s", f.Secuencial)
	
	// Reintentar Envío
	resp, err := s.sriClient.EnviarComprobante(f.XMLFirmado)
	
	if err != nil {
		s.AddLog("Envío SRI", "Error", "Fallo de red al enviar", reqLog, err.Error())
		return
	}

	respStr := fmt.Sprintf("Estado: %s", resp.Estado)
	s.AddLog("Envío SRI", "Success", fmt.Sprintf("Factura %s enviada", f.Secuencial), reqLog, respStr)

	if resp.Estado == "RECIBIDA" {
		f.EstadoSRI = "RECIBIDA"
		f.MensajeError = ""
		
		// Intentar Autorizar
		time.Sleep(1 * time.Second)
		respAuth, errAuth := s.sriClient.AutorizarComprobante(f.ClaveAcceso)
		
		if errAuth == nil {
			authStatus := "Desconocido"
			if len(respAuth.Autorizaciones.Autorizacion) > 0 {
				authStatus = respAuth.Autorizaciones.Autorizacion[0].Estado
			}
			
			s.AddLog("Autorización SRI", "Info", fmt.Sprintf("Estado Auth: %s", authStatus), f.ClaveAcceso, fmt.Sprintf("%+v", respAuth))

			for _, auth := range respAuth.Autorizaciones.Autorizacion {
				if auth.Estado == "AUTORIZADO" {
					f.EstadoSRI = "AUTORIZADO"
					break
				} else {
					f.EstadoSRI = auth.Estado
					f.MensajeError = fmt.Sprintf("Rechazo diferido: %s", auth.Mensajes.Mensaje[0].Mensaje)
				}
			}
		} else {
			s.AddLog("Autorización SRI", "Error", "Error consultando autorización", f.ClaveAcceso, errAuth.Error())
		}
	} else {
		// DEVUELTA
		f.EstadoSRI = resp.Estado
		f.MensajeError = "Devuelta en sincronización diferida."
		s.AddLog("Envío SRI", "Warning", "Factura Devuelta", reqLog, respStr)
	}

	// Guardar nuevo estado (GORM es thread-safe con pool configurado)
	db.GetDB().Save(f)
}
