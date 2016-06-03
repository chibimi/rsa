// A concurrent prime sieve

package main

import (
	"fmt"
	"math/big"
)

type key struct {
	exp int64
	mod int64
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

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

func calculateKeys(p, q int64) (publicKey key, privateKey key) {
	n := p * q
	m := (p - 1) * (q - 1)
	publicKey.mod = n
	privateKey.mod = n

	for i := max(p, q) + 2; i < m; i += 2 {
		if gcd(i, m) == 1 {
			publicKey.exp = i
			break
		}
	}

	for i := max(p, q) + 2; i < n; i++ {
		if publicKey.exp*i%m == 1 {
			privateKey.exp = i
			break
		}
	}
	return
}

func encrypt(msg []byte, publicKey key) (out []byte) {
	e := big.NewInt(publicKey.exp)
	n := big.NewInt(publicKey.mod)
	size := 5

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
	fmt.Println(out)
	return
}

func decrypt(msg []byte, privateKey key) (out []byte) {
	d := big.NewInt(privateKey.exp)
	n := big.NewInt(privateKey.mod)
	size := 5

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

// The prime sieve: Daisy-chain Filter processes.
func main() {

	p := calculatePrime(4)
	q := calculatePrime(4)

	publicKey, privateKey := calculateKeys(p, q)
	fmt.Println(publicKey, privateKey)
	code := encrypt([]byte("Emilie"), publicKey)
	decode := decrypt(code, privateKey)
	fmt.Println(string(decode))
}
