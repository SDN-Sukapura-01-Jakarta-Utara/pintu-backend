package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "generate:migration":
		generateMigration(args)
	case "generate:seeder":
		generateSeeder(args)
	case "generate:model":
		generateModel(args)
	case "generate:repository":
		generateRepository(args)
	case "generate:service":
		generateService(args)
	case "generate:controller":
		generateController(args)
	case "generate:dto":
		generateDTO(args)
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println(`
PINTU Generator Commands:

Usage: go run cmd/main.go <command> [arguments]

Commands:
  generate:migration <name>       Generate new migration file
  generate:seeder <name>          Generate new seeder file
  generate:model <name>           Generate new model file
  generate:repository <model>     Generate new repository file
  generate:service <model>        Generate new service file
  generate:controller <model>     Generate new controller file
  generate:dto <name>             Generate new DTO file

Examples:
  go run cmd/main.go generate:migration create_users_table
  go run cmd/main.go generate:model User
  go run cmd/main.go generate:repository User
  go run cmd/main.go generate:service User
  go run cmd/main.go generate:controller User
  go run cmd/main.go generate:dto User
	`)
}

func generateMigration(args []string) {
	fs := flag.NewFlagSet("migration", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Println("Error: Migration name required")
		fmt.Println("Usage: go run cmd/main.go generate:migration <name>")
		return
	}

	name := fs.Arg(0)
	if err := createMigrationFile(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Migration file created: src/database/migrations/%s.sql\n", name)
}

func generateSeeder(args []string) {
	fs := flag.NewFlagSet("seeder", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Println("Error: Seeder name required")
		fmt.Println("Usage: go run cmd/main.go generate:seeder <name>")
		return
	}

	name := fs.Arg(0)
	if err := createSeederFile(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Seeder file created: src/database/seeders/%s.go\n", name)
}

func generateModel(args []string) {
	fs := flag.NewFlagSet("model", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Println("Error: Model name required")
		fmt.Println("Usage: go run cmd/main.go generate:model <name>")
		return
	}

	name := fs.Arg(0)
	if err := createModelFile(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Model file created: src/modules/models/%s.go\n", name)
}

func generateRepository(args []string) {
	fs := flag.NewFlagSet("repository", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Println("Error: Model name required")
		fmt.Println("Usage: go run cmd/main.go generate:repository <model>")
		return
	}

	name := fs.Arg(0)
	if err := createRepositoryFile(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Repository file created: src/modules/repositories/%s_repository.go\n", toLowerFirst(name))
}

func generateService(args []string) {
	fs := flag.NewFlagSet("service", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Println("Error: Model name required")
		fmt.Println("Usage: go run cmd/main.go generate:service <model>")
		return
	}

	name := fs.Arg(0)
	if err := createServiceFile(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Service file created: src/modules/services/%s_service.go\n", toLowerFirst(name))
}

func generateController(args []string) {
	fs := flag.NewFlagSet("controller", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Println("Error: Model name required")
		fmt.Println("Usage: go run cmd/main.go generate:controller <model>")
		return
	}

	name := fs.Arg(0)
	if err := createControllerFile(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Controller file created: src/modules/controllers/%s_controller.go\n", toLowerFirst(name))
}

func generateDTO(args []string) {
	fs := flag.NewFlagSet("dto", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Println("Error: DTO name required")
		fmt.Println("Usage: go run cmd/main.go generate:dto <name>")
		return
	}

	name := fs.Arg(0)
	if err := createDTOFile(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("DTO file created: src/dtos/%s_dto.go\n", toLowerFirst(name))
}
