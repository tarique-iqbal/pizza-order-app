package container

import (
	"os"

	aEmail "email-service/internal/application/email"
	dEmail "email-service/internal/domain/email"
	iEmail "email-service/internal/infrastructure/email"
	"email-service/internal/infrastructure/messaging"
)

const emailTemplatePath = "internal/infrastructure/email/templates"

type Container struct {
	Dispatcher dEmail.EventDispatcher
	Consumer   *messaging.RabbitMQConsumer
}

func NewContainer() (*Container, error) {
	amqpURL := os.Getenv("RABBITMQ_URL")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	senderEmail := os.Getenv("SENDER_EMAIL")

	smtpSender := iEmail.NewSMTPSender(smtpHost, smtpPort, smtpUser, smtpPass, senderEmail)
	template := iEmail.NewHTMLTemplateLoader(emailTemplatePath)

	userRegisteredHandler := aEmail.NewUserRegisteredHandler(smtpSender, template)
	emailVerificationCreatedHandler := aEmail.NewEmailVerificationCreatedHandler(smtpSender, template)

	dispatcher := aEmail.NewEventDispatcher()
	dispatcher.Register("user.registered", userRegisteredHandler)
	dispatcher.Register("email.verification_created", emailVerificationCreatedHandler)

	consumer, err := messaging.NewRabbitMQConsumer(amqpURL)

	return &Container{
		Dispatcher: dispatcher,
		Consumer:   consumer,
	}, err
}
