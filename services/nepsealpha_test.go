package services

import (
	"fmt"
	"log"
	"testing"
)

func TestNepseAlpha(t *testing.T) {
	n := NewNepseAlpha("", "")

	p, err := n.Portfolio()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%+v\n", p)
}
