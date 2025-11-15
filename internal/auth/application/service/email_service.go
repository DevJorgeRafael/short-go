package service

type EmailService interface {
	SendPasswordResetCode(toEmail string, code string) error
}