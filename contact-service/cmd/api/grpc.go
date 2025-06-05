package main

import (
	"contact/contact"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ContactServer struct {
	contact.UnimplementedContactServiceServer
}

func (c *ContactServer) SendContactRequest(ctx context.Context, req *contact.ContactRequest) (*contact.ContactResponse, error) {
	// TODO: Implement sending email here
	res := &contact.ContactResponse{
		Error:   "",
		Message: "Email sent to xyz.com",
	}

	return res, nil
}

func (app *ContactConfig) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}

	s := grpc.NewServer()
	contact.RegisterContactServiceServer(s, &ContactServer{})

	log.Printf("gRPC server started on port: %s", gRPCPort)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}
}
