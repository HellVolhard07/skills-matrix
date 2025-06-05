package main

import "log"

const gRPCPort = "50001"

type ContactConfig struct{}

func main() {

	app := ContactConfig{}

	log.Printf("Starting mailer service on port %s\n", gRPCPort)

	app.gRPCListen()
}
