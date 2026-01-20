package cli

import (
	"fmt"

	"github.com/sembraniteam/esc/internal/config"
)

func printConnection(c *config.Connection) {
	fmt.Printf("Name        : %s\n", c.Name)
	fmt.Printf("Host        : %s\n", c.Host)
	fmt.Printf("User        : %s\n", c.User)
	fmt.Printf("Port        : %d\n", c.Port)
	fmt.Println()

	fmt.Println("Auth")
	fmt.Printf("  Type      : %s\n", c.Auth.Type)

	switch c.Auth.Type {
	case "key":
		if c.Auth.KeyPath == "" {
			fmt.Printf("  Key Path  : %s\n", c.Auth.KeyPath)
		}
	case "password":
		fmt.Println("  Password  : (from env)")
	case "agent":
		fmt.Println("  SSH Agent : enabled")
	}

	if len(c.Tags) > 0 {
		fmt.Println()
		fmt.Println("Tags")
		for _, t := range c.Tags {
			fmt.Printf("  - %s\n", t)
		}
	}
}
