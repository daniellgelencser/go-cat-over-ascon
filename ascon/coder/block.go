package coder

const BlockBytes = 8

type block [BlockBytes]byte

// Constructive SHIFT RIGHT
func (x block) SHIFTR(n int) block {
	var temp [7]byte
	for i := 0; i < 7; i++ {
		temp[i] = x[i] << (64 - n)
	}

	x[0] = x[0] >> n
	for i := 1; i < 8; i++ {
		x[i] = x[i] >> n
		x[i] ^= temp[i-1]
	}

	return x
}

// Constructive SHIFT LEFT
func (x block) SHIFTL(n int) block {
	var temp [7]byte
	for i := 1; i < 8; i++ {
		temp[i-1] = x[i] >> (64 - n)
	}

	for i := 0; i < 7; i++ {
		x[i] = x[i] << n
		x[i] ^= temp[i]
	}
	x[7] = x[7] >> n

	return x
}

// Constructive ROTATE
func (x block) ROTATE(l int) block {
	right := x.SHIFTR(l)
	left := x.SHIFTL(64 - l)
	return right.XOR(left)
}

// Constructive XOR
func (x block) XOR(y block) block {
	for i := 0; i < 8; i++ {
		x[i] ^= y[i]
	}
	return x
}

// Constructive XOR with partial block, truncate to most signifficant bytes
func (x block) XORP(y []byte, len int) []byte {
	for i := 0; i < len; i++ {
		x[i] ^= y[i]
	}
	return x[:len]
}

// // Constructive XOR (most signifficat bytes of y)
// func (x blockPart) XOR(y blockPart, len int) blockPart {
// 	for i := 0; i < len; i++ {
// 		x[i] ^= y[i]
// 	}
// 	return x
// }

// Destructive XOR
func (x *block) DXOR(y block) {
	for i := 0; i < 8; i++ {
		x[i] ^= y[i]
	}
}

// Destructive XOR (most signifficat bytes of y)
// func (x *blockPart) DXOR(y blockPart, len int) {
// 	for i := 0; i < len; i++ {
// 		(*x)[i] ^= y[i]
// 	}
// }

// Constructive NOT
func (x block) NOT() block {
	for i := 0; i < 8; i++ {
		x[i] = ^x[i]
	}
	return x
}

// Destructive NOT
func (x *block) DNOT() {
	for i := 0; i < 8; i++ {
		x[i] = ^x[i]
	}
}

// Constructive AND
func (x block) AND(y block) block {
	for i := 0; i < 8; i++ {
		x[i] &= y[i]
	}
	return x
}

// Destructive AND
func (x *block) DAND(y block) {
	for i := 0; i < 8; i++ {
		x[i] &= y[i]
	}
}
