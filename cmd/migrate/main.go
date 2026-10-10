package main

import (
	"fmt"
	"log"
	"os"
	"todo-service/internal/config"
	"todo-service/internal/storage/postgres"
)

func main() {
	if len(os.Args) < 2 {
		_, err := fmt.Fprintln(os.Stderr, "usage: migrate <up|down|version>")
		if err != nil {
			log.Println("Migration error:", err)
		}

		os.Exit(2)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := postgres.RunMigrations(cfg.DatabaseUrl, os.Args[1]); err != nil {
		log.Fatalf("migrate %s: %v", os.Args[1], err)
	}
}
