package main

import (
	"contact/contact"
	"context"
	"fmt"
	"os"
)

type ContactServer struct {
	contact.UnimplementedContactServiceServer
	Mailer Mailer
}

func (c *ContactServer) SendContactRequest(ctx context.Context, req *contact.ContactRequest) (*contact.ContactResponse, error) {
	msg := Message{
		From:    os.Getenv("FROM_ADDRESS"),
		To:      req.ReceiverId,
		Subject: req.Subject,
		Data:    req.Message,
	}

	err := c.Mailer.SendSMTPMessage(msg)
	if err != nil {
		return &contact.ContactResponse{
			Error:   err.Error(),
			Message: "",
		}, nil
	}

	res := &contact.ContactResponse{
		Error:   "",
		Message: fmt.Sprintf("Email sent to %s", req.ReceiverId),
	}
	return res, nil
}
