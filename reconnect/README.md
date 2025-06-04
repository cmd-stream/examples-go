# reconnect
This example demonstrates how clients can automatically restore their connection 
to the server.

## Details
The connection can be lost in two ways:
- While sending a Command, causing `Client.Send()` to return an error.
- While waiting for a Result, leading to uncertainty about whether the Command 
  was executed on the server.

In both cases, using the reconnect client allows the Command to be resent 
(assuming it is idempotent) after a short delay. If the connection was 
re-established, normal operations resume, otherwise, Client.Send() will return 
an error again.

To create a client group with reconnect clients:
```go
group, err := cmdstream.MakeClientGroup(1, codec, connFactory,
  grp.WithReconnect[T](),
  ...
)
```