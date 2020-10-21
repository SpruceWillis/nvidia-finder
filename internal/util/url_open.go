package util

import (
	"log"
	"os/exec"
)

// OpenURL open resource in the OS default program
func OpenURL(url string) {
	cmd := exec.Command("open", url)
	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
}
