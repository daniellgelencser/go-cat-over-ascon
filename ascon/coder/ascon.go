package coder

import (
	"bytes"
)

type Ascon struct {
	constants [12]block
	difusion  [5][2]int
	state     [5]block
	nonce     [2]block
	key       [2]block
	iv        block

	a int
	b int
}

func makeAscon(keyBytes [16]byte, nonceBytes [16]byte) Ascon {

	//initialization vector: 80 40 0c 06 00 00 00 00
	var iv block = block{0: 0x80, 1: 0x40, 2: 0x0C, 3: 0x06}

	keyPtr0 := (*block)(keyBytes[0:8])
	keyPtr1 := (*block)(keyBytes[8:16])
	noncePtr0 := (*block)(nonceBytes[0:8])
	noncePtr1 := (*block)(nonceBytes[8:16])

	ascon := Ascon{
		constants: [12]block{
			{7: 0xf0},
			{7: 0xe1},
			{7: 0xd2},
			{7: 0xc3},
			{7: 0xb4},
			{7: 0xa5},
			{7: 0x96},
			{7: 0x87},
			{7: 0x78},
			{7: 0x69},
			{7: 0x5a},
			{7: 0x4b},
		},
		difusion: [5][2]int{
			{19, 28},
			{61, 39},
			{01, 06},
			{10, 17},
			{07, 41},
		},
		state: [5]block{
			iv,
			*keyPtr0,
			*keyPtr1,
			*noncePtr0,
			*noncePtr1,
		},
		nonce: [2]block{*noncePtr0, *noncePtr1},
		key:   [2]block{*keyPtr0, *keyPtr1},
		iv:    iv,

		a: 12,
		b: 6,
	}

	return ascon
}

// Addition of Constants
func (ascon *Ascon) addConstant(i int, x int) {
	ascon.state[2].DXOR(ascon.constants[12-x+i])
}

// Substitution Layer - 64 "parallel" sbox
func (ascon *Ascon) sbox() {
	var temp [5]block

	ascon.state[0].DXOR(ascon.state[4]) // x0 ^= x4;
	ascon.state[4].DXOR(ascon.state[3]) // x4 ^= x3;
	ascon.state[2].DXOR(ascon.state[1]) // x2 ^= x1;

	temp[0] = ascon.state[0].NOT() // t0 = x0; t0 =~ t0;
	temp[1] = ascon.state[1].NOT() // t1 = x1; t1 =~ t1;
	temp[2] = ascon.state[2].NOT() // t2 = x2; t2 =~ t2;
	temp[3] = ascon.state[3].NOT() // t3 = x3; t3 =~ t3;
	temp[4] = ascon.state[4].NOT() // t4 = x4; t4 =~ t4;

	temp[0].DAND(ascon.state[1]) // t0 &= x1;
	temp[1].DAND(ascon.state[2]) // t1 &= x2;
	temp[2].DAND(ascon.state[3]) // t2 &= x3;
	temp[3].DAND(ascon.state[4]) // t3 &= x4;
	temp[4].DAND(ascon.state[0]) // t4 &= x0;

	ascon.state[0].DXOR(temp[1]) // t0 ^= x1;
	ascon.state[1].DXOR(temp[2]) // t1 ^= x2;
	ascon.state[2].DXOR(temp[3]) // t2 ^= x3;
	ascon.state[3].DXOR(temp[4]) // t3 ^= x4;
	ascon.state[4].DXOR(temp[0]) // t4 ^= x0;

	ascon.state[1].DXOR(ascon.state[0]) // x1 ^= x0;
	ascon.state[0].DXOR(ascon.state[4]) // x0 ^= x4;
	ascon.state[3].DXOR(ascon.state[2]) // x3 ^= x2;

	ascon.state[2].DNOT() // x2 =~ x2;
}

// Linear Diffusion Layer
func (ascon *Ascon) diffuse() {
	for i := 0; i < 5; i++ {
		temp0 := ascon.state[i].ROTATE(ascon.difusion[i][0])
		temp1 := ascon.state[i].ROTATE(ascon.difusion[i][1])
		ascon.state[i].DXOR(temp0.XOR(temp1))
	}
}

// Perform ASCON Permutations - if not ready "a" times , else "b" times
func (ascon *Ascon) permutation(x int) {
	for i := 0; i < x; i++ {
		ascon.addConstant(i, x)
		ascon.sbox()
		ascon.diffuse()
	}
}

