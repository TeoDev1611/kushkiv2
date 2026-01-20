package service

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"kushkiv2/internal/db"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
)

type MailService struct{}

func NewMailService() *MailService {
	return &MailService{}
}

// QueueEmail añade un correo a la base de datos para ser procesado por el worker.
func (s *MailService) QueueEmail(to, subject, body string, attachment []byte, attachName string) error {
	queueItem := db.EmailQueue{
		To:         to,
		Subject:    subject,
		Body:       body,
		Attachment: attachment,
		AttachName: attachName,
		Status:     "PENDIENTE",
	}
	return db.GetDB().Create(&queueItem).Error
}

// StartWorker inicia una Goroutine que procesa la cola de correos periódicamente.
func (s *MailService) StartWorker() {
	go func() {
		// Primera ejecución inmediata
		s.processQueue()
		for {
			time.Sleep(1 * time.Minute)
			s.processQueue()
		}
	}()
}

func (s *MailService) processQueue() {
	var pending []db.EmailQueue
	// Buscar correos pendientes o con error (máximo 3 reintentos)
	db.GetDB().Where("status = ? OR (status = ? AND retry_count < 3)", "PENDIENTE", "ERROR").Find(&pending)

	if len(pending) == 0 {
		return
	}

	// Obtener configuración SMTP
	var config db.EmisorConfig
	if err := db.GetDB().First(&config).Error; err != nil {
		return
	}

	if config.SMTPHost == "" || config.SMTPUser == "" {
		return 
	}

	for _, item := range pending {
		err := s.sendActualEmail(&config, item)
		if err != nil {
			item.Status = "ERROR"
			item.RetryCount++
			item.LastError = err.Error()
			fmt.Printf("MailWorker: Error enviando a %s: %v\n", item.To, err)
		} else {
			item.Status = "ENVIADO"
			item.LastError = ""
		}
		db.GetDB().Save(&item)
	}
}

func (s *MailService) sendActualEmail(config *db.EmisorConfig, item db.EmailQueue) error {
	from := mail.Address{Name: config.RazonSocial, Address: config.SMTPUser}
	to := mail.Address{Name: "", Address: item.To}

	parts := strings.Split(config.SMTPHost, ":")
	host := parts[0]

	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPass, host)

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	buf.WriteString(fmt.Sprintf("From: %s\r\n", from.String()))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to.String()))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", item.Subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
	buf.WriteString("\r\n")

	bodyPart, _ := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"text/html; charset=UTF-8"}})
	bodyPart.Write([]byte(item.Body))

	if len(item.Attachment) > 0 {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Type", "application/pdf")
		h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, item.AttachName))
		h.Set("Content-Transfer-Encoding", "base64")
		
		attachPart, _ := writer.CreatePart(h)
		attachPart.Write([]byte(base64.StdEncoding.EncodeToString(item.Attachment)))
	}

	writer.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, 
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", config.SMTPHost, tlsConfig)
	if err != nil {
		return fmt.Errorf("error TLS Dial: %v", err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("error SMTP Client: %v", err)
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("error Auth: %v", err)
	}

	if err = client.Mail(from.Address); err != nil {
		return err
	}
	if err = client.Rcpt(to.Address); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}