package main

import (
    "fmt"
    "regexp"
    "strings"
    "strconv"

    "github.com/sugoiuguu/awoo-gopher/awoo"
    "github.com/prologic/go-gopher"
)

type Handler func(gopher.ResponseWriter, *gopher.Request, []string)

type route struct {
    re *regexp.Regexp
    h  Handler
}

type Router struct {
    routes []*route
}

func main() {
    router := new(Router)

    router.HandleFunc(`/thread/(\d+)`, thread)
    router.HandleFunc(`/board/(.+)/(\d+)`, threads)
    router.HandleFunc(`/`, index)

    err := gopher.ListenAndServe("localhost:7070", router)
    if err != nil {
        panic(err)
    }
}

func (r *Router) HandleFunc(re string, h Handler) {
    r.routes = append(r.routes, &route{regexp.MustCompile(re), h})
}

func (r *Router) ServeGopher(w gopher.ResponseWriter, rq *gopher.Request) {
    for _,route := range r.routes {
        if m := route.re.FindStringSubmatch(rq.Selector); m != nil {
            route.h(w, rq, m)
            return
        }
    }
}

func index(w gopher.ResponseWriter, r *gopher.Request, p []string) {
    boards, err := awoo.Boards()
    if err != nil {
        gopher.Error(w, err.Error())
        return
    }

    w.WriteInfo(awoo.DefaultHost + " boards")

    for _,board := range boards {
        w.WriteItem(&gopher.Item{
            Type: gopher.DIRECTORY,
            Selector: fmt.Sprintf("/board/%s/0", board),
            Description: board,
        })
    }
}

func threads(w gopher.ResponseWriter, r *gopher.Request, p []string) {
    d, err := awoo.Details(p[1])
    if err != nil {
        gopher.Error(w, err.Error())
        return
    }

    w.WriteInfo(fmt.Sprintf("/%s/ - %s", d.Name, d.Description))
    for _,line := range strings.Split(d.Rules, "\n") {
        w.WriteInfo(line)
    }
    w.WriteInfo("")

    thrs, err := awoo.ThreadsPage(p[1], p[2])
    if err != nil {
        gopher.Error(w, err.Error())
        return
    }

    var desc string
    var showBoard bool

    if p[1] == "all" {
        showBoard = true
    }

    for _,t := range thrs {
        if showBoard {
            desc = fmt.Sprintf("Nr. %d | /%s/ | %d Replies | %s | %s",
                                     t.Id, t.Board, t.NrReplies, t.DatePosted, t.Title)
        } else {
            desc = fmt.Sprintf("Nr. %d | %d Replies | %s | %s",
                                     t.Id, t.NrReplies, t.DatePosted, t.Title)
        }
        w.WriteItem(&gopher.Item{
            Type: gopher.DIRECTORY,
            Selector: fmt.Sprintf("/thread/%d", t.Id),
            Description: desc,
        })
        for _,line := range strings.Split(t.Comment, "\n") {
            w.WriteInfo(line)
        }
    }

    next,_ := strconv.Atoi(p[2])
    w.WriteInfo("")
    w.WriteItem(&gopher.Item{
        Type: gopher.DIRECTORY,
        Selector: fmt.Sprintf("/board/%s/%d", p[1], next + 1),
        Description: fmt.Sprintf("%s - page %d", p[1], next + 2),
    })
}

func thread(w gopher.ResponseWriter, r *gopher.Request, p []string) {
    rp, err := awoo.Replies(p[1])
    if err != nil {
        gopher.Error(w, err.Error())
        return
    }

    w.WriteInfo(fmt.Sprintf("Replies for thread /%s/ - %d", rp[0].Board, rp[0].Id))
    w.WriteInfo("")

    var desc string

    for _,p := range rp {
        desc = fmt.Sprintf("Nr. %d | %s",
                           p.Id, p.DatePosted)
        w.WriteItem(&gopher.Item{
            Type: gopher.DIRECTORY,
            Selector: "/",
            Description: desc,
        })
        for _,line := range strings.Split(p.Comment, "\n") {
            w.WriteInfo(line)
        }
    }
}
