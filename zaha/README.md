# zaha

This package allows you to use the interfaces defined for connections, stores, and jobs to (dynamically) build a network service.

The idea is that one can run a service by using `wrkspc serve <SERVICE_NAME>` and it will run the corresponding architecture.

## architecture

This package links together a set of connections, stores, and jobs to provide a (micro) service that is exposed over a network.

This network can be both internal, as well as external.
