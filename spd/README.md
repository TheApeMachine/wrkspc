# spd

## A Dynamic API Model

Secure Private Datagram play an essential role in the concept of a Dynamic API Model
that is defined by what you write to it.

### Endpoints as Data

The *context*, or *header* of a datagram determines its canonical key.

```golang
type Datagram struct {
    Version   string
    Type      string
    Role      string
    Scope     string
    Identity  string
    Timestamp int64
}
```

The above structure with values would translate to the (example) prefix below.

```
v4.0.0/binary/datapoint/user/app.acme-corp.com/1663622730/<UUID>
```

The identity field should be governed to always correctly reflect the origin.

In this way, as long as the client has access to the endpoint, without any additional
access control layers, the client can use the role and scope fields to create new
"endpoints" under which to store data.

```
v4.0.0/binary/datapoint/profile/app.acme-corp.com/1663622730/<UUID>
```

The above example was deliberately chosen, because it will often be related to the
data under the user prefix we described earlier.

### Automatic Relational Mapping

Using a fast lookup structure by way of an in-memory radix tree it becomes possible
for us to make many alternative prefixes that all lead to the same data.

More importantly, by just permutating the prefix itself we start grouping relational
data automatically.

A simplified example below, using the *user* and *profile* prefixes.

```
v4.0.0/binary/datapoint/app.acme-corp.com/user/1663622730/<UUID>
v4.0.0/binary/datapoint/app.acme-corp.com/profile/1663622730/<UUID>
```

By simply switching the *scope* and *identity* we have created a new prefix under
`v4.0.0/binary/datapoint/app.acme-corp.com` which will return the user and the
profile, and what we have left to do for the projection is merge these two records
into the desired aggreggate structure.
