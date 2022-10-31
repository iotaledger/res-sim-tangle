# res-sim-tangle

`res-sim-tangle` is a Go-based simulator that allows to study various aspects of the [Tangle](https://assets.ctfassets.net/r1dr6vzfxhev/2t4uxvsIqk0EUau6g2sw0g/45eae33637ca92f85dd9f4a3a218e1ec/iota1_4_3.pdf).

In the simulation blocks are created through a Poisson process (i.e. random) or with fixed rate. Each block references a number of previous blocks, thus creating the block DAG. To allow through rapid simulations  we use the following approach. 

For the creation of blocks no different nodes exists. Instead, each block is at creation assigned to an ID and added to the global pool of blocks. Each block is hidden for some time and cannot be selected by the tip selection during the time it is hidden. This process emulates, to a reasonable degree, a distributed network, where the delay in processing and propagation is approximately the time a block is hidden. 

To allow for higher sampling, the simulator allows to run the simulation in parallel on multiple processors. Each of the simulation runs on the different processors is independent from the others. 

## Run the simulation

First, you need to [install Go](https://golang.org/doc/install) if it is not already installed on your machine. It is recommended that you use the most recent version of Go.

## Build

If you prefer, you can build your executable (as well as cross compiling for other architectures) by using the `go build` tool:

```
go build -o sim
```
You can then run your simulation by running:

```
./sim
```

## Parameters

|   Parameter       |       Type    | Description    |
|-------------------|:-------------:|:--------------|    
|   K               |   int         | number of tips to approve |
|   H               |   int         | network/computation delay |
|   Lambda          |   float64     | tx rate of occurance      |
|   Alpha           |   float64     | level of randomness in a random walk |
|   TangleSize      |   int         | number of txs forming the Tangle|
|   Seed            |   int64       | seed value of the PRNG|
|   TSA             |   string      | tip Selection Algorithm [`"URTS"`, `"RW"`, `"HPS"`]|
|   tsa             |   TipSelector | internal interface [`URTS`, `URW`, `BRW`, `HPS`]|
|   ConstantRate    |   bool        | if true enables constant tx rate, i.e., _1/Lambda_|
|   DataPath        |   string      | path to store results (currently unused)|
|   minCut          |   int         | initial part of the Tangle not considered for results, i.e., 30H * Lambda|
|   maxCut          |   int         | final part of the Tangle not considered for results, i.e., 30H * Lambda|
|   nRun            |   int         | number of Tangles simulated |
|   VelocityEnabled |   bool        | if true enables statistics on velocity of particles |

## Notes on parallelization

In the `main.go` file, a global variable `nParallelSims` is defined as: 
```go
var nParallelSims = runtime.NumCPU()/2 - 1
```
This variable defines the number of simulations to run in parallel. 
The GO function `runtime.NumCPU()` returns the number of logical cores of your machine. 
Potentially one could use all her/his available logical cores, but I've noticed that you start loosing performance when concurrency kicks in. For example, having a 6-core architecture implies only 6 physical cores, which might become 12 logical cores depending on your machine. Logical cores only help the processor to multitasking, thus, when it comes to parallelization it is recommended to stick with a maximum value for `nParallelSims` set to the maximum number of your _physical_ cores. One additonal core can be subtructed to have it for your web surfing while waiting for your results drinking a :coffee:

## Visualization

This is an example of visualizing the result of the simulation

The Tangle
![Tangle](images/Tangle.png)

Ghost particle path and its cone
![Tangle](images/GhostCone.png)

Visiting probability in red gradients
![Tangle](images/visitingP.png)

Random walker path
![Tangle](images/RW.png)
