# twoface

A set of high-level concurrency "primitives" for Go development.

This package captures the most common patterns we encounter in 
distributed systems and provides a friendler interface for integration.

Below is a description of each type contained in the package.

## Context

A type that wraps the native context.Context and improves developer
ergonomics.

It is most useful as a "generic" key/value store to pass data around,
be it internally or over the network.

## Job

An interface type that can be implemented if a certain workload can
be done concurrently, using a worker pool of pre-allocated goroutines.

## Pool

A type that pre-allocates a pool of goroutines that will be re-used, 
making them more useful for smaller tasks (spinning up a goroutine 
could slow down a small enough task).

Secondary benefits include keeping the amount of goroutines in check, 
with the ability to manually (or programmatically) scale the pool size, reduced
activity on the Garbage Collector, and memory efficiency.

## Scaler

A type that can automatically scale a worker pool of pre-allocated goroutines
from zero to however many the system can handle, and back down.

It allows for resources to be rebalanced based on the load on particular
parts of the program, versus just using statically sized worker pools.

## Signal

Wraps around the system interrupt signals to facilitate graceful shutdown.

## Worker

A wrapper around a goroutine that takes in a Job channel so it can take on
Job instances and run them as part of a Pool.
