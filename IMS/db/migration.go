package db

import (
	"log"
	"os"

	"github.com/omniful/go_commons/db/sql/migration"
)

func NativeMigrationUp() {
	migrator, err := migration.InitializeMigrate("file://migrations/", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}
	err = migrator.Up()
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied successfully!")
}

func NativeMigrationDown() {
	migrator, err := migration.InitializeMigrate("file://migrations/", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}
	err = migrator.Down()
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied successfully!")
}

func CleanDirtyMigration() {
	migrator, err := migration.InitializeMigrate("file://migrations/", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Init failed: %v", err)
	}

	err = migrator.ForceVersion(1)
	if err != nil {
		log.Fatalf("Failed to force clean version: %v", err)
	}

	log.Println("Forced to clean version 1")
}
