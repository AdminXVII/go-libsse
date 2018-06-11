package main

import (
    "log"
    "net/http"
    "strconv"
    "time"
    "os"

    "github.com/AdminXVII/go-sse"
)

func main() {
    s := sse.NewServer(&sse.Options{
        // Increase default retry interval to 10s.
        RetryInterval: 10 * 1000,
        // CORS headers
        Headers: map[string]string {
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "GET, OPTIONS",
            "Access-Control-Allow-Headers": "Keep-Alive,X-Requested-With,Cache-Control,Content-Type,Last-Event-ID",
        },
        // Print debug info
        Logger: log.New(os.Stdout,
          "go-sse: ",
          log.Ldate|log.Ltime|log.Lshortfile),
        // Add pertinent info first
        InitClient: func(client *sse.Client, LastEventId string) bool {
          client.SendMessage(&sse.Message{Id: "42", Data: "This is the answer to life, to the universe and to everything else",})
          return true
        },
    })

    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.Handle("/events/", s)

    go func () {
        for {
            s.SendMessage(&sse.Message{Data: time.Now().String()})
            time.Sleep(5 * time.Second)
        }
    }()

    go func () {
        i := 0
        for {
            i++
            s.SendMessage(&sse.Message{Data: strconv.Itoa(i)})
            time.Sleep(5 * time.Second)
        }
    }()

    log.Println("Listening at :3000")
    http.ListenAndServe(":3000", nil)
}
