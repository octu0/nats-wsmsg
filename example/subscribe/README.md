## websocket subscription

# run subscription

By connecting to specific topic, the connection on websocket will wait until a message arrives

```
url := "ws://localhost/ws/sub/test"
conn, _, _ := websocket.DefaultDialer.Dial(url, nil)
defer conn.Close()

for {
  _, msg, err := c.ReadMessage()
  if err != nil {
    log.Fatal(err)
    return
  }
  fmt.Printf("recv: %s\n", msg, t)
}
```

# publish topic

Use the POST method for any topic name (test this time) and set the value.

```
$ curl -X POST http://localhost:8080/ws/pub/test -d 'hello'
```
