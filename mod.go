package main

import (
    "flag"
    "fmt"
    "os"
    "sort"
    "runtime"
    "hash/fnv"

    dgo "github.com/bwmarrin/discordgo"
    "github.com/pkg/browser"
)

var (
    uid, tag, token string
    prn, me bool
)

const (
    opLogin = "auth.keyIn"
    opReadPassword = "auth.read.credentials"
    opReadToken = "auth.read.token"
    opGetRelations = "get.relations"
    opGetMe = "get.@me"
    opGetUser = "get.target"
    opOpenBrowser = "browser.launch"
)


type Err struct {
    Op string
    Cause error
}

func (e Err) Error() string {
    return fmt.Sprintf("%s: %s", e.Op, e.Cause)
}

func (e Err) Hash() uint32 {
    hrx := fnv.New32a()
    hrx.Write([]byte(e.Error()))
    return hrx.Sum32()
}

func (e Err) FCk() {
    if e.Cause != nil {
        fmt.Fprintf(os.Stderr, "fatal: %s\n", e.Error())
        os.Exit(int(e.Hash()))
    }
}

func main() {
    flag.StringVar(&uid, "t", "", "Target user snowflake")
    flag.StringVar(&tag, "g", "", "Target user tag")
    flag.BoolVar(&me, "me", false, "Set the user's account as the target")
    flag.StringVar(&token, "T", "", "Authentication token")
    flag.BoolVar(&prn, "p", false, "Print-only")
    var (
        client *dgo.Session
        err error
    )
    flag.Parse()
    if tag == "" && uid == "" && !me {
        Err{opReadToken, fmt.Errorf("no target provided; please provide one of -g, -t, or -me")}.FCk()
    }
    if token == "" {
        token = os.Getenv("DCPFP_TOKEN")
        if len(token) > 0 && token[0] == '"' && token[len(token)-1] == '"' {
            token = token[1:len(token)-1] // elide quotes
        }
    }
    if token == "" {
        env := "DCPFP_TOKEN"
        if runtime.GOOS == "windows" {
            env = fmt.Sprintf("%%%s%%", env)
        } else {
            env = fmt.Sprintf("$%s", env)
        }
        Err{opReadToken, fmt.Errorf("no token provided; please set one of -T or %s", env)}.FCk()
    }
    if client == nil {
        client, err = dgo.New(token)
        Err{opLogin, err}.FCk()
    }
    var u *dgo.User
    if uid == "" && me {
        // set implicit user self
        var err error
        u, err = client.User("@me")
        Err{opGetMe, err}.FCk()
    }
    if tag != "" && uid == "" && u == nil {
        // search for uid and retrieve user if none given or implied
        relations, err := client.RelationshipsGet()
        Err{opGetRelations, err}.FCk()
        if len(relations) == 0 {
            Err{opGetRelations, fmt.Errorf("you have no friends ;/")}.FCk()
        }
        // find the Discord tag with the closest Levenshtein
        // distance to the search term
        levs := make(map[*dgo.Relationship]int, len(relations))
        for i, r := range relations {
            uname := relations[i].User.Username
            discriminator := relations[i].User.Discriminator
            ltag := fmt.Sprintf("%s#%s", uname, discriminator)
            levs[r] = lev(tag, ltag)
        }
        sort.Slice(relations, func(i, j int) bool {
            return levs[relations[i]] < levs[relations[j]]
        })
        uid = relations[0].User.ID
    }
    if u == nil {
        u, err = client.User(uid)
        Err{opGetUser, err}.FCk()
    }
    uri := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.jpg?size=2048",
        u.ID, u.Avatar)
    if prn {
        fmt.Println(uri)
    } else {
        err := browser.OpenURL(uri)
        Err{opOpenBrowser, err}.FCk()
    }
}
