package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Result represents a verified center, its twin prime pair, and its additive "ancestry"
type Result struct {
	Center       int      `json:"center"`
	PrimePair    [2]int   `json:"primePair"`
	Combinations [][2]int `json:"combinations"`
}

func main() {
	// 1. Setup Flags
	maxFlag := flag.Int("max", 1000, "The maximum number to check up to")
	jsonFile := flag.String("db", "conjecture.json", "JSON file for additive genealogy")
	primeFile := flag.String("primes", "primes.txt", "Text file for persisted primes")
	optimized := flag.Bool("optimized", false, "If true, finds only the pair closest to midpoint and saves to conjecture_optimized.json")
	flag.Parse()

	limit := *maxFlag
	start := time.Now()

	// 2. Load Data from Persistence
	primes, primeMap, lastPrime := loadPrimes(*primeFile)
	results := loadResults(*jsonFile)

	// 3. Build Pool and Map for lookups
	pool := []int{}
	provenCenters := make(map[int]bool)
	for _, r := range results {
		pool = append(pool, r.Center)
		provenCenters[r.Center] = true
	}
	sort.Ints(pool)

	lastCenter := 2
	if len(pool) > 0 {
		lastCenter = pool[len(pool)-1]
	}

	if lastCenter >= limit {
		fmt.Printf("Already reached %d. Increase -max to continue.\n", lastCenter)
		return
	}

	// 4. Update Prime Database (Incremental Sieve)
	if lastPrime < limit+1 {
		fmt.Printf("Updating primes from %d to %d...\n", lastPrime, limit+1)
		pf, _ := os.OpenFile(*primeFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		for n := lastPrime + 1; n <= limit+1; n++ {
			if isPrimeIncremental(n, primes) {
				primes = append(primes, n)
				primeMap[n] = true
				fmt.Fprintf(pf, "%d\n", n)
			}
		}
		pf.Close()
	}

	// 5. Main Conjecture Loop
	fmt.Printf("Processing centers from %d to %d...\n", lastCenter, limit)
	for n := lastCenter + 2; n <= limit; n += 2 {
		if primeMap[n-1] && primeMap[n+1] {
			var currentCombinations [][2]int

			// MIDPOINT OPTIMIZATION
			// We search for na starting from the midpoint (n/2) downwards to 0.
			// This leverages the "Midpoint Density" observation.
			mid := n / 2

			// If mid is odd, the nearest even center is mid-1 or mid+1.
			// Since na + nb = n, if na is mid-i, then nb is mid+i.
			for i := 0; i <= mid; i++ {
				na := mid - i
				// We only care about even centers (and the seed 2)
				if na == 1 {
					continue
				}

				nb := n - na

				if provenCenters[na] && provenCenters[nb] {
					// To keep results consistent with your previous version (na <= nb)
					if na <= nb {
						currentCombinations = append(currentCombinations, [2]int{na, nb})
					}

					// PERFORMANCE TWEAK:
					if *optimized {
						break
					}
				}
			}

			if len(currentCombinations) > 0 {
				// To keep the JSON file ordered as before
				sort.Slice(currentCombinations, func(i, j int) bool {
					return currentCombinations[i][0] < currentCombinations[j][0]
				})

				newResult := Result{
					Center:       n,
					PrimePair:    [2]int{n - 1, n + 1},
					Combinations: currentCombinations,
				}
				results = append(results, newResult)
				pool = append(pool, n)
				provenCenters[n] = true
				fmt.Printf("[SUCCESS] Center %d: %d combinations found (nearest to mid: %v)\n",
					n, len(currentCombinations), currentCombinations[len(currentCombinations)-1])
			} else {
				fmt.Printf("\n[CONJECTURE FAILED] Center %d has no combinations.\n", n)
				saveResultsCompact(*jsonFile, results)
				return
			}
		}
	}

	saveResultsCompact(*jsonFile, results)
	fmt.Printf("\nDone. Verified up to %d in %v.\n", limit, time.Since(start))
}

// saveResultsCompact writes JSON where each pair in combinations gets its own line
func saveResultsCompact(filename string, results []Result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Center < results[j].Center
	})

	var sb strings.Builder
	sb.WriteString("[\n")

	for i, r := range results {
		pairBytes, _ := json.Marshal(r.PrimePair)

		sb.WriteString("  {\n")
		sb.WriteString(fmt.Sprintf("    \"center\": %d,\n", r.Center))
		sb.WriteString(fmt.Sprintf("    \"primePair\": %s,\n", string(pairBytes)))
		sb.WriteString("    \"combinations\": [\n")

		for j, comb := range r.Combinations {
			combBytes, _ := json.Marshal(comb)
			sb.WriteString("      " + string(combBytes))
			if j < len(r.Combinations)-1 {
				sb.WriteString(",")
			}
			sb.WriteString("\n")
		}

		sb.WriteString("    ]\n")
		sb.WriteString("  }")

		if i < len(results)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]")

	os.WriteFile(filename, []byte(sb.String()), 0644)
}

func isPrimeIncremental(n int, existingPrimes []int) bool {
	if n < 2 {
		return false
	}
	for _, p := range existingPrimes {
		if p*p > n {
			break
		}
		if n%p == 0 {
			return false
		}
	}
	return true
}

func loadPrimes(filename string) ([]int, map[int]bool, int) {
	list := []int{2, 3}
	m := map[int]bool{2: true, 3: true}
	file, err := os.Open(filename)
	if err != nil {
		return list, m, 3
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val, _ := strconv.Atoi(scanner.Text())
		if val > 3 {
			list = append(list, val)
			m[val] = true
		}
	}
	return list, m, list[len(list)-1]
}

func loadResults(filename string) []Result {
	var results []Result
	file, err := os.ReadFile(filename)
	if err != nil {
		return []Result{{Center: 2, PrimePair: [2]int{1, 3}, Combinations: [][2]int{}}}
	}
	json.Unmarshal(file, &results)
	return results
}
