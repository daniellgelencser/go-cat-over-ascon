package main

import (
	"fmt"

	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/ascon/coder"
)

func main() {
	var buf []byte = []byte("mattis aliquam faucibus purus in massa tempor nec feugiat nisl pretium fusce id velit ut tortor pretium viverra suspendisse potenti nullam ac tortor vitae purus faucibus ornare suspendisse sed nisi lacus sed viverra tellus in hac habitasse platea dictumst vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt lobortis feugiat vivamus at augue eget arcu dictum varius duis at consectetur lorem donec massa sapien faucibus et molestie ac feugiat sed lectus vestibulum mattis ullamcorper velit sed ullamcorper morbi tincidunt ornare massa eget egestas purus viverra accumsan in nisl nisi scelerisque eu ultrices vitae auctor eu augue ut lectus arcu bibendum mattis aliquam faucibus purus in massa tempor nec feugiat nisl pretium fusce id velit ut tortor pretium viverra suspendisse potenti nullam ac tortor vitae purus faucibus ornare suspendisse sed nisi lacus sed viverra tellus in hac habitasse platea dictumst vestibulum rhoncus est pellentesque elit ullamcorper dignissim cras tincidunt lobortis feugiat vivamus at augue eget arcu dictum varius duis at consectetur lorem donec massa sapien faucibus et molestie ac feugiat sed lectus vestibulum")

	// []byte{
	// 	0x12,
	// 	0x34,
	// 	0x56,
	// 	0x78,
	// 	0x90,
	// 	0xab,
	// 	0xcd,
	// 	0xef,
	// 	0x08,
	// 	0x21,
	// 	0x87,
	// 	0x12,
	// 	0x34,
	// 	0x56,
	// 	0x78,
	// 	0x90,
	// 	0xab,
	// 	0xcd,
	// 	// 0xef,
	// 	// 0x08,
	// 	// 0x21,
	// 	// 0x87,
	// 	// 0x99,
	// 	// 0x99,
	// }

	var key []byte = //make([]byte, 16)
	[]byte{
		0x12,
		0x34,
		0x12,
		0x34,
		0x12,
		0x34,
		0x12,
		0x34,
		0x12,
		0x34,
		0x12,
		0x34,
		0x12,
		0x34,
		0x12,
		0x34,
	}

	var nonce []byte = //make([]byte, 16)
	[]byte{
		0xab,
		0xcd,
		0xab,
		0xcd,
		0xab,
		0xcd,
		0xab,
		0xcd,
		0xab,
		0xcd,
		0xab,
		0xcd,
		0xab,
		0xcd,
		0xab,
		0xcd,
	}

	cyphertext, tag := coder.Encrypt(key, nonce, buf)

	plaintext := coder.Decrypt(key, nonce, cyphertext, tag)

	fmt.Printf("\nCyphertext: \n%x\n", cyphertext)
	fmt.Printf("\nBuffer (%d mod 8 = %d):	\n%s\n", len(buf), len(buf)%coder.BlockBytes, buf)
	fmt.Printf("\nPlaintext: \n%s\n", plaintext)
	fmt.Printf("\nTag: %X\n", tag)
}
