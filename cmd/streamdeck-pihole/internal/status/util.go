package status

import (
	"fmt"
	"log"
)

func handleError(msg string, a ...interface{}) error {
	err := fmt.Errorf(msg, a...)
	log.Printf("[error]: %v", err)
	return err
}