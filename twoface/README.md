# twoface

A set of high-level concurrency "primitives" for Go development.

This package captures the most common patterns we encounter in distributed
systems and provides a friendler interface for integration.

Below is a description of each type contained in the package.

## Pool

A type that pre-allocates a pool of goroutines that will be re-used, making
them more useful for smaller tasks (spinning up a goroutine could slow down
a small enough task).

Secondary benefits include keeping the amount of goroutines in check, with
the ability to manually (or programmatically) scale the pool size, reduced
activity on the Garbage Collector, and memory efficiency.

