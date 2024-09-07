package attestation

// "github.com/google/go-tpm/legacy/tpm2"
// "github.com/google/go-tpm/legacy/tpm2"

// "github.com/google/go-tpm/legacy/tpm2"

// func EstablishTrust() {
// 	// Open a connection to the TPM device
// 	tpm, err := tpm2.OpenTPM("/dev/tpm0")
// 	if err != nil {
// 		log.Fatalf("Failed to open TPM: %v", err)
// 	}
// 	// sim, err := simulator.Get()
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to start TPM simulator: %v", err)
// 	// }
// 	// defer sim.Close()

// 	pcrSel := tpm2.PCRSelection{Hash: tpm2.AlgAES, PCRs: []int{0}}

// 	// Record the current PCR value and timestamp as the trust anchor
// 	trustAnchorPCRValue, err := tpm2.ReadPCRs(tpm, pcrSel)
// 	if err != nil {
// 		fmt.Printf("Failed to read PCR: %v\n", err)
// 		return
// 	}

// 	trustAnchorTimestamp := time.Now()

// 	fmt.Printf("%v %v\n", trustAnchorPCRValue, trustAnchorTimestamp)

// 	// Register the trust anchor with the Verifier (e.g., store it in a secure location)
// 	// ...

// 	// Continue normal operations
// 	// ...
// }
