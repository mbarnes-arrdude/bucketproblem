# bucketproblem
code challenge - temporary

# Demo Instructions

1. Install a golang development environment from http://golang.org (it's easier than you think)
2. Download, clone or `go get -t https://github.com/mbarnes-arrdude/bucketproblem`
3. Compile the binaries and put them in your path (if $GOBIN is not set or in your path)
4. Run `calcbucket 5 3 4`
5. Run `runbucket` for an interactive way to launch simulations with long run times

Note: runcalc requires a terminal to run in and will not launch in most IDEs

# The Problem
## The Good Die Hard
This project is a code challenge to solve an abstract problem made famous in the movie Die Hard 3. In the movie, a pair of heroes are run ragged across New York City by a mad bombsman putting them through feats of strength, character, intelligence and agility. At one such challenge, they find themselves at a park next to a fountain, with a 3 gallon bucket, a five gallon bucket, and a bomb ready to explode, killing hundreds of innocent men, women and children if 4 gallons is not placed on a scale in 30 seconds.

Now, I want to note that the first time I saw the movie I thought the problem easy. Any container with parallel sides and mirror symmetry in its profile may be poured to exactly half its volume if tilted so that the level of the water touches the lip and the back edge. (5 * 0.5) + (3 * 0.5) = 2.5 + 1.5 = 4 after all. Just as the people in the theater did not appreciate my exclamation of that solution, the person who put the question to me did not either.

This project is a library with 2 demo implementations of the abstract bucket challenge solved for all real positive integers. The challenge is to predict the best sequence of operations to satisfy any problem involving any 2 sized buckets and any desired result. Some of the soft requirements I was told were scalability and a solution for arbitrary number size (beyond the domain of int64). Once the best route is chosen, the challenge requires an output of the table of operations.
