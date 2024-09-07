package coder

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

func Printer() {

	a := RandomBytes(32)
	b := RandomBytes(32)

	fmt.Printf("Alice private:	%X\n", a)
	fmt.Printf("Bob private:	%X\n", b)

	aPub := ComputePublicKey(a)
	bPub := ComputePublicKey(b)

	fmt.Printf("Alice public:	%X\n", aPub)
	fmt.Printf("Bob public:	%X\n", bPub)

	aSh, _ := DeriveSharedKey(a, bPub)
	bSh, _ := DeriveSharedKey(b, aPub)

	if bytes.Equal(aSh, bSh) {
		fmt.Printf("Shared key:	%X\n", aSh)
		return
	}

	fmt.Printf("Alice shared:	%X\n", aSh)
	fmt.Printf("Bob shared:	%X\n", bSh)
}

// computePublicKey computes the public key corresponding to the given private key.
func ComputePublicKey(privateKey []byte) []byte {
	var publicKey []byte = make([]byte, 32)
	curve25519.ScalarBaseMult((*[32]byte)(publicKey), (*[32]byte)(privateKey))
	return publicKey
}

// deriveSharedKey derives a shared key between a private key and a public key.
func DeriveSharedKey(privateKey []byte, publicKey []byte) ([]byte, error) {
	sharedKey, err := curve25519.X25519(privateKey, publicKey)
	if err != nil {
		return sharedKey, err
	}
	return sharedKey, nil
}
