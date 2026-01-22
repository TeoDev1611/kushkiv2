package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
)

type MailService struct{}

func NewMailService() *MailService {
	return &MailService{}
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

func (s *MailService) SendInvoiceEmail(config SMTPConfig, to string, razonSocial string, pdfContent []byte, secuencial string) error {
	subject := fmt.Sprintf("Comprobante Electrónico - %s - Factura %s", razonSocial, secuencial)
	filename := fmt.Sprintf("FACTURA-%s.pdf", secuencial)
	
	// Plantilla HTML Profesional
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; color: #333; line-height: 1.6; margin: 0; padding: 0; }
		.container { max-width: 600px; margin: 20px auto; border: 1px solid #e2e8f0; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 6px rgba(0,0,0,0.05); }
		.header { background-color: #34d399; padding: 30px; text-align: center; color: white; }
		.content { padding: 40px; background-color: white; }
		.footer { background-color: #f8fafc; padding: 20px; text-align: center; font-size: 12px; color: #64748b; border-top: 1px solid #e2e8f0; }
		.summary { background-color: #f1f5f9; padding: 20px; border-radius: 8px; margin: 25px 0; }
		.summary-row { display: block; margin-bottom: 8px; border-bottom: 1px solid #e2e8f0; padding-bottom: 8px; }
		.summary-row:last-child { border-bottom: none; }
		.label { font-weight: bold; color: #475569; width: 100px; display: inline-block; }
		.value { color: #0f172a; font-weight: 600; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1 style="margin:0; font-size: 28px; font-weight: 300; letter-spacing: 1px;">Nuevo Comprobante</h1>
		</div>
		<div class="content">
			<p style="font-size: 16px;">Estimado cliente,</p>
			<p>Le informamos que <strong>%s</strong> ha generado un nuevo comprobante electrónico a su nombre.</p>
			
			<div class="summary">
				<div style="font-weight:bold; margin-bottom:15px; color:#059669; text-transform: uppercase; font-size: 12px; letter-spacing: 1px;">Detalles del Documento:</div>
				<div class="summary-row">
					<span class="label">Tipo:</span> <span class="value">Factura</span>
				</div>
				<div class="summary-row">
					<span class="label">Número:</span> <span class="value">%s</span>
				</div>
				<div class="summary-row">
					<span class="label">Estado:</span> <span class="value" style="color: #059669;">Autorizado SRI</span>
				</div>
			</div>

			<p>Adjunto a este correo encontrará el archivo PDF (RIDE) con el detalle completo de su transacción.</p>
			
			<p style="text-align: center; margin-top: 30px; color: #94a3b8;">
				<small>Gracias por su confianza.</small>
			</p>
		</div>
		<div class="footer">
			Este es un correo automático generado por el sistema de facturación.<br>
			<strong>Kushki App - Tecnología hecha en Ecuador</strong>
		</div>
	</div>
</body>
</html>
`, razonSocial, secuencial)

	return s.sendMailWithAttachment(config, to, subject, body, pdfContent, filename)
}

func (s *MailService) SendTestEmail(config SMTPConfig, to string) error {
	subject := "Kushki App - Prueba de Conexión Exitosa"
	body := `
<!DOCTYPE html>
<html>
<body style="font-family: sans-serif; text-align: center; padding: 20px;">
	<h2 style="color: #34d399;">¡Conexión Exitosa!</h2>
	<p>Si estás leyendo esto, tu configuración de correo en Kushki App funciona correctamente.</p>
	<p style="color: #64748b;">Ya puedes enviar facturas a tus clientes.</p>
</body>
</html>`
	
	// Para test enviamos sin adjunto o con adjunto vacío si la función base lo requiere?
	// sendMailWithAttachment expects an attachment.
	// Let's create a sendMailWithoutAttachment or simply pass empty attachment.
	// Passing empty attachment works with current impl, creates an empty file.
	// Better to make attachment optional in sendMailWithAttachment or just send a dummy.
	
	return s.sendMailWithAttachment(config, to, subject, body, []byte("test"), "test.txt")
}

func (s *MailService) sendMailWithAttachment(config SMTPConfig, to, subject, body string, attachment []byte, filename string) error {
	if config.Host == "" || config.Port == 0 {
		return fmt.Errorf("configuración SMTP incompleta")
	}

	auth := smtp.PlainAuth("", config.User, config.Password, config.Host)
	
	// Buffer para construir el mensaje
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Headers básicos
	headers := make(map[string]string)
	headers["From"] = config.User
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary())

	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")

	// 1. Cuerpo del mensaje (HTML)
	partHeader := make(map[string][]string)
	partHeader["Content-Type"] = []string{"text/html; charset=UTF-8"}
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return err
	}
	part.Write([]byte(body))

	// 2. Adjunto (PDF)
	partHeader = make(map[string][]string)
	partHeader["Content-Type"] = []string{"application/pdf"}
	partHeader["Content-Disposition"] = []string{fmt.Sprintf("attachment; filename=\"%s\"", filename)}
	partHeader["Content-Transfer-Encoding"] = []string{"base64"}
	
	part, err = writer.CreatePart(partHeader)
	if err != nil {
		return err
	}
	
	encoder := base64.NewEncoder(base64.StdEncoding, part)
	encoder.Write(attachment)
	encoder.Close()

	writer.Close()

	// Enviar
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err = smtp.SendMail(addr, auth, config.User, []string{to}, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
