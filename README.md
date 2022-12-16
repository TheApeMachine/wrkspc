# wrkspc

A single binary that takes you from development into production, and beyond.

## Development Guidelines

Should you want to develop on this code, here are a couple of pointers.

### ford

It is a custom execution environment where objects can be connected,
and send messages to each other, even across distributed systems.

It is a hierarchy that looks something like below.

```yaml
Workspace:
  Workloads:
    - Workload:
        Assemblies:
          - Assembly:
              Abstracts:
                - Abstract
                - Abstract
                - ...
          - Assembly:
              Abstracts:
                - Abstract
                - Abstract
                - ...
          - ...
    - Workload:
        Assemblies:
          - Assembly:
              Abstracts:
                - Abstract
                - Abstract
                - ...
          - Assembly:
              Abstracts:
                - Abstract
                - Abstract
                - ...
          - ...
    - ...
```

#### drknow.Abstract

[package](../drknow)

The best way to think about it is that the `Abstract` type is a building
block that composes some basic behaviors together.

The implementation of the inner workings of a behavior is left open to
modification.

**Abstracts** can exchange messages with other **Abstracts**.

By grouping them into `Assembly` types, we end up with a pipeline of
message sending operations.

Messages can be *mutated* (though in an append-only fashion), before
they enter the next stage of the process.

By grouping Assemblies into `Workload` types, we obtain yet another level
of composition, which can be used to perform operations on the results
coming out of multiple Assemblies.

Finaly the `Workspace` type holds it all together and acts as the top
level controller of the entire process.
