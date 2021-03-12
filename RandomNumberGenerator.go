package DCP

import (
	"fmt"
	"math"
	"math/big"
)
// https://play.golang.org/p/XXiTIpRBAu
const max = 1000000

func RandomLargePrimeGenerator() *big.Int {
	var x, y, n int
	nsqrt := math.Sqrt(max)

	is_prime := [max]bool{}

	for x = 1; float64(x) <= nsqrt; x++ {
		for y = 1; float64(y) <= nsqrt; y++ {
			n = 4*(x*x) + y*y
			if n <= max && (n%12 == 1 || n%12 == 5) {
				is_prime[n] = !is_prime[n]
			}
			n = 3*(x*x) + y*y
			if n <= max && n%12 == 7 {
				is_prime[n] = !is_prime[n]
			}
			n = 3*(x*x) - y*y
			if x > y && n <= max && n%12 == 11 {
				is_prime[n] = !is_prime[n]
			}
		}
	}

	for n = 5; float64(n) <= nsqrt; n++ {
		if is_prime[n] {
			for y = n * n; y < max; y += n * n {
				is_prime[y] = false
			}
		}
	}

	is_prime[2] = true
	is_prime[3] = true

	primes := make([]int, 0, 1270606)
	for x = 0; x < len(is_prime)-1; x++ {
		if is_prime[x] {
			primes = append(primes, x)
		}
	}

	// primes is now a slice that contains all the
	// primes numbers up to max

	// let's print them
	for _, x := range primes {
		fmt.Println(x)
	}

	return big.NewInt(0)
}