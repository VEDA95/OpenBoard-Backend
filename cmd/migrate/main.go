package main

import (
	"VEDA95/open_board/api/internal/config"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	argsWithoutProgram := os.Args[1:]

	if len(argsWithoutProgram) > 2 {
		log.Fatal("Too many arguments...")
	}

	if len(argsWithoutProgram) < 1 {
		log.Fatal("Please provide one of the following arguments: init, up, down, step <count>")
	}

	if err := config.LoadEnvConfigs("./env"); err != nil {
		log.Fatal(err)
	}

	mainAction := argsWithoutProgram[0]

	if mainAction == "init" {
		if len(argsWithoutProgram) < 2 {
			log.Fatal("A name needs to be provided for the migration being created")
		}

		cmd := exec.Command("migrate", "create", "-ext", "sql", "-dir", "./migrations", "-seq", argsWithoutProgram[1])
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(`Migration files for migration: "%s" were successfully created!`, argsWithoutProgram[1])
		return
	}

	dbUrl := os.Getenv("DATABASE_URL")

	if len(dbUrl) == 0 {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	if strings.Contains(dbUrl, "postgres://") {
		dbUrl = strings.Replace(dbUrl, "postgres://", "pgx5://", 1)
	}

	migration, err := migrate.New(
		"file://./migrations",
		dbUrl,
	)

	if err != nil {
		log.Panic(err)
	}

	if mainAction == "up" {
		err := migration.Up()

		if err != nil {
			log.Panic(err)
		}

		fmt.Println("UP migration completed!")
		return
	}

	if mainAction == "down" {
		err := migration.Down()

		if err != nil {
			log.Panic(err)
		}

		fmt.Println("DOWN migration completed!")
		return
	}

	if mainAction == "step" {
		if len(argsWithoutProgram) < 2 {
			log.Fatal("Step count needs to be provided")
		}

		stepCount, err := strconv.ParseInt(argsWithoutProgram[1], 0, 64)

		if err != nil {
			log.Fatal("Unable to parse the value provided for the step count")
		}

		err = migration.Steps(int(stepCount))

		if err != nil {
			log.Panic(err)
		}

		fmt.Println("STEP migration completed!")
		return
	}
}
