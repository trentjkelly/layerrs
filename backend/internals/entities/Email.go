package entities

type LoginEmailInfo struct {
	EmailSender     string
	EmailRecipients []string
	EmailSubject    string
	EmailBodyHTML   string
}
