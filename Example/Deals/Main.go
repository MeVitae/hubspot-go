package main

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"os"

	"gitub.com/MeVitae/hubspot-go"
)

var (
	name        = flag.String("name", "", "Name of repo to create in authenticated user's GitHub account.")
	description = flag.String("description", "", "Description of created repo.")
	private     = flag.Bool("private", false, "Will created repo be private.")
	autoInit    = flag.Bool("auto-init", false, "Pass true to create an initial commit with empty README.")
)

func main() {
	flag.Parse()
	
	token := os.Getenv("HUBSPOT_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	ctx := context.Background()
	client := hubspot.NewClient(nil).WithAuthToken(token)

	r := &github.Repository{Name: name, Private: private, Description: description, AutoInit: autoInit}
	repo, _, err := client.Repositories.Create(ctx, "", r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully created new repo: %v\n", repo.GetName())
}