package goo

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

	b, err := enc.Decode("940c2f8fea0d366ddaef4bc26661abeebf3bce75f7ee2effffcf02037467452a2e37766aec0a175375513655d023d4f7ada15625a282a0ae5f352586c4c57dad63b08ffaff05ee531cddbc724d26accd")
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(string(b))
}
