package util

import (
	"log"
	"os/exec"
)

// OpenUrl open resource in the OS default program
func OpenUrl(url string) {
	cmd := exec.Command("open", url)
	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
}
