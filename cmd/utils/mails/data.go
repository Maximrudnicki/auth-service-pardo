package mails

type EmailData struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Name    string   `json:"username"`
}