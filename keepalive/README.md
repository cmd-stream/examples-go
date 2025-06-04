# keepalive
keepalive example demonstrates how the client can keep a connection alive
with the server even when there are no commands to send. This is done by
configuring the clients with the `WithKeepalive` option:
```go
group, err := cmdstream.MakeClientGroup(clientsCount, codec, connFactory,
  grp.WithClientOps[T](
    ... 
    cln.WithKeepalive(
      dcln.WithKeepaliveTime(KeepaliveTime),
      dcln.WithKeepaliveIntvl(KeepaliveIntvl),
    ),
  ),
)
```
