package service

import (
	"fmt"
	"kushkiv2/internal/db"
	"kushkiv2/pkg/crypto"
	"kushkiv2/pkg/logger"
	"kushkiv2/pkg/pdf"
	"kushkiv2/pkg/sri"
	"kushkiv2/pkg/util"
	"kushkiv2/pkg/xml"
	"time"
)

type InvoiceService struct {
	db        *db.EmisorConfig // Cache de configuración (opcional)
	sriClient *sri.SRIClient
}

func NewInvoiceService() *InvoiceService {
	return &InvoiceService{
		sriClient: sri.NewSRIClient(),
	}
}

func (s *InvoiceService) GetNextSecuencial() (string, error) {
	var lastFacturas []db.Factura
	// Usamos Find con Limit 1 para evitar el error "record not found" si está vacía
	db.GetDB().Order("created_at desc").Limit(1).Find(&lastFacturas)
	
	if len(lastFacturas) == 0 {
		// Si no hay facturas, empezamos en 1
		return "000000001", nil
	}
	lastFactura := lastFacturas[0]

	var currentSec int
	fmt.Sscanf(lastFactura.Secuencial, "%d", &currentSec)
	return fmt.Sprintf("%09d", currentSec+1), nil
}

// EmitirFactura coordina el flujo completo de facturación.
func (s *InvoiceService) EmitirFactura(dto *db.FacturaDTO) error {
	// 1. Obtener Configuración del Emisor
	var config db.EmisorConfig
	if err := db.GetDB().First(&config).Error; err != nil {
		return fmt.Errorf("emisor no configurado: %v", err)
	}
	
	// Validar RUC del Emisor (CRÍTICO: debe tener 13 dígitos)
	if len(config.RUC) != 13 {
		return fmt.Errorf("configuración inválida: El RUC del emisor tiene %d dígitos, debe tener 13", len(config.RUC))
	}

	// VALIDACIONES NORMATIVA SRI 2025/2026

	// Regla 6: Validación de productos
	if len(dto.Items) == 0 {
		return fmt.Errorf("error validación: la factura debe tener al menos un ítem")
	}
	for _, item := range dto.Items {
		if len(item.Codigo) == 0 {
			return fmt.Errorf("normativa 2025: todos los ítems deben tener código principal (SKU)")
		}
	}

	// Regla 2: Uso sistema financiero > $1000
	// Calculamos el total preliminar para validar
	var totalValidacion float64
	for _, item := range dto.Items {
		base := util.Round(item.Cantidad*item.Precio, 2)
		impuesto := util.Round(base*(item.PorcentajeIVA/100), 2)
		totalValidacion += base + impuesto
	}

	if config.Ambiente == 2 && totalValidacion >= 1000.00 && dto.FormaPago == "01" {
		return fmt.Errorf("normativa SRI: facturas superiores a $1,000 requieren uso del sistema financiero (no '01')")
	}

	// Regla 3: Consumidor Final > $50
	if dto.ClienteID == "9999999999999" {
		if totalValidacion > 50.00 {
			return fmt.Errorf("normativa SRI: consumidor final no permitido para montos mayores a $50 (se requieren datos reales)")
		}
	}
	
	// Si la forma de pago viene vacía, asignamos "01" por defecto (si cumple reglas)
	if dto.FormaPago == "" {
		dto.FormaPago = "01" 
	}


	// 2. Preparar Datos y Cálculos (Regla 1: IVA Dinámico)
	var detallesXML []xml.Detalle
	var totalConImpuestos []xml.TotalImpuesto
	
	// Mapa para agrupar bases imponibles por código de impuesto
	// Key: CodigoPorcentaje ("0", "2", "4", "5"), Value: BaseImponible
	basesImponibles := make(map[string]struct {
		Base  float64
		Valor float64
		Tarifa float64
	})

	for _, item := range dto.Items {
		precioTotalSinImpuesto := util.Round(item.Cantidad*item.Precio, 2)

		// Crear detalle XML
		detalle := xml.Detalle{
			CodigoPrincipal:        item.Codigo,
			Descripcion:            item.Nombre,
			Cantidad:               item.Cantidad, // Cantidad permitimos hasta 6, no redondeamos agresivamente aquí
			PrecioUnitario:         item.Precio,   // Unitario hasta 6
			Descuento:              0.00,
			PrecioTotalSinImpuesto: precioTotalSinImpuesto,
			Impuestos:              []xml.Impuesto{},
		}

		// Cálculo Impuesto Dinámico
		valorIVA := util.Round(precioTotalSinImpuesto*(item.PorcentajeIVA/100), 2)

		detalle.Impuestos = append(detalle.Impuestos, xml.Impuesto{
			Codigo:           "2", // Siempre 2 para IVA en Ecuador
			CodigoPorcentaje: item.CodigoIVA, // "0", "2", "4", "5"
			Tarifa:           item.PorcentajeIVA,
			BaseImponible:    precioTotalSinImpuesto,
			Valor:            valorIVA,
		})

		// Acumular para totales
		actual := basesImponibles[item.CodigoIVA]
		actual.Base = util.Round(actual.Base+precioTotalSinImpuesto, 2)
		actual.Valor = util.Round(actual.Valor+valorIVA, 2)
		actual.Tarifa = item.PorcentajeIVA
		basesImponibles[item.CodigoIVA] = actual

		detallesXML = append(detallesXML, detalle)
	}

	// Construir resumen de impuestos (TotalConImpuestos)
	var importeTotal, totalSinImpuestos float64
	
	for codigo, datos := range basesImponibles {
		totalConImpuestos = append(totalConImpuestos, xml.TotalImpuesto{
			Codigo:           "2",
			CodigoPorcentaje: codigo,
			BaseImponible:    datos.Base,
			Valor:            datos.Valor,
		})
		importeTotal += datos.Base + datos.Valor
		totalSinImpuestos += datos.Base
	}
	
	// Redondeo final de totales globales
	importeTotal = util.Round(importeTotal, 2)
	totalSinImpuestos = util.Round(totalSinImpuestos, 2)

	// 3. Formateo Estricto SRI (Padding)
	// RECALCULAR SECUENCIAL: Ignoramos el del DTO por ser inseguro (concurrencia)
	// y obtenemos el verdadero siguiente disponible.
	realSec, _ := s.GetNextSecuencial()
	var nEstab, nPtoEmi, nSec int
	fmt.Sscanf(config.Estab, "%d", &nEstab)
	fmt.Sscanf(config.PtoEmi, "%d", &nPtoEmi)
	fmt.Sscanf(realSec, "%d", &nSec)

	estabStr := fmt.Sprintf("%03d", nEstab)
	ptoEmiStr := fmt.Sprintf("%03d", nPtoEmi)
	secuencialStr := fmt.Sprintf("%09d", nSec)

	// Actualizar DTO para reflejar el real usado
	dto.Secuencial = secuencialStr

	// 4. Generar Clave de Acceso (49 dígitos)
	// Algoritmo estándar del SRI: Fecha + Tipo + RUC + Ambiente + Serie + Secuencial + Código + TipoEmisión + DigitoVerificador
	fechaEmision := time.Now()
	fechaStr := fechaEmision.Format("02012006")
	tipoDoc := "01" // Factura
	ruc := config.RUC
	ambiente := fmt.Sprintf("%d", config.Ambiente)
	serie := estabStr + ptoEmiStr
	emision := "1"          // Normal
	codigoNum := "12345678" // Código de seguridad fijo (puede ser aleatorio)

	clavePrevia := fechaStr + tipoDoc + ruc + ambiente + serie + secuencialStr + codigoNum + emision
	digito := util.CalcularDigitoModulo11(clavePrevia)
	claveAcceso := fmt.Sprintf("%s%d", clavePrevia, digito)

	// Determinar dirección matriz (fallback a Razon Social si vacía)
	dirMatriz := config.Direccion
	if dirMatriz == "" {
		dirMatriz = config.RazonSocial
	}

		// 5. Construir XML Completo con campos formateados (REFACTORIZADO: 100% DINÁMICO)

		facturaXML := &xml.FacturaXML{

			Version: "1.1.0",

			ID:      "comprobante",

			InfoTributaria: xml.InfoTributaria{

				Ambiente:        ambiente,

				TipoEmision:     emision,

				RazonSocial:     config.RazonSocial,

				NombreComercial: config.NombreComercial, // Ahora dinámico

				Ruc:             ruc,                     // Variable unificada con ClaveAcceso

				ClaveAcceso:     claveAcceso,

				CodDoc:          tipoDoc,

				Estab:           estabStr,

				PtoEmi:          ptoEmiStr,

				Secuencial:      secuencialStr,

				DirMatriz:       dirMatriz,
				ContribuyenteRimpe: config.ContribuyenteRimpe,
				AgenteRetencion:    config.AgenteRetencion,
			},

			InfoFactura: xml.InfoFactura{

				FechaEmision:                fechaEmision.Format("02/01/2006"),

				DirEstablecimiento:          dirMatriz, // Usamos DirMatriz también aquí por defecto

				ObligadoContabilidad:        "NO",

				TipoIdentificacionComprador: "05", // Cédula por defecto

				RazonSocialComprador:        dto.ClienteNombre,

				IdentificacionComprador:     dto.ClienteID,

				DireccionComprador:          dto.ClienteDireccion,

				TotalSinImpuestos:           totalSinImpuestos,

				TotalDescuento:              0.00,

				TotalConImpuestos:           totalConImpuestos,

				ImporteTotal:                importeTotal,

				Moneda:                      "DOLAR",

				Pagos: []xml.Pago{
					{
						FormaPago:    dto.FormaPago,
						Total:        importeTotal,
						Plazo:        dto.Plazo,
						UnidadTiempo: dto.UnidadTiempo,
					},
				},

			},

			Detalles: detallesXML,

		}

	if config.Obligado {
		facturaXML.InfoFactura.ObligadoContabilidad = "SI"
	}

	// Campos Adicionales (Email, Teléfono, Dirección extra)
	if dto.ClienteEmail != "" {
		facturaXML.InfoAdicional = append(facturaXML.InfoAdicional, xml.CampoAdicional{Nombre: "Email", Value: dto.ClienteEmail})
	}
	if dto.ClienteTelefono != "" {
		facturaXML.InfoAdicional = append(facturaXML.InfoAdicional, xml.CampoAdicional{Nombre: "Telefono", Value: dto.ClienteTelefono})
	}
	if dto.Observacion != "" {
		facturaXML.InfoAdicional = append(facturaXML.InfoAdicional, xml.CampoAdicional{Nombre: "Observacion", Value: dto.Observacion})
	}
	// Si la dirección es larga o se requiere explícitamente en adicionales también:
	if dto.ClienteDireccion != "" {
		facturaXML.InfoAdicional = append(facturaXML.InfoAdicional, xml.CampoAdicional{Nombre: "Direccion", Value: dto.ClienteDireccion})
	}

	xmlData, err := xml.GenerateXML(facturaXML)

	if err != nil {

		return err

	}

	// Calcular subtotales para DB
	var subtotalGravado, subtotalCero, totalIVA float64
	
	if datos, ok := basesImponibles["0"]; ok {
		subtotalCero = datos.Base
	}
	
	// Sumar cualquier otra base gravada (2, 4, 5, etc)
	for codigo, datos := range basesImponibles {
		if codigo != "0" {
			subtotalGravado += datos.Base
			totalIVA += datos.Valor
		}
	}

	// 6. Crear Registro en DB (Factura)
	facturaDB := &db.Factura{
		ClaveAcceso:  claveAcceso,
		Secuencial:   secuencialStr,
		FechaEmision: fechaEmision,
		ClienteID:    dto.ClienteID,
		Total:        importeTotal,
		Subtotal15:   subtotalGravado, // Reutilizamos campo para Base Gravada
		Subtotal0:    subtotalCero,
		IVA:          totalIVA,
		EstadoSRI:    "PENDIENTE",
	}

	// 6. Firmar XML
	// Descifrar contraseña
	p12Pass, err := crypto.Decrypt(config.P12Password)
	if err != nil {
		return fmt.Errorf("error descifrando contraseña de firma: %v", err)
	}

	signer, err := crypto.NewSignerFromFile(config.P12Path, p12Pass)
	if err != nil {
		return fmt.Errorf("error cargando firma: %v", err)
	}

	xmlFirmado, err := signer.SignXML(xmlData)
	if err != nil {
		return fmt.Errorf("error firmando xml: %v", err)
	}

	facturaDB.XMLFirmado = xmlFirmado

	// 7. Enviar al SRI (Recepción)
	respRecepcion, err := s.sriClient.EnviarComprobante(xmlFirmado)
	
	// Manejo de Errores de Red (Contingencia Offline)
	if _, isNetworkError := err.(*sri.NetworkError); isNetworkError {
		facturaDB.EstadoSRI = "PENDIENTE_ENVIO"
		facturaDB.MensajeError = "SRI Offline: Documento guardado para envío posterior."
		// No retornamos error, continuamos para generar PDF y guardar en DB
	} else if err != nil {
		// Error técnico fatal (ej. XML mal formado localmente)
		facturaDB.EstadoSRI = "ERROR_TECNICO"
		facturaDB.MensajeError = err.Error()
	} else {
		// Conexión exitosa, procesar respuesta SRI
		if respRecepcion.Estado == "RECIBIDA" {
			facturaDB.EstadoSRI = "RECIBIDA"
			
			// Esperar un momento antes de pedir autorización (latencia del SRI)
			time.Sleep(2 * time.Second)

			// 8. Solicitar Autorización
			respAuth, errAuth := s.sriClient.AutorizarComprobante(claveAcceso)
			if _, isNetErr := errAuth.(*sri.NetworkError); isNetErr {
				// Recibida pero falló la consulta de autorización
				facturaDB.EstadoSRI = "RECIBIDA" // Se queda así, el worker verificará luego
				facturaDB.MensajeError = "Documento recibido. Verificación de autorización pendiente por red."
			} else if errAuth != nil {
				facturaDB.EstadoSRI = "ERROR_AUTH"
				facturaDB.MensajeError = errAuth.Error()
			} else {
				// Buscar autorización válida
				autorizado := false
				for _, auth := range respAuth.Autorizaciones.Autorizacion {
					if auth.Estado == "AUTORIZADO" {
						facturaDB.EstadoSRI = "AUTORIZADO"
						facturaDB.MensajeError = "" // Limpiar errores previos
						// Aquí podríamos guardar el XML autorizado que devuelve el SRI (tiene fecha y número)
						autorizado = true
						break
					} else {
						// Concatenar mensajes de rechazo
						msg := fmt.Sprintf("[%s]", auth.Estado)
						for _, m := range auth.Mensajes.Mensaje {
							msg += fmt.Sprintf(" %s: %s;", m.Identificador, m.Mensaje)
						}
						facturaDB.MensajeError = msg
						facturaDB.EstadoSRI = auth.Estado
					}
				}
				if !autorizado && facturaDB.EstadoSRI == "RECIBIDA" {
					// Caso raro: Recibida pero sin respuesta clara de autorización
					facturaDB.MensajeError = "Documento recibido pero no se obtuvo respuesta de autorización."
				}
			}

		} else {
			// DEVUELTA
			facturaDB.EstadoSRI = respRecepcion.Estado
			msg := ""
			for _, comp := range respRecepcion.Comprobantes.Comprobante {
				for _, m := range comp.Mensajes.Mensaje {
					msg += fmt.Sprintf("%s: %s (%s); ", m.Identificador, m.Mensaje, m.InformacionAdicional)
				}
			}
			facturaDB.MensajeError = msg
		}
	}

	// 9. Generar RIDE (PDF) si no hubo error fatal técnico (Offline sí genera PDF)
	if facturaDB.EstadoSRI != "ERROR_TECNICO" {
		pdfBytes, errPdf := pdf.GenerarRIDE(*facturaXML, config.LogoPath, config.PDFTheme)
		if errPdf != nil {
			logger.Error("Error generando RIDE: %v", errPdf)
		} else {
			facturaDB.PDFRIDE = pdfBytes
		}
	}

	// 10. Auto-Guardar Cliente (Upsert)
	// Si el cliente no existe, lo creamos. Si existe, actualizamos nombre/dirección.
	var cliente db.Client
	var tipoID string
	
	// Inferencia simple de TipoID
	if dto.ClienteID == "9999999999999" {
		tipoID = "07" // Consumidor Final
	} else if len(dto.ClienteID) == 13 {
		tipoID = "04" // RUC
	} else if len(dto.ClienteID) == 10 {
		tipoID = "05" // Cédula
	} else {
		tipoID = "06" // Pasaporte / Otro
	}

	// Buscamos si existe
	if err := db.GetDB().First(&cliente, "id = ?", dto.ClienteID).Error; err != nil {
		// No existe, crear nuevo
		cliente = db.Client{
			ID:        dto.ClienteID,
			TipoID:    tipoID,
			Nombre:    dto.ClienteNombre,
			Direccion: dto.ClienteDireccion,
			Email:     dto.ClienteEmail,
			Telefono:  dto.ClienteTelefono,
		}
		db.GetDB().Create(&cliente)
	} else {
		// Ya existe, actualizamos datos básicos para mantenerlos al día
		cliente.Nombre = dto.ClienteNombre
		cliente.Direccion = dto.ClienteDireccion
		if dto.ClienteEmail != "" {
			cliente.Email = dto.ClienteEmail
		}
		if dto.ClienteTelefono != "" {
			cliente.Telefono = dto.ClienteTelefono
		}
		db.GetDB().Save(&cliente)
	}

	// 11. Guardar Factura en Base de Datos
	if err := db.GetDB().Create(facturaDB).Error; err != nil {
		return fmt.Errorf("error guardando factura en DB: %v", err)
	}

	// 12. Guardar Items de Factura para Reportería
	for _, item := range dto.Items {
		facturaItem := db.FacturaItem{
			FacturaClave:   claveAcceso,
			ProductoSKU:    item.Codigo,
			Nombre:         item.Nombre,
			Cantidad:       item.Cantidad,
			PrecioUnitario: item.Precio,
			Subtotal:       item.Cantidad * item.Precio,
			PorcentajeIVA:  item.PorcentajeIVA,
		}
		db.GetDB().Create(&facturaItem)
	}

	// Actualizar DTO de retorno con la clave generada

	dto.ClaveAcceso = claveAcceso

	return nil

}