func (ascon *Ascon) initialize() {
	ascon.permutation(ascon.a)
	ascon.state[3].DXOR(ascon.key[0])
	ascon.state[4].DXOR(ascon.key[1])
}

// Process Whole Padded Plaintext
func (ascon *Ascon) processPlaintext(plaintext []byte) []byte {

	var cyphertext []byte
	textLen := len(plaintext)
	l := textLen % BlockBytes
	lastByte := textLen

	if l > 0 {
		// do padding P||1||0 r-1(|P| % r)
		tmp := make([]byte, BlockBytes-l)
		tmp[0] = 0x80
		plaintext = append(plaintext, tmp...)
		lastByte -= l
	} else {
		lastByte -= BlockBytes
	}

	for i := 0; i < lastByte; i += BlockBytes {

		// Sr ← Sr ⊕ Pi - xor and store
		ascon.state[0].DXOR(*(*block)(plaintext[i : i+BlockBytes]))

		//Ci ← Sr - append
		cyphertext = append(cyphertext, ascon.state[0][:]...)

		//S ← pb (S) - permutate
		ascon.permutation(ascon.b)
	}

	// Sr ← Sr ⊕ Pt - xor and store
	ascon.state[0].DXOR(*(*block)(plaintext[lastByte:]))

	// Ct ← \Sr/l - append
	cyphertext = append(cyphertext, ascon.state[0][:]...)

	// truncate on return
	return cyphertext[:textLen]
}

func (ascon *Ascon) processCyphertext(cyphertext []byte) []byte {

	var plaintext []byte
	textLen := len(cyphertext)
	l := textLen % BlockBytes
	lastByte := textLen

	if l > 0 {
		lastByte -= l
	} else {
		lastByte -= BlockBytes
	}

	for i := 0; i < lastByte; i += BlockBytes {

		cyphertextBlock := *(*block)(cyphertext[i : i+BlockBytes])

		// Pi ← Sr ⊕ Ci - xor and store
		tmp := ascon.state[0].XOR(cyphertextBlock)
		plaintext = append(plaintext, tmp[:]...)

		// S ← Ci || Sc - replace
		ascon.state[0] = cyphertextBlock

		// S ← pb (S) - permute
		ascon.permutation(ascon.b)
	}

	if l == 0 {

		cyphertextBlock := *(*block)(cyphertext[lastByte:])

		// Pi ← Sr ⊕ Ci - xor and store
		tmp := ascon.state[0].XOR(cyphertextBlock)
		plaintext = append(plaintext, tmp[:]...)

		// S ← Ci || Sc - replace
		ascon.state[0] = cyphertextBlock

		return plaintext
	}

	// Pt ← \Sr/l ⊕ Ct - truncate (first l bits) and xor
	tmp1 := ascon.state[0].XORP(cyphertext[lastByte:], l)
	plaintext = append(plaintext, tmp1...)

	// S ← Ct || (/Sr\r−l ⊕ (1 || 0r−1−l )) || Sc truncate (last r-l bits) and xor
	copy(ascon.state[0][:], cyphertext[lastByte:])
	ascon.state[0][l] ^= 0x80

	return plaintext
}

// Return TAG []byte
func (ascon *Ascon) finalize() []byte {

	//S ← pa (S ⊕ (0r || K || 0c−k ))
	ascon.state[0].DXOR(ascon.key[0])
	ascon.state[1].DXOR(ascon.key[1])
	ascon.permutation(ascon.a)

	ascon.state[3].DXOR(ascon.key[0])
	ascon.state[4].DXOR(ascon.key[1])

	return append(ascon.state[3][:], ascon.state[4][:]...)
}

// Returns Cyphertext and Tag
func Encrypt(key []byte, nonce []byte, plaintext []byte) ([]byte, []byte) {

	ascon := makeAscon(*(*[16]byte)(key), *(*[16]byte)(nonce))

	ascon.initialize()

	cyphertext := ascon.processPlaintext(plaintext)

	tag := ascon.finalize()

	return cyphertext, tag
}

// Returns Plaintext if tag matches
func Decrypt(key []byte, nonce []byte, cyphertext []byte, tag []byte) []byte {

	ascon := makeAscon(*(*[16]byte)(key), *(*[16]byte)(nonce))

	ascon.initialize()

	plaintext := ascon.processCyphertext(cyphertext)

	tag1 := ascon.finalize()

	if !bytes.Equal(tag, tag1) {
		return nil
	}

	return plaintext
}
