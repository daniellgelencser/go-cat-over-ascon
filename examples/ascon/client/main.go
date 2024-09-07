package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go-attested-coap-over-ascon/v3/ascon"
)

func main() {
	co, err := ascon.Dial("localhost:5688")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	path := "/a"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// for {
	// ctx, cancel = context.WithTimeout(context.Background(), time.Second)

	fmt.Println("\nSentGet")
	resp, err := co.Get(ctx, path)
	// fmt.Println("SentGet")
	// resp, err = co.Get(ctx, path)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	body, err := resp.ReadBody()
	log.Printf("\nResponse payload: \n%v\nbody: \n%v", resp.String(), string(body))
	time.Sleep(2 * time.Second)
	// }
}
