# bucketproblem
code challenge - temporary

# Demo Instructions
## Untested at time of writing

1. Install a golang development environment (Its easier than you think) http://golang.org
2. Download, clone or `go get -t https://github.com/mbarnes-arrdude/bucketproblem`
3. Compile the binaries and put them in your path (if $GOBIN is not set or in your path)
4. Run `calcbucket 5 3 4`
5. Run `runbucket` for an interactive way to launch simulations with long run times

Note: runcalc requires a terminal to run in and will not launch in most IDEs

# The Problem
## The Good Die Hard
This project is a code challenge to solve an abstract problem made famous in the movie Die Hard 3. In the movie a pair of heros are run ragged across New York City by a mad bombsman who puts them through feats of strength, character, intelligence and agility. At one such challenge, they find themselves at a park, next to a fountain, with a 3 gallon bucket, a five gallon bucket, and a bomb ready to explode killing hundreds of innocent men, women and children if 4 gallons is not placed on a scale in 30 seconds.

Now, I want to note that the first time I saw the movie I thought the problem easy. Any container with parallel sides and mirror symmetry in its profile may be poured to exactly half its volume if tilted so that the level of the water touches the lip and the back edge. (5 * 0.5) + (3 * 0.5) = 2.5 + 1.5 = 4 after all. Just as the people in the theater did not appreciate my excalaimation of that solution, the person who put the question to me to solve as a coding challenge did not either.

This project is a library with 2 demo implementations of the abstract bucket challenge solved for all real positive integers. The challenge is to predict the best sequence of operations to satisfy any problem involving any 2 sized buckets and any desired result. Some of the soft requirements I was told were scalability and a solution for arbitrary number size (beyond the domainof int64). Once the best route is chosen, the challenge requires an output of the table of operations.

### Priorities
The hard rules for the challenge are in order of precedence:
1. Functionality
2. Efficiency (Time, Space)
3. Code Quality / Design / Patterns
4. Testability
5. UI/UX design

### Considerations
**Functionality** OK, it has to work... and work right. The problem is complex to break down. The only reasonable routes through the process of filling one bucket dumping it into another and then emptying and/or filling are limited:

1. Fill the big one, pour into the small one til it fills, empty the small one and repeating til the big one is ready to fill again.
2. Conversely, fill the small one and continue to fill and pour it into the large one until it is full.

Repeating either of these processes will get you to the answer. The problem breaks down to modular math. The ratio of the bucket sizes describes a bicyclic modular number series. This is the series of numbers that are in a range of Big Bucket A's volume times Big Bucket B's divided by their greatest common denominator.

Any other move squence would involve going backwards to a previous move.

The cycle is identical as a segment to the previous segment on the integer number line, and also has the property that the ratios may be inverted in cycle. That is to say that the cycle of remainders of the mod of either bucket to a position in the series is unique and that the order of those states has a reverse symmetry. Therefore the problem is solvable algorithmicly.

The answer is to use the extended version of Euclids Algorithm for Greatest Common Denominator. This is a very old algorithm for reducing pairs of large numbers into their GCDs. One takes the larger number and divides the smaller into it getting the mod. The smaller number is then multiplied back to the mod and the process repeats until the number cannot be reduced further. The result is the GCD of the two numbers. Re-integrating those steps will then result in an identity of the form Ay + Bx ≡ 1(modA). The distance to travel on the series of numbers in the domain is (x(modA) * B * desired)modA. The distance the other direction is the inverse of the result (modA).

The total number of steps is 2 for the initial fillup and pour + 2 for each "count" as additional fills and pours, and another +2 for each empty and pour at a rate of B/A additional 2 steps per pour. If the desired amount is larger than the smaller bucket, you will have one less pour and fill.

**Efficiency (Time, Space)** I have an algorithm that is not quite 1n. Euclid's algorithm is not linear but thanks to Knuth, most big int libraries have very fast implementations. As long as I don't try and store the results table space is not an issue. Then there is the run-time of the simulation. If the solution is a big integer, it could take some time to run through the series. Space is a non-issue except for the size of buffers asynchronously feeding the UI and the size of the solution table. The default buffer size for the simulation is set high because it often greatly overruns ANYTHING that writes to a screen buffer. It would not take much to make the buffer sizes configurable or even adjustable, there would just be longer "miss" periods in the client.

If time is not an issue the design is lends itself to 

**Code Quality / Design / Patterns** There are philsophies, other philsophies, the competing philsophies, then there is me... I will generally learn the expectations of a shop and match those. Kind peer review and suggestions are always welcomed so that I can quickly learn expectations and preferences.

