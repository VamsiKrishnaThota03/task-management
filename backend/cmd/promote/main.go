// Command promote changes a user's role, e.g. to create the first admin:
//
//	go run ./cmd/promote -username alice            # make alice an admin
//	go run ./cmd/promote -username alice -role user # demote again
package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"task-management-system/internal/database"
	"task-management-system/internal/models"
)

func main() {
	username := flag.String("username", "", "username to update (required)")
	role := flag.String("role", models.RoleAdmin, "new role: admin or user")
	flag.Parse()

	if *username == "" || (*role != models.RoleAdmin && *role != models.RoleUser) {
		flag.Usage()
		os.Exit(2)
	}

	_ = godotenv.Load()
	dsn := os.Getenv("DB_URI")
	if dsn == "" {
		log.Fatal("DB_URI is not set (see .env.example)")
	}

	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}

	res := db.Model(&models.User{}).Where("username = ?", *username).Update("role", *role)
	if res.Error != nil {
		log.Fatal(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Fatalf("no user named %q", *username)
	}
	log.Printf("%s is now %s", *username, *role)
}
