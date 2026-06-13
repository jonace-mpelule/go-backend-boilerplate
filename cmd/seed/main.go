package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/username/project-name/ent/user"
	"github.com/username/project-name/internal/config"
	"github.com/username/project-name/internal/permissions"
	"github.com/username/project-name/internal/platform/db"
	"github.com/username/project-name/internal/utils"
)

func main() {
	email := flag.String("email", "admin@example.com", "admin email")
	password := flag.String("password", "ChangeMe123!", "admin password")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client, err := db.New(context.Background(), cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	passwords := utils.NewPasswordHasher()
	hash, err := passwords.Hash(*password)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	existing, err := client.Ent.User.Query().Where(user.EmailEQ(*email)).Only(ctx)
	if err == nil {
		if err := client.Ent.User.UpdateOneID(existing.ID).
			SetPasswordHash(hash).
			SetRole("super_admin").
			SetPermissions(permissions.All()).
			SetUpdatedAt(time.Now()).
			Exec(ctx); err != nil {
			log.Fatal(err)
		}
		log.Printf("updated admin user %s", *email)
		return
	}

	created, err := client.Ent.User.Create().
		SetEmail(*email).
		SetPasswordHash(hash).
		SetRole("super_admin").
		SetPermissions(permissions.All()).
		Save(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("created admin user %s (%s)", created.Email, created.ID)
}
