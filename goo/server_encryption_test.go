package goo

import (
	"log"
	"testing"
)

func TestEncryption_Decode(t *testing.T) {
	enc := &Encryption{
		Key:    "1a63a991a95bea49e35e0c6ab1da66fd",
		Secret: "1a6395991ae3bea49ab1da65e0c6a6fd",
	}

	str, err := enc.Encode([]byte(""))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(">>1>>", str)

	enc2 := &Encryption{
		Key:    "29a9916a6395bea4aae35e0cb1da66fd",
		Secret: "2a6e35e0c69a991a395bea4ab1da66fd",
	}

	b, err := enc2.Decode(str)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(">>2>>", string(b))
}
