package goo_http

import (
	"log"
	"testing"
)

func TestEncryption_Encode(t *testing.T) {
	enc := &Encryption{
		Key:    "65bea91a9a90cf43",
		Secret: "231a05f6faf619fb",
	}

	str, err := enc.Encode([]byte(`12345`))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(str)
}

func TestEncryption_Decode(t *testing.T) {
	enc := &Encryption{
		Key:    "65bea91a9a90cf43",
		Secret: "231a05f6faf619fb",
	}

	b, err := enc.Decode("")
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(string(b))
}
