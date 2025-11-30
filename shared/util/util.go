package util

import (
	"fmt"
	"io"
	"log"
)

// GetRandomAvatar returns a random avatar URL from the randomuser.me API
func GetRandomAvatar(index int) string {
	return fmt.Sprintf("https://randomuser.me/api/portraits/lego/%d.jpg", index)
}

func CloseAndLog(c io.Closer, name string) {
	if err := c.Close(); err != nil {
		log.Printf("Error closing %s: %v", name, err)
	}
}
