package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fogfish/blueprint-serverless-golang/pkg/api"
	gurlhttp "github.com/fogfish/gurl/v2/http"
)

func main() {
	host := os.Args[1]
	fmt.Printf("==> integration test against %s\n", host)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test 1: Health check - endpoint responds
	fmt.Print("[1] GET /petshop/pets responds 200 ... ")
	client := api.NewPetShop(gurlhttp.New(gurlhttp.WithClient(&http.Client{Timeout: 10 * time.Second})), host)
	pets, err := client.List(ctx)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("PASS")

	// Test 2: Response contains seeded pets
	fmt.Printf("[2] Response contains pets (got %d) ... ", len(pets.Pets))
	if len(pets.Pets) == 0 {
		fmt.Println("FAIL: no pets returned")
		os.Exit(1)
	}
	fmt.Println("PASS")

	// Test 3: Pet structure is valid
	fmt.Print("[3] Pet has required fields ... ")
	p := pets.Pets[0]
	if p.ID == "" || p.Category == "" || p.Price == 0 || p.Url == "" {
		fmt.Printf("FAIL: incomplete pet %+v\n", p)
		os.Exit(1)
	}
	fmt.Println("PASS")

	// Test 4: Pagination cursor present
	fmt.Print("[4] Pagination cursor present ... ")
	if pets.Next == nil || len(*pets.Next) == 0 {
		fmt.Println("FAIL: no cursor")
		os.Exit(1)
	}
	fmt.Println("PASS")

	// Test 5: Lookup individual pet
	fmt.Print("[5] GET /petshop/pets/{id} returns pet ... ")
	pet, err := client.Pet(ctx, pets.Pets[0].Url)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		os.Exit(1)
	}
	if pet.ID != pets.Pets[0].ID {
		fmt.Printf("FAIL: expected %s got %s\n", pets.Pets[0].ID, pet.ID)
		os.Exit(1)
	}
	fmt.Println("PASS")

	fmt.Println("\n==> all integration tests passed")
}
