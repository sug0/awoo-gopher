package awoo

const DefaultHost = "dangeru.us"
var defcli = NewClient(DefaultHost, true)

func Boards() ([]string, error) {
    return defcli.Boards()
}

func Details(board string) (*BoardDetails, error) {
    return defcli.Details(board)
}

func Threads(board string) ([]*Post, error) {
    return defcli.Threads(board)
}

func ThreadsPage(board, page string) ([]*Post, error) {
    return defcli.ThreadsPage(board, page)
}

func ThreadMetadata(threadId string) (*Post, error) {
    return defcli.ThreadMetadata(threadId)
}

func Replies(threadId string) ([]*Post, error) {
    return defcli.Replies(threadId)
}

func NewThread(board, title, comment string) error {
    return defcli.NewThread(board, title, comment)
}

func NewReply(board, threadId, comment string) error {
    return defcli.NewReply(board, threadId, comment)
}
