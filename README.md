#Introduction#

*corridor* is library containing Go language functions for the implementation of a concurrent genetic algorithm for the multi-objective corridor location problem. This problem involves finding the least cost connected pathway through a discrete search domain in which each location is characterized by one or more measures of cost. The library requires that the user provide a predefined search domain, objective function(s), and input parameters specifying the nature of the problem (i.e. desired start and destination locations).

##Installation##

The project is hosted as a publicly available GitHub repository. Providing that your local client GOPATH and GOROOT variables have been previously defined, the repository can be cloned and built using the following single shell command:

$ go get github.com/ericdfournier/corridor

##Description##

The work contained in this library is based upon the MOGADOR algorithm that was first introduced by Zhang & Armstrong (2008) in: http://www.envplan.com/abstract.cgi?id=b32167 . It also contains additional modifications to the initialization routine introduced by Fournier (2014) in: placeholderURL .

##Input Format##

All inputs must be formatted as comma delimited value (CSV) files.

###Example Search Domain###

The search domain should be encoded in a binary format with cells in the feasible search domain set to a value of 1 and cells outside of the feasible search domain set to a value of 0 as below. The user need not generate a "buffer zone" of zero encoded cells surrounding the feasible search domain as this is done automatically by the algorithm at runtime.

````
$ cat searchDomain.csv

0, 0, 0, 0, 0,
0, 1, 1, 1, 0,
0, 1, 1, 1, 0,
0, 1, 1, 1, 0,
0, 0, 0, 0, 0;
````

###Example Search Objectives###

The user should note that the objective values for cells that are outside of the search domain will be automatically set to be equal to an arbitrarily high value. Specifically, the objective scores for the locations which are outside of the feasible search domain values are set to be equal to the total number of cells (feasible and otherwise) contained within the entire search domain. For an example illustration of how this work, please see below.

````
$ cat objective1.csv

25, 25, 25, 25, 25,
25,  2,  3,  3, 25,
25,  1,  2,  5, 25,
25,  1,  1,  4, 25,
25, 25, 25, 25, 25;
````

````
$ cat objective2.csv

25, 25, 25, 25, 25,
25,  4,  1,  3, 25,
25,  5,  3,  6, 25,
25,  2,  1,  3, 25,
25, 25, 25, 25, 25;
````

###Example Source and Destination Subscripts###

The source and destination subscript files should be formatted to contain, separately, the row and column subscripts corresponding to the location of the either the source or the destination within the context of the input search domain grid. These subscripts should be stored as two comma separated values on a single line of the input .csv file as below.

````
$ cat sourceSubs.csv

1,1
````

````
$ cat destinationSubs.csv

3,3
````
##Output Format##

If the Algorithm fails to converge upon a solution within the given iteration limit, an error message will be printed to the console and a basic log.cv file will be written to the local directory. This log file contains information about the computational runtime and the total number of evolutionary iterations that were executed (which in this case will be equal to the maximum number of evolutions specified by the user).

If the Algorithm successfully converges upon a solution within the given iteration limit, a success message will be printed to the console and two files will be written to the local directory. The first is a log.csv file which contains information about the same information quoted previously. The second is an output solution file which contains the row and column subscripts for each step along the solution corridor. Additionally, subsequent rows within this output file will contain the individual, stepwise, objective scores for each of the objectives, for each step along the solution corridor.

###Example Output###

A possible output solution file for the previously constructed example problem is illustrated below.

