package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "encoding/json"
)

 type GithubUser struct {
    Login   string  `json:"login"`
    Name    string  `json:"name"`
    Email   string  `json:"email"`
    EventsURL   string  `json:"events_url"`
 }

func main() {
    const BASE = "https://api.github.com/users/"
    var name string = ""

    if len(os.Args) <= 1 {
        fmt.Fprintf(os.Stderr, "Usage: ./gh-contribs [username]\n")
        os.Exit(1)
    } else {
        name = os.Args[1]
    }

    var reqURI string = BASE + name
    fmt.Fprintf(os.Stdout, "URI endpoint: %s\n", reqURI)

    resp, err := http.Get(reqURI)

    if err != nil {
        log.Fatal(err)
    } else {
        defer resp.Body.Close()

        decoder := json.NewDecoder(resp.Body)
        var gu GithubUser
        err := decoder.Decode(&gu)
        if err != nil {
            os.Exit(1)
        }

        fmt.Fprintf(os.Stdout, "login:%s, name:%s, email:%s, events_url:%s\n", gu.Login, gu.Name, gu.Email, gu.EventsURL)
    }
}

