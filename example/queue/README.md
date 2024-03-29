## websocket queue subscription

![output2](https://user-images.githubusercontent.com/42143893/50048371-8d663880-010d-11e9-833a-eeb3cbdcf294.gif)

# run queue subscription

queue subscription will only receive one message even if there are multiple workers.

worker1 groupA:

```
url := "ws://localhost/ws/sub/test/groupA"
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

worker2 groupA:

```
url := "ws://localhost/ws/sub/test/groupA"
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
$ curl -X POST http://localhost:8080/pub/test -d 'hello'
```
