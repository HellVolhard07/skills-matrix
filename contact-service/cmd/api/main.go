package main

import (
	"contact/contact"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

const gRPCPort = "50001"

type ContactConfig struct{}

func main() {

	app := ContactConfig{}

	log.Printf("Starting mailer service on port %s\n", gRPCPort)

	app.gRPCListen()
}

func (app *ContactConfig) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}

	s := grpc.NewServer()
	contact.RegisterContactServiceServer(s, &ContactServer{
		Mailer: app.createMailer(),
	})

	log.Printf("gRPC server started on port: %s", gRPCPort)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}
}

func (app *ContactConfig) createMailer() Mailer {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mailer{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}
	return m
}
