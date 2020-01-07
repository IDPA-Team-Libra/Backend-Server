package mail

import (
	"fmt"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

//Mail struct to hold auth and message information for a mail
type Mail struct {
	Sender   string
	Receaver string
	Pass     string
	Message  string
	Subject  string
	UserID   string
}

//Configuration contains auth elements for the mail
type Configuration struct {
	Sender string
	Pass   string
	UserID string
}

var mailConfiguration Configuration

//SetMailConfiguration sets the global mail config
func SetMailConfiguration(configuration Configuration) {
	mailConfiguration = configuration
}

//LoadMailConfiguration gets the global mail config
func LoadMailConfiguration() Configuration {
	return mailConfiguration
}

//NewMail create a new mail
func NewMail(Receaver string) Mail {
	mail := Mail{
		Receaver: Receaver,
	}
	return mail
}

//ApplyConfiguration applys the configuration to the mail instance
func (mail *Mail) ApplyConfiguration(configuration Configuration) {
	mail.Sender = configuration.Sender
	mail.Pass = configuration.Pass
	mail.UserID = configuration.UserID
}

//SendEmail sends the mail over the mailjet api [gmail blacklisted the server, so move to a non smtp api]
func (mail *Mail) SendEmail() {
	mailjetClient := mailjet.NewMailjetClient(mail.UserID, mail.Pass)
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: mail.Sender,
				Name:  "project-libra-company",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: mail.Receaver,
					Name:  mail.Receaver,
				},
			},
			Subject:  mail.Subject,
			TextPart: "",
			HTMLPart: mail.Message,
			CustomID: "project-libra-company",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		fmt.Println(err)
	}
}

//SendWelcomeEmail sends the welcome email for a new user
func (mail *Mail) SendWelcomeEmail() {
	message :=
		`
		Wilkommen neuer Nutzer<br /><br />

		Es freut uns, dass Sie sich für libra entschieden haben. 
		Wir hoffen, dass Ihnen der Service gefällt und Ihren Ansprüchen entspricht. 
		<br /><br />
		Wir werden Sie über Ihre vertagten Transaktionen auf dem laufenden halten 
		und Sie darüber hinaus so wenig wie möglich stören.
		<br /><br />
		Ebenfalls halten wir es uns vor, Ihnen von Zeit zu Zeit eine Mail mit anfallenden Informationen zu senden.<br /><br />
		Ihr Libra-Team
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

//SendDelayedTransactionEmail sends an update to the user regarding his delayed transactions
func (mail *Mail) SendDelayedTransactionEmail(totalOperations int, results []string) {
	mail.Subject = "Informationen zu Ihren vertagten Transaktionen"
	message :=
		`
		Ihre vertagten Transaktionen sind soeben abgeschlossen worden. <br />
		Hier sind die Resultate Ihrer Käufe und Verkäufe: <br /><br />
		<br />
		<ul>
	`
	for index := range results {
		message += fmt.Sprintf("<li>%s</li><br />", results[index])
	}
	message += "</ul><br /> Wir hoffen natürlich, dass Sie damit zufrieden sind.<br />"
	message += `<b>Ihr Libra-Team</b>
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
