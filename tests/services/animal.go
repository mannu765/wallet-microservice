package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Animal interface
type Animal interface {
	Eat()
	Move()
	Speak()
}

// Cow struct
type Cow struct{}

func (c Cow) Eat()   { fmt.Println("grass") }
func (c Cow) Move()  { fmt.Println("walk") }
func (c Cow) Speak() { fmt.Println("moo") }

// Bird struct
type Bird struct{}

func (b Bird) Eat()   { fmt.Println("worms") }
func (b Bird) Move()  { fmt.Println("fly") }
func (b Bird) Speak() { fmt.Println("peep") }

// Snake struct
type Snake struct{}

func (s Snake) Eat()   { fmt.Println("mice") }
func (s Snake) Move()  { fmt.Println("slither") }
func (s Snake) Speak() { fmt.Println("hsss") }

func main() {
	animals := make(map[string]Animal) // store animals by name
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ") // prompt
		if !scanner.Scan() {
			break
		}

		parts := strings.Fields(scanner.Text())
		if len(parts) != 3 {
			fmt.Println("Invalid command. Format: newanimal <name> <type> OR query <name> <info>")
			continue
		}

		command, name, arg := parts[0], parts[1], parts[2]

		switch command {
		case "newanimal":
			var a Animal
			switch arg {
			case "cow":
				a = Cow{}
			case "bird":
				a = Bird{}
			case "snake":
				a = Snake{}
			default:
				fmt.Println("Unknown animal type")
				continue
			}
			animals[name] = a
			fmt.Println("Created it!")

		case "query":
			a, ok := animals[name]
			if !ok {
				fmt.Println("Animal not found")
				continue
			}
			switch arg {
			case "eat":
				a.Eat()
			case "move":
				a.Move()
			case "speak":
				a.Speak()
			default:
				fmt.Println("Unknown query")
			}

		default:
			fmt.Println("Unknown command")
		}
	}
}
