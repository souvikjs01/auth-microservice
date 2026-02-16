package repository

type EmailRepository struct{}

func NewEmailRepository() *EmailRepository {
	return &EmailRepository{}
}

func (r *EmailRepository) SendEmail(email string) (string, error) {
	// Implement the logic to send email here
	return "Email sent successfully", nil
}
