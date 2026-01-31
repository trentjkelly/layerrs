package authRepository

import (
	"os"
	"log"

	"github.com/resend/resend-go/v2"
	"github.com/trentjkelly/layerrs/internals/entities"
)

type VerificationEmailRepository struct {
	resend *resend.Client
}

// Constructor for the VerificationEmailRepository
func NewVerificationEmailRepository() *VerificationEmailRepository {
	resendClient := resend.NewClient(os.Getenv("EMAIL_API_KEY"))

	return &VerificationEmailRepository{
		resend: resendClient,
	}
}

// Creates a verification email and sends it to the given email address
func (v *VerificationEmailRepository) SendEmail(emailInfo entities.LoginEmailInfo) error {
	params := &resend.SendEmailRequest{
        From:    emailInfo.EmailSender,
        To:      emailInfo.EmailRecipients,
        Subject: emailInfo.EmailSubject,
        Html:    emailInfo.EmailBodyHTML,
    }

    sent, err := v.resend.Emails.Send(params)
	if err != nil {
		return err
	}

	log.Println(sent)

	return nil
}
