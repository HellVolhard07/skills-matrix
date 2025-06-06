package main

import (
	"broker/contact"
	"broker/search"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action  string         `json:"action"`
	Auth    AuthPayload    `json:"auth,omitempty"`
	Contact ContactPayload `json:"contact,omitempty"`
	Search  SearchPayload  `json:"search,omiempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ContactPayload struct {
	SenderId   string `json:"senderid"`
	ReceiverId string `json:"receiverid"`
	Subject    string `json:"subject"`
	Message    string `json:"message"`
}

type SearchPayload struct {
	SkillName   string `json:"skillname"`
	Category    string `json:"category"`
	Proficiency string `json:"proficiency"`
}

func (app *Config) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "contact":
		app.contactHandler(w, requestPayload.Contact)
	case "search":
		app.searchHandler(w, requestPayload.Search)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we will send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable to read response.Body into
	var jsonFromService jsonResponse

	// decode json from auth-service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) contactHandler(w http.ResponseWriter, cp ContactPayload) {

	conn, err := grpc.NewClient("contact-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := contact.NewContactServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = c.SendContactRequest(ctx, &contact.ContactRequest{
		SenderId:   cp.SenderId,
		ReceiverId: cp.ReceiverId,
		Subject:    cp.Subject,
		Message:    cp.Message,
	})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Email sent",
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) searchHandler(w http.ResponseWriter, s SearchPayload) {
	conn, err := grpc.NewClient("search-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := search.NewSearchServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.SearchUsersBySkill(ctx, &search.SearchRequest{
		SkillName:   s.SkillName,
		Category:    s.Category,
		Proficiency: s.Proficiency,
	})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Optionally convert the gRPC Users to a format that's JSON-friendly
	users := make([]map[string]any, len(res.Users))
	for i, u := range res.Users {
		skills := make([]string, len(u.Skills))
		copy(skills, u.Skills)
		users[i] = map[string]any{
			"id":         u.Id,
			"name":       u.Name,
			"department": u.Department,
			"skills":     skills,
		}
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Search results retrieved",
		Data:    users,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
