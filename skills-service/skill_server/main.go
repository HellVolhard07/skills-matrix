package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	skillpb "google.golang.org/grpc/SKILLS_MTRX/protobuf" // Adjust this import path
)

// server implements the SkillServiceServer interface
type server struct {
	userSkills                              map[string][]*skillpb.Skill
	allSkills                               map[string]*skillpb.Skill
	skillpb.UnimplementedSkillServiceServer // Embed this for forward compatibility
}

func NewSkillServiceServer() *server {
	return &server{
		userSkills: make(map[string][]*skillpb.Skill),
		allSkills:  make(map[string]*skillpb.Skill),
	}
}

func (s *server) AddSkillToUser(ctx context.Context, req *skillpb.AddSkillRequest) (*skillpb.Skill, error) {
	log.Printf("AddSkillToUser called for user: %s, skill: %s", req.GetUserId(), req.GetSkill().GetName())

	// Basic validation
	if req.GetUserId() == "" || req.GetSkill() == nil || req.GetSkill().GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id, skill, and skill.id cannot be empty")
	}

	s.userSkills[req.GetUserId()] = append(s.userSkills[req.GetUserId()], req.GetSkill())
	s.allSkills[req.GetSkill().GetId()] = req.GetSkill() // Keep track of all skills

	return req.GetSkill(), nil
}

func (s *server) RemoveSkillFromUser(ctx context.Context, req *skillpb.RemoveSkillRequest) (*skillpb.Empty, error) {
	log.Printf("RemoveSkillFromUser called for user: %s, skill ID: %s", req.GetUserId(), req.GetSkillId())

	if req.GetUserId() == "" || req.GetSkillId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id and skill_id cannot be empty")
	}

	// Remove skill from user's list
	userSkills := s.userSkills[req.GetUserId()]
	found := false
	for i, skill := range userSkills {
		if skill.GetId() == req.GetSkillId() {
			s.userSkills[req.GetUserId()] = append(userSkills[:i], userSkills[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return nil, status.Errorf(codes.NotFound, "skill with ID %s not found for user %s", req.GetSkillId(), req.GetUserId())
	}

	return &skillpb.Empty{}, nil
}

func (s *server) GetSkillsByUser(ctx context.Context, req *skillpb.GetSkillsByUserRequest) (*skillpb.SkillList, error) {
	log.Printf("GetSkillsByUser called for user: %s", req.GetUserId())

	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id cannot be empty")
	}

	skills, ok := s.userSkills[req.GetUserId()]
	if !ok {
		return &skillpb.SkillList{Skills: []*skillpb.Skill{}}, nil // Return empty list if no skills
	}

	return &skillpb.SkillList{Skills: skills}, nil
}

func (s *server) GetAllSkills(ctx context.Context, req *skillpb.Empty) (*skillpb.SkillList, error) {
	log.Println("GetAllSkills called")

	var allSkills []*skillpb.Skill
	for _, skill := range s.allSkills {
		allSkills = append(allSkills, skill)
	}

	return &skillpb.SkillList{Skills: allSkills}, nil
}

// Main function to run the gRPC server
func main() {
	lis, err := net.Listen("tcp", ":50051") // Listen on port 50051
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	skillpb.RegisterSkillServiceServer(s, NewSkillServiceServer())

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
