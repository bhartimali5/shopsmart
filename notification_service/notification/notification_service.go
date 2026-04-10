package notification

// User represents the recipient of a notification
type User struct {
	Name  string
	Email string
}

// EmailClient is the interface all doubles will implement
type EmailClient interface {
	Send(to, subject, body string) error
}

// Logger is a simple logging interface
type Logger interface {
	Log(message string)
}

// NotificationService sends notifications using an email client and logger
type NotificationService struct {
	emailClient EmailClient
	logger      Logger
}

func NewNotificationService(emailClient EmailClient, logger Logger) *NotificationService {
	return &NotificationService{emailClient: emailClient, logger: logger}
}

func (n *NotificationService) SendWelcome(user User) error {
	err := n.emailClient.Send(
		user.Email,
		"Welcome to ShopSmart! 🎉", // changed subject
		"Hi "+user.Name+", welcome aboard!",
	)
	if err != nil {
		n.logger.Log("failed to send welcome email to " + user.Email)
		return err
	}
	n.logger.Log("welcome email sent to " + user.Email)
	return nil
}
