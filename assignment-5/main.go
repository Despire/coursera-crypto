package main

import (
	"fmt"
	"math/big"
)

func main() {
	// P the prime.
	P, ok := big.NewInt(0).SetString("13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006084171", 10)
	if !ok {
		panic("failed P")
	}

	// G the generator.
	G, ok := big.NewInt(0).SetString("11717829880366207009516117596335367088558084999998952205599979459063929499736583746670572176471460312928594829675428279466566527115212748467589894601965568", 10)
	if !ok {
		panic("failed G")
	}

	H, ok := big.NewInt(0).SetString("3239475104050450443565264378728065788649097520952449527834792452971981976143292558073856937958553180532878928001494706097394108577585732452307673444020333", 10)
	if !ok {
		panic("failed H")
	}

	// from the description of the assignment we're solving
	// B = 2^20
	// h / g^x1 = (g^B)^x0 => where (x0, x1) e {0 ... B-1}
	// build the left side 2^20 operations
	// continue with the right side.

	cache := make(map[string]int, 1<<20)

	for x1 := 0; x1 < 1<<20; x1++ {
		// g^x1
		r := big.NewInt(0).Exp(G, big.NewInt(int64(x1)), P)
		// 1/g^x1
		r.ModInverse(r, P)
		// h * 1 / g^x1
		r.Mul(H, r)
		// r % P
		r.Mod(r, P)
		cache[r.String()] = x1
	}

	// B = 2^20
	B := big.NewInt(2)
	B.Exp(B, big.NewInt(20), P)
	// g^B
	g := big.NewInt(0).Exp(G, B, P)

	for x0 := 0; x0 < 1<<20; x0++ {
		// calc (g^B)^x0
		r := big.NewInt(0).SetBytes(g.Bytes())
		r.Exp(r, big.NewInt(int64(x0)), P)

		if x1, ok := cache[r.String()]; ok {
			fmt.Printf("x0: %v\tx1: %v\n", x0, x1)
		}
	}

	// found
	// x0: 357984
	// x1: 787046
	x0 := big.NewInt(357984)
	x1 := big.NewInt(787046)

	// x = x0*B +x1
	x := big.NewInt(0).Mul(x0, B)
	x.Add(x, x1)

	fmt.Println(x.String())

	// g^x
	got := big.NewInt(0).Exp(G, x, P)

	fmt.Println(got.String())
	fmt.Println(H.String())
}
