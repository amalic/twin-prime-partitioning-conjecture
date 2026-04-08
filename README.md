# The Malic Conjecture on Twin Prime Additive Partitioning

**Proposed by Alexander Malic (2026)**

## 1. Definitions
Let $\mathbb{P}$ be the set of all prime numbers.
Let $T$ be the set of **Twin Prime Centers**, defined as the set of even integers whose neighbors are both prime:

$$T = \{ n \in \mathbb{N} \mid (n-1) \in \mathbb{P} \text{ and } (n+1) \in \mathbb{P} \}$$

*Note: For the purposes of this genealogy,* $n=2$ *is defined as the initial seed, corresponding to the pair* $(1, 3)$*.*

## 2. The Conjecture
For every twin prime center $n \in T$ where $n > 2$, there exists at least one pair of twin prime centers $(n_a, n_b)$ such that $n$ is their sum:

$$\forall n \in T, n > 2 \implies \exists n_a, n_b \in T : n = n_a + n_b$$

where $n_a, n_b < n$.

*Note on Terminology: While the set* $T$ *is not "additively closed" in the strict sense (e.g.,* $4+4=8$*, and* $8 \notin T$*), the conjecture posits that every element of* $T$ *is partitionable into two other elements from the same set.*

## 3. Relationship to Goldbach and Twin Prime Conjectures
The **Malic Conjecture on Twin Prime Additive Partitioning** represents a **restricted refinement** of the Strong Goldbach Conjecture. While Goldbach asserts that every even integer $2k > 2$ is the sum of two primes, this framework restricts the summands specifically to the centers of twin prime pairs ($T$).

It is important to note that this conjecture implies the **Twin Prime Conjecture**; for the partitioning of increasingly large centers to remain possible, the set $T$ must be infinite. It suggests that the additive properties of primes are robust enough to persist even when restricted to this significantly sparser subset.

## 4. Density and Brun's Constant
The validity of this conjecture is deeply linked to the density of the set $T$. It is well known from the work of Viggo Brun (1919) that the sum of the reciprocals of twin primes converges to Brun's Constant $B_2$:

$$B_2 = \sum_{p, p+2 \in \mathbb{P}} \left( \frac{1}{p} + \frac{1}{p+2} \right) \approx 1.90216$$

The convergence of this series implies that $T$ is a "thin" set compared to the set of all primes. The Malic Conjecture suggests that despite this relative sparseness, the distribution of $T$ is sufficiently dense—as predicted by the **First Hardy-Littlewood Conjecture** (specifically the $k=2$ case)—to maintain this additive partitioning property for all $n > 2$.

## 5. Computational Verification
Extensive computational analysis has been performed to verify this additive property. Using an incremental sieve and an additive genealogy tracking algorithm in Go, the conjecture has been verified for all $n \in T$ up to:

* **Exhaustive Genealogy ($n = 1,000,000$):** Full additive combinations stored in `conjecture.json`.
* **Optimized Midpoint Search ($n = 100,000,000$):** High-speed verification sample stored in `conjecture_optimized.json`.
* **Verified Upper Bound ($n = 1,000,000,000$):** Successfully verified to $10^9$ via local execution (see Data Availability).
* **Observations:** 
    * The number of valid pairs $(n_a, n_b)$ for a given $n$ generally scales with the magnitude of $n$, suggesting that "isolated" centers become increasingly improbable as the numerical space expands.
    * **Midpoint Density:** Empirical data suggests that for most $n$, there exists at least one pair where $n_a \approx n_b \approx n/2$. This indicates that the additive richness of twin prime centers remains concentrated around the arithmetic mean, similar to the behavior observed in the standard Goldbach Conjecture.

### Data Availability and Performance
This repository provides two primary datasets for research and verification:
1. **`conjecture.json`**: The full additive genealogy for all centers up to $1,000,000$.
2. **`conjecture_optimized.json`**: A midpoint-optimized verification sample up to $100,000,000$.

**Regarding the $10^9$ (1 Billion) Verification:**
The verification for $n = 1,000,000,000$ was completed successfully. However, due to GitHub's 100 MB file size limit, the resulting JSON file (exceeding 500 MB) is not hosted directly in this repository. 

* **Performance Benchmark:** On a MacBook Air (M4, 16GB RAM), the optimized verification up to $10^9$ completes in approximately **10 minutes**.
* **Reproducibility:** The $10^9$ result can be independently reproduced by running the provided Go source code with the `-optimized` flag.

### Running the Verification
To reproduce the results or extend the verification, use the following command:

```bash
go run main.go -max=1_000_000
```

Optimized midpoint search up to 100 Million

```bash
go run main.go -max=100_000_000 -optimized=true -db=conjecture_optimized.json
```

Verify to 1 Billion (Approx. 10 mins on M4 hardware)

```bash
go run main.go -max=1_000_000_000 -optimized=true -db=my_local_verification.json -primes=my_local_primes.txt
```

The source code, methodology, and resulting datasets (including the JSON-formatted genealogy) are available for review in this repository.

---

**License:** This work and the associated dataset are licensed under a [Creative Commons Attribution 4.0 International License (CC BY 4.0)](https://creativecommons.org/licenses/by/4.0/deed.en).

**Software License:** The source code in this repository is licensed under the [MIT License](LICENSE).

**Author:** Alexander Malic