**Testibility** Totally important and while NYI some great consideration to speccing test has been done and expressed in the files that will contain them. Unit test coverage is expressed in rspec-ish comments in corresponding ztest{gofile}.go to each source file in the library. Naive returns are not tested. Tests are not written but a pass at coverage in unit tests is done.

Much of the testing of a multi bracketted algorithm is acceptance tests. There are a few comments that describe bracketing inputs and comparing past output. This is necessary for proving the algorithms in the big.Int number space. Small bugs could make a naive set of tests pass because the result is similar. For instance it could be very easy to have the algorithm change so that there is a miss-calculated length of the series but everything else looks ok and no one knows until a simulation is run with the right inputs. To the clients everything would be OK, but we could experience cache misses or buffer underruns with careless code elsewhere falling prey to what other processes ignore.

**UI/UX Design** I am not interviewing to be a front end developer though I have some skillz in design. I did want to show off the speed of the simulation with the chosen language and the ability to start/stop and pause multiple simulations managed in the same process.  Both CLI based demos use a thrid party framework, the first is to make injestion of arguements easier to extend. I can imagine a lot of features to the command. Having a framework to handle flag and arguments as well as configs is nice to have early and makes refactoring easier.

The `runbucket` UI is a little more fancy and uses a terminal's built in color and drawing capabilities. It will require a full terminal emulator to run. It was a good choice to implement this one. Running long simulations to veryify the algorithm's predictions was hard with just a single process executable and invisible to me unless I wanted a spew of a couple tens of millions of lines thrown at me. I wanted to be able to run multiple problems at once and monitor them together without a synchronous screen output bog down my speed.

As a result the library implements a backgrounded process for solution and simulation and reads double-buffer style the event stream from those processes through Go channels. Implementations using the library can subscribe to the current controllers channels and monitor the simulation without slowing it down much except the shared CPU. If you put big problems into the simulator, expect your CPU fan to turn on... I wrote it to be greedy. Changing the niceness is only a matter of tuning buffer sizes and/or increasing sleep time between steps. Also implementing the controller as an interface and extending (or limiting) it's behaviors is fairly trivial, far more naive and synchronous implementation would be very fast to write and could be plugged into the rest of the code. Inversion of control pivots on this functional unit.

**Soft: Scalability**
An algorithm for finding any one position is posible, but streaming an answer is a more interesting problem if one pretends there could be a use for running a simulation of this type for a large population of users. The soft requirements talk about scalability... it isn't much of a question except in a scenario like this. So in projecting the problem - I imagine a cloud service that clients connect to, provide problems and get back solutions.

Oof simulations of the problem with really big numbers can run for a LONG time. Small ones are fairly fast especially with a (yet to be implmented) int64 based solution. Even faster would be some set of small number solutions stored in a memory cache. What is this for? I don't know but lets assume I have the need to simulate bicyclic patterns expressed with really large numbers. Perhaps I am mimicing the perfect aggregated signal from two rotational processes for a comparitor looking for irregularities?

From this imagineering I project the possible need for the ability to monitor and administrate the simulation for a lot of really big problems. They would need to stream their simulation states into kinesis (SQS might be fast enough but more code) and use firehose to spew them into giant files on S3 from which simulations can be played back or examined for features. I hope this is not a requirement! It would get expensive. From that feed we might also serve up a websocket that mirrors the stream for the client to receive live data from the simulation.

To better scale the problem in needing lots of smaller requests handled quickly I could also see some swarm of faster int64 solvers feeding an LRU on DynoDB for inputs of int32 number space and less. SQS might suffice as transport.

Finally for super speedy delivery of smaller scale problems, it might be worth it to serve pre-compiled solutions from a memory cache populated with a 3 dimensional bracket of all integers 0 < X <= 0xff (yes 0... we want precompiled errors too). I might even want to be able to connect to a stream of state objects coming from the simulation running on a server synchoronously or asynchronously.

**Soft: Arbitrary Number Size**
A big int implementation is not available on the web and as a soft requirement it makes sense to develop this first. It is not as fast as a int64 implemenation but it handles a LOT more inputs. Any implementation of int64 will be a derivative of the same algorithm in Big Int but easier. Big Int problems can be made to run a REALLY long time for testing for memory leaks. The form of the big.Int library is functional and develops an algorithm more easilty translated into scala for in memory processors like Ignite (for REALLY REALLY REALLY big numbers) or most libraries that leverage SSE instructions or other custom math hardware like a crypto-card or other ASIC or FPGA designed for quick modular math.
