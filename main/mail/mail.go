package mail

import (
	"fmt"

	gmail "github.com/go-mail/mail"
)

type Mail struct {
	Sender    string
	Receaver  string
	GmailPass string
	Message   string
	Subject   string
}

type SMPTServer struct {
	host string
	port string
}

type MailConfiguration struct {
	Sender    string
	GmailPass string
}

var mailConfiguration MailConfiguration

func SetMailConfiguration(configuration MailConfiguration) {
	mailConfiguration = configuration
}

func LoadMailConfiguration() MailConfiguration {
	return mailConfiguration
}

func (s *SMPTServer) serverName() string {
	return s.host + ":" + s.port
}

func NewMail(Receaver string) Mail {
	mail := Mail{
		Receaver: Receaver,
	}
	return mail
}

func (mail *Mail) ApplyConfiguration(configuration MailConfiguration) {
	mail.Sender = configuration.Sender
	mail.GmailPass = configuration.GmailPass
}

func (mail *Mail) SendEmail() {
	m := gmail.NewMessage()
	m.SetHeader("From", mail.Sender)
	m.SetHeader("To", mail.Receaver)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Message)
	d := gmail.NewDialer("smtp.gmail.com", 587, mail.Sender, mail.GmailPass)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

//SendWelcomeEmail sends the welcome email for a new user
func (mail *Mail) SendWelcomeEmail() {
	message :=
		`
		Wilkommen neuer Nutzer<br /><br />

		Es freut uns das Du dich für libra entschieden hast. 
		Wir hoffen das Dir der Service gefällt und Deinen Ansprüchen entspricht. 
		<br /><br />
		Wir werden Dich über Deine vertagten Transaktionen auf dem laufenden halten 
		und Dich darüber hinaus so wenig wie möglich stören.
		<br /><br />
		Ebenfalls halten wir es uns vor, Dir von Zeit zu Zeit eine Mail mit anfallenden Informationen zu senden.<br /><br />
		Dein Libra-Team
		<br /><br />
		<i>Libra - The way to go</i>
		<br />
		<pre>
		| |    (_) |            
		| |     _| |__  _ __ __ _ 
		| |    | | '_ \| '__/ _  |
		| |____| | |_) | | | (_| |
		|______|_|_.__/|_|  \__,_|
		</pre>
	`
	mail.Message = message
	mail.Subject = "Hallo neuer Nutzer - Wilkommen bei Libra"
	mail.SendEmail()
}

func (mail *Mail) SendDelayedTransactionEmail(totalOperations int, results []string) {
	mail.Subject = "Informationen zu deinen vertagten Transaktionen"
	message :=
		`
		Hey... Du hast da einige vertagte Transaktionen gehabt, die heute durchgelaufen sind. <br />
		Hier sind mal die Resultate Deiner Käufe und Verkäufe: <br /><br />
		<br />
		<ul>
	`
	for index := range results {
		message += fmt.Sprintf("<li>%s</li><br />", results[index])
	}
	message += "</ul><br /> Wir hoffen natürlich, dass du damit zufrieden bist.<br />"
	message += `<b>Dein Libra-Team</b>
		<br /><br />
		<i>Libra - The way to go</i>
		<br />
		<pre>
		| |    (_) |            
		| |     _| |__  _ __ __ _ 
		| |    | | '_ \| '__/ _  |
		| |____| | |_) | | | (_| |
		|______|_|_.__/|_|  \__,_|
		</pre>`
	mail.Message = message
	mail.SendEmail()
}
