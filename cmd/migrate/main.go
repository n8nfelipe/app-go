package main

import (
	"fmt"
	"log"
	"os"

	"github.com/n8nfelipe/app-go.git/db"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		if err := db.Up(); err != nil {
			log.Fatalf("migrate up: %v", err)
		}
		fmt.Println("migrations aplicadas")
	case "down":
		if err := db.Down(); err != nil {
			log.Fatalf("migrate down: %v", err)
		}
		fmt.Println("migrations revertidas")
	case "version":
		version, dirty, err := db.Version()
		if err != nil {
			log.Fatalf("migrate version: %v", err)
		}
		if version == 0 {
			fmt.Println("nenhuma migration aplicada")
			return
		}
		fmt.Printf("versão: %d (dirty: %v)\n", version, dirty)
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "uso: go run ./cmd/migrate <up|down|version>")
}
