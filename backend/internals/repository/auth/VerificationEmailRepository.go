package authRepository

import (
	"os"
	"log"

	"github.com/resend/resend-go/v2"
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
func (v *VerificationEmailRepository) SendVerificationEmail(email string) (error) {
	params := &resend.SendEmailRequest{
        From:    "onboarding@resend.dev",
        To:      []string{"trentjkelly1@gmail.com"},
        Subject: "Hello World",
        Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
    }

    sent, err := v.resend.Emails.Send(params)
	if err != nil {
		return err
	}

	log.Println(sent)

	return nil
}
