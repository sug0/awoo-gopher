package awoo

import (
    "time"
)

type TimeStamp int64

func (t TimeStamp) ToTime() time.Time {
    return time.Unix(int64(t), 0)
}

func (t TimeStamp) String() string {
    return t.ToTime().String()
}

func (t TimeStamp) GoString() string {
    return t.String()
}

type BoardDetails struct {
    Name        string `json:"name"`
    Description string `json:"desc"`
    Rules       string `json:"rules"`
}

type Post struct {
    Id         int       `json:"post_id"`
    Board      string    `json:"board"`
    IsOp       bool      `json:"is_op"`
    Comment    string    `json:"comment"`
    DatePosted TimeStamp `json:"date_posted"`
    CapCode    string    `json:"capcode"`
    Title      string    `json:"title"`
    LastBumped TimeStamp `json:"last_bumped"`
    IsLocked   bool      `json:"is_locked"`
    NrReplies  int       `json:"number_of_replies"`
    IsSticky   bool      `json:"sticky"`
    Stickiness int       `json:"stickiness"`
    Hash       string    `json:"hash"`
}
