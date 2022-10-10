package goo

import (
	"log"
	"testing"
)

func TestEncryption_Decode(t *testing.T) {
	enc := &Encryption{
		Key:    "",
		Secret: "",
	}

	b, err := enc.Decode("")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(b))
}
