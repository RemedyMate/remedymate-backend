package interfaces

// IMailService defines the interface for sending emails
type IMailService interface {
	Send(to, subject, htmlBody string) error
}
