package rsa

import "math/big"

//Represent a RSA key
type Key struct {
	Exp int64
	Mod int64
}

//Return maw value between 2 int
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

//Return gcd between 2 int
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

//Split a message in bytes array of equal sizes
func splitMessage(bytes []byte, size int) (chunks [][]byte) {
	for len(bytes) > size {
		chunks = append(chunks, bytes[:size])
		bytes = bytes[size:]
	}
	if len(bytes) > 0 {
		chunks = append(chunks, bytes)
	}
	return chunks
}

//Return a public and a private key calculated from 2 random primes
func CalculateKeys(size int) (publicKey Key, privateKey Key) {
	p := CalculatePrime(size)
	q := CalculatePrime(size)
	n := p * q
	m := (p - 1) * (q - 1)
	publicKey.Mod = n
	privateKey.Mod = n

	for i := max(p, q) + 2; i < m; i += 2 {
		if gcd(i, m) == 1 {
			publicKey.Exp = i
			break
		}
	}

	for i := max(p, q) + 2; i < n; i++ {
		if publicKey.Exp*i%m == 1 {
			privateKey.Exp = i
			break
		}
	}
	return
}

//encrypt the giver message with the given key
func Encrypt(msg []byte, publicKey Key) (out []byte) {
	e := big.NewInt(publicKey.Exp)
	n := big.NewInt(publicKey.Mod)
	size := (n.BitLen() + 7) / 8

	//3 is arbitrary
	splitedMsg := splitMessage(msg, 3)

	for i := 0; i < len(splitedMsg); i++ {
		chunck := new(big.Int)
		chunck = chunck.SetBytes(splitedMsg[i])

		chunck = chunck.Exp(chunck, e, n)

		t := chunck.Bytes()

		if len(t) < size {

			//LeftPad with 0 if the output is to small
			temp := make([]byte, size)
			copy(temp[size-len(t):], t)
			t = temp
		}
		out = append(out, t...)

	}
	return
}

//decrypt the giver message with the given key
func Decrypt(msg []byte, privateKey Key) (out []byte) {
	d := big.NewInt(privateKey.Exp)
	n := big.NewInt(privateKey.Mod)
	size := (n.BitLen() + 7) / 8

	splitedMsg := splitMessage(msg, size)

	for i := 0; i < len(splitedMsg); i++ {
		chunck := new(big.Int)
		chunck = chunck.SetBytes(splitedMsg[i])
		chunck = chunck.Exp(chunck, d, n)
		for _, v := range chunck.Bytes() {
			out = append(out, v)
		}
	}
	return out

}
