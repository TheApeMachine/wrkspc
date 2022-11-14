# zaha

This package is responsible for turning configuration structures
into buildable architectures for services.

## architecture

The `Architecture` type is a structure that defines a set of nested
`io.ReadWriter` objects, which are the primitives to build a service.
