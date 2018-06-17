package awoo

import (
    "fmt"
    "errors"
    "net/http"
    "net/url"
    "encoding/json"
    "strings"
)

var (
    ErrMsgTooLong = errors.New("awoo: message length exceeds 500 characters")
)

type Client struct {
    host string
    str  string
}

func NewClient(host string, tls bool) *Client {
    if tls {
        host = "https://" + host
    } else {
        host = "http://" + host
    }
    return &Client{host: host, str: host}
}

func (c *Client) String() string {
    return c.str
}

func (c *Client) GoString() string {
    return c.String()
}

func (c *Client) get(path string) (resp *http.Response, err error) {
    return http.Get(fmt.Sprintf("%s%s", c.host, path))
}

func (c *Client) post(path string, form url.Values) (resp *http.Response, err error) {
    return http.Post(fmt.Sprintf("%s%s", c.host, path),
                     "application/x-www-form-urlencoded",
                     strings.NewReader(form.Encode()))
}

func (c *Client) decode(dest interface{}, path string) error {
    resp, err := c.get(path)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    d := json.NewDecoder(resp.Body)

    if err = d.Decode(dest); err != nil {
        return err
    }
    return nil
}

func (c *Client) Boards() ([]string, error) {
    var rsp []string
    if err := c.decode(&rsp, "/api/v2/boards"); err != nil {
        return nil, err
    }
    return rsp, nil
}

func (c *Client) Details(board string) (*BoardDetails, error) {
    var rsp BoardDetails
    if err := c.decode(&rsp, fmt.Sprintf("/api/v2/board/%s/detail", board)); err != nil {
        return nil, err
    }
    return &rsp, nil
}

func (c *Client) Threads(board string) ([]*Post, error) {
    return c.ThreadsPage(board, "0")
}

func (c *Client) ThreadsPage(board, page string) ([]*Post, error) {
    var rsp []*Post
    if err := c.decode(&rsp, fmt.Sprintf("/api/v2/board/%s?page=%s", board, page)); err != nil {
        return nil, err
    }
    return rsp, nil
}

func (c *Client) ThreadMetadata(threadId string) (*Post, error) {
    var rsp Post
    if err := c.decode(&rsp, fmt.Sprintf("/api/v2/thread/%s/metadata", threadId)); err != nil {
        return nil, err
    }
    return &rsp, nil
}

func (c *Client) Replies(threadId string) ([]*Post, error) {
    var rsp []*Post
    if err := c.decode(&rsp, fmt.Sprintf("/api/v2/thread/%s/replies", threadId)); err != nil {
        return nil, err
    }
    return rsp, nil
}

func (c *Client) NewThread(board, title, comment string) error {
    if len(comment) > 500 {
        return ErrMsgTooLong
    }
    resp, err := c.post("/post", url.Values{
        "board": {board},
        "title": {title},
        "comment": {comment},
    })
    if err != nil {
        return err
    }
    resp.Body.Close()
    return nil
}

func (c *Client) NewReply(board, threadId, comment string) error {
    if len(comment) > 500 {
        return ErrMsgTooLong
    }
    resp, err := c.post("/reply", url.Values{
        "board": {board},
        "parent": {threadId},
        "content": {comment},
    })
    if err != nil {
        return err
    }
    resp.Body.Close()
    return nil
}
