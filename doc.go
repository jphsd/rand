/*
Package rand contains additional math/rand/Source64 implementations for use with math/rand's Rand type.

The standard math/rand Source64 uses an Additive Lagged Fibonacci Generator (ALFG) orignally designed by
DP Mitchell and JA Reeds for Bell Lab's Plan9 operating system.

SplitMix64 (period 2^64)
  A fixed-increment version of Java 8's SplittableRandom generator
  See http://dx.doi.org/10.1145/2714064.2660195 and
  http://docs.oracle.com/javase/8/docs/api/java/util/SplittableRandom.html
  It is a very fast generator passing BigCrush, and it can be useful if
  for some reason you absolutely want 64 bits of state; otherwise, we
  rather suggest to use a xoroshiro128+ (for moderately parallel
  computations) or xorshift1024* (for massively parallel computations)
  generator.
  - Sebastiano Vigna

XOshiro256 (xoshiro256** in the literature) (period 2^256 - 1)
  According to Vigna, this is faster and produces better output than the xorshift family of generators.
  (See https://en.wikipedia.org/wiki/Xorshift and https://prng.di.unimi.it/)

PCG (Permuted Congruential Generator) (period 2^128)
  An implementation of the PCG XSL RR 128/64 (LCG) generator described in Melissa O'Neill's paper.
  (See http://www.pcg-random.org/pdf/toms-oneill-pcg-family-v1.02.pdf)
  There's an alternative implementation in golang.org/x/exp/rand.

Relative speeds (not very scientific)
  Without/with locking
  SplitMix64 0.345/1.540
  XOshiro256 0.379/1.717
  ALFG       0.429/1.748
  PCG        0.454/1.728

The statistical testing suite Big Crush/TU01 is described at http://simul.iro.umontreal.ca/testu01/tu01.html
and available at https://github.com/umontreal-simul/TestU01-2009/
*/
package rand
