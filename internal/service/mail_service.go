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
	subject := fmt.Sprintf("Factura enviada por %s", razonSocial)
	filename := fmt.Sprintf("FACTURA-%s.pdf", secuencial)
	
	// Branding de Kushki en el cuerpo
	body := fmt.Sprintf(`
<html>
<body>
	<p>Estimado cliente,</p>
	<p>Adjunto encontrará su factura electrónica emitida por <strong>%s</strong>.</p>
	<br>
	<p>Este documento ha sido generado y enviado automáticamente por <strong>Kushki Facturador</strong>.</p>
	<p><small>Kushki - Simplificando tus cobros y facturación.</small></p>
</body>
</html>
`, razonSocial)

	return s.sendMailWithAttachment(config, to, subject, body, pdfContent, filename)
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
