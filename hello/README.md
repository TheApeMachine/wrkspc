# Hello (Operator?)

This is a work in progress concept to make any type be able to talk to any other type.

It is an ongoing effort, combined with `Hefner Pipes` to make concurrency workflow self-similar
on and off the "machine".

In other words, it starts with trying to make Go channels and streaming network sockets look, feel,
and act the same.

The Operator type in the hello package is likely going to be responsible for just making the
connections between types, and assigning some "translation" method between their individual messages.
