# travel audience Go challenge

## Build / run
```bash
go get github.com/dmitry-yakutkin/tachallenge/server
cd $GOPATH/src/github.com/dmitry-yakutkin/tachallenge/server
go build
./server
```

## Rationale
### Set implementation
The main (and maybe the most controversial) thing this solution is based on is a custom sync.Map-based set implementation.
It serves 2 purposes:
* Safe container storage for concurrent environment.
* No more we need to worry about elements duplication. It's not a really big issue, but we still solve it this way.
This solution might have some drawbacks in terms of specifics of it's usage area - sync.Map is mostly used in Go compiler itself internally for resolving cache contention problems.
Regular Go's map + RWMutex could be used, as well as a custom tree-based structure to even get items sorted on insert, but those would need a bit more additional work and time.
### Ambiguity
Initial task description is a bit ambiguous in terms of how it defines the desired behaviour in cases when we exceeded handler timeout.
Current implementation suposses that we're trying to get as much as we can within specified timeframe, if it seems like we're going to exceed it - return just what we have now.
That's why actual http client requests timeout is a bit lower that one specified in the task description - for us to have some flexibility in terms of decision making while processing.
### Processing pipeline
Links processing pipeline is straightforward: spawn goroutines for each link processing, also spawn timer for timeout. Exit either if all links processing is finished or timeout has expired.
The fact that currently we just spawn as many goroutines as we get links might lead to some problems, but as it wasn't really specified in the task descriptions and Go's goroutines can actually be spawned in big quantities, I decided not to work additional on introducing pool or something else for this purpose.