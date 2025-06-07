package main

import (
	"fmt"
	"log"
	"net"
	"skills/skills"

	"google.golang.org/grpc"
)

const gRPCPort = "50053"

type SkillConfig struct{}

func main() {

	app := SkillConfig{}

	log.Printf("Starting skills service on port %s\n", gRPCPort)

	app.gRPCListen()

}

func (app *SkillConfig) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}

	s := grpc.NewServer()
	skills.RegisterSkillServiceServer(s, &SkillsServer{
		userSkills: make(map[string][]*skills.Skill),
		allSkills:  make(map[string]*skills.Skill),
	})

	log.Printf("gRPC server started on port: %s", gRPCPort)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to listen to gRPC: %v", err)
	}
}
