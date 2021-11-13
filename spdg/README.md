# Secure Private Datagrams

**EXPLAIN THAT NAME:** A (cloud) Datagram is one "unit" of cloud native data. Eventually privacy and
security will be implemented as first-class citizens of the type itself.

A Cloud Native type that aims to be used (almost) everywhere to transport data.

It wraps any data (as bytes) with a Context header that holds meta-data about the inner type.

This has a few benefits.

1. You can marshal the whole thing to bytes for storage (in S3 compatible storage for instance) and
   still predict which type it is.

2. Once unmarshaled back into Datagram format the Context header gives you the inner data type.

3. It can transport anything from a string to a pointer, to a network connection itself. If it gob
   encodes (which everything does) it can be transported. This is an advantage over gRPC.

4. Using this for almost all inputs and outputs for methods and functions creates a "mono typed"
   system, which has proven to be really flexible and easy to work with.

5. Datagrams are nestable, one can be stored like any other type into another Datagram, so in that
   way it has a native sense of "batching".

The Datagram Context header also provides the building block (its fields) to generate a "prefix"
based slug as seen in many cloud storage solutions.

This has shown to be a very workable feature when it comes to both storage and well as retrieval.

Allowing for an easy "single endpoint" design, it means there is not need to define API endpoints
in the traditional sense, is compatible with most storage solutions and ultimately designed to use
extramely fast big data structures such as Radix Trees.
