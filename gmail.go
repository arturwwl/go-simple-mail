package mail

import (
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	// AuthGmail implements gmail oauth2
	AuthGmail AuthType = iota
)

type GmailConfig struct {
	ClientID     string
	ClientSecret string
}

// GmailConnect ...
func (server *SMTPServer) GmailConnect() (*gmail.Service, error) {
	config := oauth2.Config{
		ClientID:     server.GmailConfig.ClientID,
		ClientSecret: server.GmailConfig.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	return srv, nil
}

// GmailSend sends the composed email
func (email *Email) GmailSend(client *gmail.Service) error {
	var message gmail.Message
	message.Raw = email.GetMessage()

	// Send the message
	_, err := client.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return err
	}
	return nil
}
