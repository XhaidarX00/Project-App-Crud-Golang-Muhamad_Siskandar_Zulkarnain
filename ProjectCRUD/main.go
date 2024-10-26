package main

import (
	service "ProjectCRUD/Service"
	utils "ProjectCRUD/Utils"
	"context"
	"fmt"
	"time"
)

func main() {
	utils.ClearScreen()
	educationService, err := service.NewEducationService("education_db.json")
	if err != nil {
		fmt.Printf("Error loading database: %v\n", err)
		return
	}

	if educationService.AuthenticateUser() {
		ctx := context.Background()
		deadline := time.Now().Add(3000 * time.Second)
		ctxwithdeadline, cancel := context.WithDeadline(ctx, deadline)
		defer cancel()
		educationService.Menu(ctxwithdeadline)
	}
}
