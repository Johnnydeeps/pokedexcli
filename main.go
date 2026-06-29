package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			if err := scanner.Err(); err != nil {
				return
			}
			return
		}

		text := scanner.Text()

		words := cleanInput(text)
		if len(words) == 0 {
			continue
		}

		cliInput := words[0]
		fmt.Printf("Your command was: %s\n", cliInput)
	}
}
