package main

import (
    "fmt"
    "net/http"
    "os"
    "encoding/json"
    "log"
    "time"
)

 /*
  * Internal model of JSON from
  * Github API /users/ ...
  */
 type GithubUser struct {
    Login       string  `json:"login"`
    Name        string  `json:"name"`
    Email       string  `json:"email"`
    EventsURL   string  `json:"events_url"`
 }

 type GithubUserEvents []struct {
    Id          string  `json:"id"`
    Type        string  `json:"type"`
    CreatedAt   time.Time `json:"created_at"`
 }
 type Payload struct {
    PushId      string  `json:"push_id"`
    Size        int     `json:"size"`
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

    client := &http.Client{}
    req,_ := http.NewRequest("GET", reqURI, nil)
    req.Header.Add("User-Agent", "gh-contribs")
    resp, err := client.Do(req)

    if err != nil {
        os.Exit(1)
    } else {
        defer resp.Body.Close()

        decoder := json.NewDecoder(resp.Body)
        var gu GithubUser
        err := decoder.Decode(&gu)
        if err != nil {
            os.Exit(1)
        }

        fmt.Fprintf(os.Stdout,
            "\nlogin:\t%s, \nname:\t%s, \nemail:\t%s, \nevents:\t%s\n",
            gu.Login,
            gu.Name,
            gu.Email,
            gu.EventsURL)

        req,_ := http.NewRequest("GET", reqURI + "/events?per_page=100&page=0", nil)
        resp, err := client.Do(req)

        if err != nil {
            os.Exit(1)
        } else {
            var gue GithubUserEvents

            decoder := json.NewDecoder(resp.Body)
            err := decoder.Decode(&gue)
            if err != nil {
                log.Fatal(err)
                os.Exit(1)
            }

            var currentStreak, highestStreak, day int = 0, 0, 0
            for _, g := range gue {

                if day == g.CreatedAt.Day() {
                    //*XXX*/fmt.Println(day, "==", g.CreatedAt.Day())
                    currentStreak++
                } else {
                    //*XXX*/fmt.Println(day, "!=", g.CreatedAt.Day())
                    if highestStreak < currentStreak {
                        //*XXX*/fmt.Println(highestStreak, "<", currentStreak)
                        highestStreak = currentStreak
                    }
                    currentStreak = 0
                }
                day = g.CreatedAt.Day()
                fmt.Fprintf(os.Stdout, "[%02d] %s:\t%s\n",
                    currentStreak + 1,
                    g.CreatedAt.String()[0:10],
                    g.Type)
            }

            // ie highest number of commits in a day
            fmt.Fprintf(os.Stdout, "streak: %d\n", highestStreak)

        }
    }
}

