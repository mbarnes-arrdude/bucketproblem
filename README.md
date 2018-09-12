# bucketproblem
code challenge - temporary

#Demo Instructions

1. Install a golang development environment (Its easier than you think) http://golang.org
2. Download, clone or `go get -t https://github.com/mbarnes-arrdude/bucketproblem`
3. compile the binaries and put them in your path (if $GOBIN is not set or in your path)
4. run `calcbucket 5 3 4`
5. run `runbucket` for an interactive way to launch simulations with long run times

# The Problem
## The Good Die Hard
This project is a code challenge to solve an abstract problem made famous in the movie Die Hard 3. In the movie a pair of heros are run ragged across New York City by a mad bombsman who puts them through feats of strength, character, intelligence and agility. At one such challenge, they find themselves at a park, next to a fountain, with a 3 gallon bucket, a five gallon bucket, and a bomb ready to explode killing hundreds of innocent men, women and children if 4 gallons is not placed on a scale in 30 seconds.

Now, I want to note that the first time I saw the movie I thought the problem easy. Any container with parallel sides and mirror symmetry in its profile may be poured to exactly half its volume if tilted so that the level of the water touches the lip and the back edge. 5*0.5 +3*0.5 = 2.5 + 1.5 = 4! Just as the people in the theater did not appreciate my excalaimation of that solution while the heros fumbled to find the best method, the person who put the question to me to solve as a coding challenge did not either.

This project is a library with 2 demo implementations of the abstract bucket challenge solved for all real positive integers. The challenge is to predict the best sequence of operations to satisfy any problem involving any 2 sized buckets and any desired result. Some of the soft requirements I was told were scalability and a solution for arbitrary number size (beyond the domainof int64). Once the best route is chosen, the challenge requires an output of the table of operations.

### Priorities
The hard rules for the challenge are in order of precedence:
1. Functionality
2. Efficiency (Time, Space)
3. Code Quality / Design / Patterns
4. Testability
5. UI/UX design

### Considerations
*Functionality* OK, it has to work... and work right. The problem is complex to break down. The only reasonable routes through the process of filling one bucket dumping it into another and then emptying and/or filling are only 2 in form:

1. Fill the big one, pour into the small one til it fills, empty the small one and repeating til the big one is ready to fill again.
2. Alternately fill the small one and continue to fill and pour it into the large one until it is full.

Repeating either of these processes will get you to the answer. The problem breaks down to modular math. The ratio of the bucket sizes describes a bicyclic modular number series. This is the series of numbers that are in a range of Big Bucket A's volume times Big Bucket B's divided by their greatest common denominator.

The cycle is identical as a segment to the previous segment on the integer number line, and also has the property that the ratios may be inverted in cycle. That is to say that the cycle of remainders of the mod of either bucket to a position in the series is unique and that the order of those states has a reverse symmetry. Therefore the problem is solvable algorithmicly.

The answer is to use the extended version of Euclids Algorithm for Greatest Common Denominator. This is a very old algorithm for reducing pairs of large numbers into their GCDs. One takes the larger number and divides the smaller into it getting the mod. The smaller number is then multiplied back to the mod and the process repeats until the number cannot be reduced further. The result is the GCD of the two numbers. Re-integrating those steps will then result in an identity of the form Ay + Bx ≡ 1(modA). The distance to travel on the series of numbers in the domain is (x(modA) * B * desired)modA. The distance the other direction is the inverse of the result (modA).

The total number of steps is 2 for the initial fillup and pour + 2 for each "count" as additional fills and pours, and another +2 for each empty and pour at a rate of B/A additional 2 steps per pour. If the desired amount is larger than the smaller bucket, you will have one less pour and fill.
