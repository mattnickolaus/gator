package main

import (
	"fmt"

	"github.com/mattnickolaus/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	err = c.SetUser("mattn")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	c, err = config.Read()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("Config:\n %+v\n", c)
}
