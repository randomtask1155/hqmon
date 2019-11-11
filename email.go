package main

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	smail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

/*
   send the email using the users gmail account
*/
func SendGmail(FromName, FromEmail, ToName, ToEmail, Subject, Message, PlainMessage string) error {
	message := smail.NewSingleEmail(&smail.Email{Name: FromName, Address: FromEmail}, Subject, &smail.Email{Name: ToName, Address: ToEmail}, PlainMessage, Message)

	client := sendgrid.NewSendClient(sendGridAPIKey)
	response, err := client.Send(message)
	if err != nil {
		logger.Println(err)
		return err
	}

	// error code reference https://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html
	if response.StatusCode > 400 {
		return fmt.Errorf("Failed to send email:%d: %s\n", response.StatusCode, response.Body)
	}
	return nil
}