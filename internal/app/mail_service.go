package app

type MailService interface {
	SendVerificationMail(email string) error
	SendPasswordResetMail(email string) error
}

type mailService struct{}

func NewMailService() MailService {
	return &mailService{}
}
