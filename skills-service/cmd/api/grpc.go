package main

import (
	"context"
	"log"
	"skills/skills"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server implements the SkillServiceServer interface
type SkillsServer struct {
	skills.UnimplementedSkillServiceServer
	userSkills map[string][]*skills.Skill
	allSkills  map[string]*skills.Skill
}

func (s *SkillsServer) AddSkillToUser(ctx context.Context, req *skills.AddSkillRequest) (*skills.Skill, error) {
	log.Printf("AddSkillToUser called for user: %s, skill: %s", req.GetUserId(), req.GetSkill().GetName())

	// Basic validation
	if req.GetUserId() == "" || req.GetSkill() == nil || req.GetSkill().GetId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id, skill, and skill.id cannot be empty")
	}

	s.userSkills[req.GetUserId()] = append(s.userSkills[req.GetUserId()], req.GetSkill())
	s.allSkills[req.GetSkill().GetId()] = req.GetSkill() // Keep track of all skills

	return req.GetSkill(), nil
}

func (s *SkillsServer) RemoveSkillFromUser(ctx context.Context, req *skills.RemoveSkillRequest) (*skills.Empty, error) {
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

	return &skills.Empty{}, nil
}

func (s *SkillsServer) GetSkillsByUser(ctx context.Context, req *skills.GetSkillsByUserRequest) (*skills.SkillList, error) {
	log.Printf("GetSkillsByUser called for user: %s", req.GetUserId())

	if req.GetUserId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id cannot be empty")
	}

	skillList, ok := s.userSkills[req.GetUserId()]
	if !ok {
		return &skills.SkillList{Skills: []*skills.Skill{}}, nil // Return empty list if no skills
	}

	return &skills.SkillList{Skills: skillList}, nil
}

func (s *SkillsServer) GetAllSkills(ctx context.Context, req *skills.Empty) (*skills.SkillList, error) {
	log.Println("GetAllSkills called")

	var allSkills []*skills.Skill
	for _, skill := range s.allSkills {
		allSkills = append(allSkills, skill)
	}

	return &skills.SkillList{Skills: allSkills}, nil
}
