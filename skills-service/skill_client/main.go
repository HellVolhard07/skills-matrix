package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	skillpb "google.golang.org/grpc/SKILLS_MTRX/protobuf" // Adjust this import path
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := skillpb.NewSkillServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Add Skills
	log.Println("--- Adding Skills ---")
	skill1 := &skillpb.Skill{Id: "1", Name: "Go Programming", Description: "Expertise in Go language"}
	res1, err := c.AddSkillToUser(ctx, &skillpb.AddSkillRequest{UserId: "jdaniel", Skill: skill1})
	if err != nil {
		log.Printf("could not add skill: %v", err)
	} else {
		log.Printf("Added skill: %s", res1.GetName())
	}

	skill2 := &skillpb.Skill{Id: "2", Name: "C++", Description: "Development expertise using C++"}
	res2, err := c.AddSkillToUser(ctx, &skillpb.AddSkillRequest{UserId: "dsingh", Skill: skill2})
	if err != nil {
		log.Printf("could not add skill: %v", err)
	} else {
		log.Printf("Added skill: %s", res2.GetName())
	}

	skill3 := &skillpb.Skill{Id: "3", Name: "Microservices Architecture", Description: "Writing microservices in Go"}
	res3, err := c.AddSkillToUser(ctx, &skillpb.AddSkillRequest{UserId: "skhandel", Skill: skill3})
	if err != nil {
		log.Printf("could not add skill: %v", err)
	} else {
		log.Printf("Added skill: %s", res3.GetName())
	}

	// Get Skills
	log.Println("\n--- Getting Skills for existing user ---")
	userSkills, err := c.GetSkillsByUser(ctx, &skillpb.GetSkillsByUserRequest{UserId: "jdaniel"})
	if err != nil {
		log.Printf("could not get skills by user: %v", err)
	} else {
		log.Printf("Skills for requested user")
		for _, skill := range userSkills.GetSkills() {
			log.Printf("- ID: %s, Name: %s", skill.GetId(), skill.GetName())
		}
	}

	// Get All
	log.Println("\n--- Printing All Skills ---")
	allSkills, err := c.GetAllSkills(ctx, &skillpb.Empty{})
	if err != nil {
		log.Printf("could not get all skills: %v", err)
	} else {
		log.Printf("All Skills:")
		for _, skill := range allSkills.GetSkills() {
			log.Printf("- ID: %s, Name: %s", skill.GetId(), skill.GetName())
		}
	}

	// Remove Skill
	log.Println("\n--- Removing a skill from user---")
	_, err = c.RemoveSkillFromUser(ctx, &skillpb.RemoveSkillRequest{UserId: "dsingh", SkillId: "2"})
	if err != nil {
		log.Printf("could not remove skill: %v", err)
	} else {
		log.Println("Skill '2' removed successfully from user")
	}

	// Verify Removal of skill
	log.Println("\n--- Verifying Skills for User after removal ---")
	userSkillsAfterRemoval, err := c.GetSkillsByUser(ctx, &skillpb.GetSkillsByUserRequest{UserId: "dsingh"})
	if err != nil {
		log.Printf("could not get skills by user: %v", err)
	} else {
		log.Printf("Skills for user after removal:")
		if len(userSkillsAfterRemoval.GetSkills()) == 0 {
			log.Println("No skills found.")
		} else {
			for _, skill := range userSkillsAfterRemoval.GetSkills() {
				log.Printf("- ID: %s, Name: %s", skill.GetId(), skill.GetName())
			}
		}
	}
}
