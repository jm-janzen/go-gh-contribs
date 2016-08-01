package main

import (
    "fmt"
    "log"
    "net/http"
    "io"
    "os"
)

/*
 * TODO
 *  1)  Fix range error if no args
 */

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
        _, err := io.Copy(os.Stdout, resp.Body)
        if err != nil {
            log.Fatal(err)
        }
    }
}

