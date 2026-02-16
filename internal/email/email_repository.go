package email

type EmailRepository interface {
	SendEmail(email string) (string, error)
}
