package mail

import "gopkg.in/gomail.v2"

type Mail struct {
	Sender    string
	Receaver  string
	GmailPass string
	Message   string
	Subject   string
}
type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) serverName() string {
	return s.host + ":" + s.port
}

func NewMail(Sender string, GmailPass string, Message string, Subject string) Mail {
	mail := Mail{
		Sender:    Sender,
		GmailPass: GmailPass,
		Message:   Message,
		Subject:   Subject,
	}
	return mail
}

func (mail *Mail) SendEmail(Receaver string) {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.Sender)
	m.SetHeader("To", Receaver)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Message)
	d := gomail.NewDialer("smtp.gmail.com", 587, mail.Sender, mail.GmailPass)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