````
[Line #1: Solution #1 Corridor Row Subscripts] 1, 1, 1, 2, 3;
[Line #2: Solution #1 Corridor Column Subscripts] 1, 2, 3, 3, 3;
[Line #3: Solution #1 Objective 1 Scores] 2, 3, 3, 5, 4;
[Line #4: Solution #1 Objective 2 Scores] 4, 1, 3, 6, 3;
[Line #5: Solution #2 Corridor Row Subscripts] ...
[Line #6: Solution #2 Corridor Column Subscripts] ...
[Line #7: Solution #2 Objective 1 Scores] ...
[Line #8: Solution #2 Objective 2 Scores] ...
...
````

This pattern is repeated for each output solution requested from the final population by the user. Solutions are automatically sorted by fitness score such that the first solution is the best, the second is the second best, etc.

#Benchmarking#

Two benchmark suites have been developed for this package. The first in a single run benchmark. Which evaluates the performance of the algorithm for a contrived problem specification on a particular machine given a single set of evolutionary runtime parameters. This "Single" suite is usefull for getting a feel for the scaling relationships between population size, runtime, and solution quality. 

The second benchmark suite is a monte carlo based simulation which takes are particular population size setting and uses repeated solution runs to deliver an estimate of the expected variation in average solution qaulity between runs due to the stochastic nature of the evolutionary optimization process. Sample usage of both benchmark suites are provided below.

##Single Benchmark Examples##

A "Single" sample test suite has been built into the package which allows the user to benchmark both the runtime performance of the algorithm as well as the output solution quality under various parameter settings. This is done in the following set of examples using a contrived problem specification in which there is a single known, globally optimal solution, set amidst a decision space containing randomly distributed costs. This globally optimal solution, relative to the start (S) and destination (D) points, is plotted below:

Globally Optimal Solution, F = [0.0, 0.0, 0.0]:

````
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 S 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 0 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 0 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 D 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
````

###Small Population Size###

Running the benchmark with a Small population size [P = 1,000], as in the following command, will deliver the following output. Note the characteristics of the population at convergence reveals that a near optimal solution was found. 

````
$ go test -bench=.SingleSmall
````

Final Population Distribution at Convergence [Small Population]:

````
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0 1000    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0 1000    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0 1000    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0 1000    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0    0 1000    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0    0    0 1000    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0    0    0 1000    0    0    0    0    0    0    0    0    0    0    0    4    3    0    0    0    0    0]
[   0    0    0    0    0    0 1000 1000  998 1000 1000 1000 1000 1000 1000 1000 1000  996  997 1000    2    0    0    0]
[   0    0    0    0    0    0    0    0    2    0    0    0    0    0    0    0    0    0    1    0  998    2    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0  998    2    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0  998    2    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0 1000    0    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0 1000    2    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0  998    2    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0 1000    0    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0 1000    0    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0 1000    0    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0 1000    0    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
[   0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0    0]
````

For this particular Small sized benchmark run, the mean fitness values for all of the individual solutions in the final output population were:

````
F = [8.235, 39.484, 83.719]
````

56 Evolutions were required to achieve convergence and the elapsed runtime, on this particular machine, was 3.84 seconds.

### Medium Population Size###

Running the benchmark with a Medium population size [P = 10,000], as in the following command, will deliver something like the following output. Here, the quality of the output solution has improved, and in some cases may deliver the globally optimal solution.

````
$ go test -bench=.SingleMedium
````

Final Population Distribution at Convergence [Medium Population]:

````
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0 10000     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0 10000     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0  9997     3     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0 10000     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0  9996     4     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0  9996     4     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0     0 10000     1     1     0     0     0     0     0     0     1     0     0     0     1     1     0     0     0     0     0]
[    0     0     0     0     0  9999  9998  9996 10000 10000 10000 10000 10000 10089 10000 10000 10000  9999  9999 10905  6410     0     0     0]
[    0     0     0     0     0     0     1     4     0     0     0     0     0     0     0     1     1     0     0     4  9997     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0 10000     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0 10000     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     4  9996     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     4  9995     1     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0 10000     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0 10000     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0  9999     1     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0 10000     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0 10000     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
[    0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0     0]
````

For this particular Medium sized benchmark run the mean fitness values for all of the individual solutions in the final output population were:

````
F = [4.0409, 6.09, 23.1304]
````

48 Evolutions were required to achieve convergence and the elapsed runtime, on this particular machine, was 34.305 seconds.

###Large Population Size###

Finally, running the benchmark with a Large population size [P = 100,000], as in the following command, will deliver something like the following output. Here, the quality of the output solution has improved to the point in which it will nearly guarantee the delivery of the globally optimal solution for this problems specification.

````
$ go test -bench=.SingleLarge
````

Final Population Distribution at Convergence [Large Population]:

````
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      0 100000      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      0 100000      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      5  99994      1      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      2  99992      6      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0     11  99985      4      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      0  99991      9      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      0  99989     17      2      2      2      2      4      4      3      1      1      2      3      2      5      1      3      0      0      0      0]
[     0      0      0      2  99993  99998  99996 100000 100000 100000 100000 100000 100000 100000 100000 100000 100000  99992  99996 100000  73873      0      0      0]
[     0      0      0      0      1      0      2      3     11      1      2      1      3      2      3      6      4      3      3     22 100000      0      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      1  99998      1      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      7  99988      5      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      5 100000      1      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      3 100000     13      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      1 100000      0      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      3 100000      1      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0 100000     12      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0 100000      0      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0 100000      0      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
[     0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0      0]
````

For this particular Large sized benchmark run the mean fitness values for all of the individual solutions in the final output population were:

````
F = [0.01451, 0.02635, 0.04418]
````

58 Evolutions were required to achieve convergence and the elapsed runtime, on this particular machine, was 375.358 seconds.

###Comments###

With these three single benchmark run examples, note the roughly linear scaling of runtime with increased population size relative to the roughly linear improvement in mean output solution quality with population size. This is an explicit tradeoff of the genetic algorithm approach: larger populations take longer to achieve convergence yet are more likely to deliver globally optimal solutions.

##Monte Carlo Benchmark Example##

A "MonteCarlo" sample test suite has also been built into the corridor package for evaluating the expected distribution of runtimes between multiple separate solution runs for the same problem specification. The Monte Carlo benchmark suite repeats the single benchmark solution process described below for a fixed number of [N = 100] simulation runs. In the process, it compiles statistics about the mean and standard deviation of average fitness values for the set of final solution populations generated during the 100 simulation runs. At conclusions, it prints these descriptive statistics to the terminal.

###Small Population Size###

Monte Carlo simulations take a long time to execute. As a result, only a small population based example is presented below. However, the benchmark code can be modifed to run this type of simulation with any desired number of samples and with a population of any size. 

````
$ go test -bench=.MonteCarloSmall

...
Sample Size (N) = 100
Mean Population Aggregate Fitnesses | Minimum Fitness = 375.8316299999999 | 0.0
Standard Deviation of Aggregate Fitnesses = 145.1107156924006
Mean Runtime in Seconds = 4.0307659093499995
Standard Deviationof Runtimes in Seconds = 0.736894917206145
````

###Comments###

The Monte Carlo simulations benchmark suite is useful in helping to determine the distribution of the quality of output solutions that can be expected from certain combinations of problem specifications and evolutionary parameter settings. Due to its long runtime however, it is not advocated for use in most basic package testing scenarios.

#Author#

This project was developed by Eric Daniel Fournier [@ericdfournier] as part of his doctoral dissertation research at the University of California, Santa Barbara's Donald Bren School of Environmental Science & Management. The author would like to acknowledge the generous financial support of the Walton Family Foundation's Sustainable Water Markets Fellowship Program in making this development effort possible.

#Contact and Support#

If you have any questions about the usage of this library or would like to discuss the details of its implementation please email me@ericdfournier.com Please submit bug reports and other feature requests as issues via the GitHub repo. Thank you for your interest in this project!

