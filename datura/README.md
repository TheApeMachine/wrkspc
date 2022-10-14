# datura

An abstraction around prefix based key/value storage built on top of the [store interface](https://github.com/TheApeMachine/wrkspc/blob/master/passepartout/store.go?ts=4).

Optimally designed for S3 compatible storage (AWS S3, GCP Cloud Storage, MinIO, DataFlare R2, etc.).

Cached based on [Immutable Radix Tree](https://github.com/hashicorp/go-immutable-radix).

Uses auto scaling goroutine pools.

## S3 Compatible Storage

Connection and data wrangling methods for any S3 compatible storage solution, including but not limited 
to AWS S3, GCP Cloud Storage, MinIO, CloudFlare R2.

Represents the persistent data layer of the *wrkspc* environment.

## Radix Tree

Fastest prefix based caching mechanism with `O(k)` lookups.

The implementation is **1:1** compatible (and interchangeable) with the S3 logic.

## Radix Forest

Experimental HA orchestration for radix trees.
