package main

import (
    "flag"
    "fmt"
    "os"
    "sort"
    "runtime"

    dgo "github.com/bwmarrin/discordgo"
    "github.com/pkg/browser"
)

var (
    uid, tag, token string
    prn bool
)

func main() {
    flag.StringVar(&uid, "t", "", "Target user snowflake")
    flag.StringVar(&tag, "g", "", "Target user tag")
    flag.StringVar(&token, "T", "", "Authentication token")
    flag.BoolVar(&prn, "p", false, "Print-only")
    flag.Parse()
    if tag == "" && uid == "" {
        fmt.Fprintf(os.Stderr, "fatal: no target provided; please provide one of -g or -t\n")
        os.Exit(1)
    }
    if token == "" {
        token = os.Getenv("DCPFP_TOKEN")
    }
    if token == "" {
        env := "DCPFP_TOKEN"
        if runtime.GOOS == "windows" {
            env = fmt.Sprintf("%%%s%%", env)
        } else {
            env = fmt.Sprintf("$%s", env)
        }
        fmt.Fprintf(os.Stderr, "fatal: no token provided; please set one of -T or %s\n", env)
        os.Exit(2)
    }
    client, err := dgo.New(token)
    if tag != "" {
        // emplace uid if none given
        relations, err := client.RelationshipsGet()
        if err != nil {
            fmt.Fprintf(os.Stderr, "fatal: retrieve relations: %s\n", err)
            os.Exit(3)
        }
        if len(relations) == 0 {
            fmt.Fprintf(os.Stderr, "fatal: you have no friends ;/\n")
            os.Exit(4)
        }
        // find the Discord nick#tag with the closest Damerau—Levenshtein
        // distance to the search term
        levs := make(map[*dgo.Relationship]int, len(relations))
        for i, r := range relations {
            uname := relations[i].User.Username
            discriminator := relations[i].User.Discriminator
            ltag := fmt.Sprintf("%s#%s", uname, discriminator)
            levs[r] = dlev(tag, ltag)
        }
        sort.Slice(relations, func(i, j int) bool {
            return levs[relations[i]] < levs[relations[j]]
        })
        uid = relations[0].User.ID
    }
    if err != nil {
        fmt.Fprintf(os.Stderr, "fatal: login failure: %s\n", err)
        os.Exit(4)
    }
    u, err := client.User(uid)
    if err != nil {
        fmt.Fprintf(os.Stderr, "fatal: retrieve user: %s\n", err)
        os.Exit(4)
    }
    uri := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.jpg?size=2048",
        u.ID, u.Avatar)
    if prn {
        fmt.Println(uri)
    } else {
        err := browser.OpenURL(uri)
        if err != nil {
            fmt.Fprintf(os.Stderr, "failure spawning browser: %s\n", err)
            os.Exit(5)
        }
    }
}