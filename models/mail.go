package models

type Mail struct {
	From        string
	To          string
	Subject     string
	ContentType string
	Content     string
}

type MailConfig struct {
	Host          string
	Port          int
	Username      string
	Password      string
	AllowInsecure bool
}
