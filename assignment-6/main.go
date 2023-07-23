package main

import (
	"fmt"
	"math/big"
)

const (
	N1 = "179769313486231590772930519078902473361797697894230657273430081157732675805505620686985379449212982959585501387537164015710139858647833778606925583497541085196591615128057575940752635007475935288710823649949940771895617054361149474865046711015101563940680527540071584560878577663743040086340742855278549092581"
	N2 = "648455842808071669662824265346772278726343720706976263060439070378797308618081116462714015276061417569195587321840254520655424906719892428844841839353281972988531310511738648965962582821502504990264452100885281673303711142296421027840289307657458645233683357077834689715838646088239640236866252211790085787877"
	N3 = "720062263747350425279564435525583738338084451473999841826653057981916355690188337790423408664187663938485175264994017897083524079135686877441155132015188279331812309091996246361896836573643119174094961348524639707885238799396839230364676670221627018353299443241192173812729276147530748597302192751375739387929"
)

func main() {
	solve1()
	solve2()
	solve3()
	solve4()
}

func solve1() {
	N, ok := big.NewInt(0).SetString(N1, 10)
	if !ok {
		panic("failed N1")
	}

	A := big.NewInt(0).Sqrt(N)
	A.Add(A, big.NewInt(1))

	X := big.NewInt(0).Set(A)
	X.Mul(X, X)
	X.Sub(X, N)
	X.Sqrt(X)

	P := big.NewInt(0).Set(A)
	P.Sub(A, X)

	Q := big.NewInt(0).Set(A)
	Q.Add(A, X)

	PQ := big.NewInt(0).Mul(P, Q)

	if PQ.String() == N.String() {
		fmt.Println(P.String())
	}
}

func solve2() {
	N, ok := big.NewInt(0).SetString(N2, 10)
	if !ok {
		panic("failed N2")
	}

	sqrt := big.NewInt(0).Sqrt(N)
	sqrt.Add(sqrt, big.NewInt(1))

	// A - sqrt(N) < 2^20
	// -A + sqrt(N) > 2^20
	// sqrt(N) > A - 2^20

	for i := 0; i < 1<<20; i++ {
		A := big.NewInt(0).Set(sqrt)
		A.Add(A, big.NewInt(int64(i)))

		X := big.NewInt(0).Set(A)
		X.Mul(X, X)
		X.Sub(X, N)
		X.Sqrt(X)

		P := big.NewInt(0).Set(A)
		P.Sub(A, X)

		Q := big.NewInt(0).Set(A)
		Q.Add(A, X)

		PQ := big.NewInt(0).Mul(P, Q)

		if PQ.String() == N.String() {
			fmt.Println(P.String())
			break
		}
	}
}

func solve3() {
	N, ok := big.NewInt(0).SetString(N3, 10)
	if !ok {
		panic("failed N3")
	}

	N_24 := big.NewInt(0).Mul(N, big.NewInt(24))

	A := big.NewInt(0).Set(N_24)
	A.Sqrt(A)
	A.Add(A, big.NewInt(1))

	X := big.NewInt(0).Set(A)
	X.Mul(X, X)    // A^2
	X.Sub(X, N_24) // - 24
	X.Sqrt(X)

	Q := big.NewInt(0).Sub(A, X)
	Q.Div(Q, big.NewInt(6))

	P := big.NewInt(0).Add(A, X)
	P.Div(P, big.NewInt(4))

	R := big.NewInt(0).Mul(P, Q)

	if R.String() == N.String() {
		fmt.Println(Q.String())
	}
}

func solve4() {
	N, ok := big.NewInt(0).SetString(N1, 10)
	if !ok {
		panic("failed N1")
	}

	A := big.NewInt(0).Sqrt(N)
	A.Add(A, big.NewInt(1))

	X := big.NewInt(0).Set(A)
	X.Mul(X, X)
	X.Sub(X, N)
	X.Sqrt(X)

	P := big.NewInt(0).Set(A)
	P.Sub(A, X)

	Q := big.NewInt(0).Set(A)
	Q.Add(A, X)

	challenge := "22096451867410381776306561134883418017410069787892831071731839143676135600120538004282329650473509424343946219751512256465839967942889460764542040581564748988013734864120452325229320176487916666402997509188729971690526083222067771600019329260870009579993724077458967773697817571267229951148662959627934791540"
	encryption := "65537"

	E, ok := big.NewInt(0).SetString(encryption, 10)
	if !ok {
		panic("failed E")
	}

	C, ok := big.NewInt(0).SetString(challenge, 10)
	if !ok {
		panic("failed C")
	}

	// compute phi.
	left := big.NewInt(0)
	left.Sub(P, big.NewInt(1))

	right := big.NewInt(0)
	right.Sub(Q, big.NewInt(1))

	PHI := big.NewInt(0).Mul(left, right)

	D := big.NewInt(0).ModInverse(E, PHI)

	// decrypt
	C.Exp(C, D, N)

	msg := C.Bytes()
	msg = msg[len(msg)-29:]
	fmt.Println(string(msg))
}
