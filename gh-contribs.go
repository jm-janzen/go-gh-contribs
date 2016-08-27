package main

import (
    "fmt"
    "net/http"
    "os"
    "encoding/json"
    "log"
    "time"
    "strconv"
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
    req,_ := http.NewRequest("GET",
    reqURI +
        "?client_id=" + "7f46d3c17268b610d86d" +
        "&client_secret=" + "44530b218e3e6ce30084ec7e51a36b585dc41b59",
    nil)
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

            //var responses []GithubUserEvents
            /*
             * XXX
             *  Fails after page=10
             */
            var totalStreaks int = 0
            for i := 0; i < 100; i++ {
                reqEventsURI := reqURI + "/events?page=" + strconv.Itoa(i + 1) +
                    //"&per_page=100" +
                    "&client_id=" + "7f46d3c17268b610d86d" + 
                    "&redirect_uri=http://jmjanzen.com/gh-contribs" + 
                    "&client_secret=" + "" /* XXX insert id here */
                fmt.Println("querying: %s", reqEventsURI)
                req,_ = http.NewRequest("GET",  reqEventsURI, nil)
                resp,err := client.Do(req)
                // TODO load up responses arr
                if err != nil {
                    os.Exit(1)
                } else {
                    var gue GithubUserEvents

                    decoder := json.NewDecoder(resp.Body)
                    err := decoder.Decode(&gue)
                    if err != nil {
                        fmt.Fprintf(os.Stdout, "%s", decoder)
                        log.Print(err)
                        i = 100
                    }

                    var currentStreak, highestStreak, day int = 0, 0, 0
                    for _, g := range gue {
                        fmt.Fprintf(os.Stdout, "[%02d] %s:\t%s\n", currentStreak, g.CreatedAt.String()[0:10], g.Type)

                        if day == g.CreatedAt.Day() {
                            currentStreak++
                        } else {
                            fmt.Println("new day")
                            if highestStreak < currentStreak {
                                fmt.Println(highestStreak, currentStreak)
                                highestStreak = currentStreak
                                totalStreaks += highestStreak
                            }
                            currentStreak = 0
                        }
                        day = g.CreatedAt.Day()
                    }

                    // ie highest number of commits in a day
                    fmt.Fprintf(os.Stdout, "highest streak: %d\n", highestStreak)

                }
                // don't over-spam Gh API
                //time.Sleep(1 * time.Second)
            }
            fmt.Fprintf(os.Stdout, "total streaks: %d\n", totalStreaks)
    }
}

