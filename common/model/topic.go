package model

type Topic struct {
    CountReply int
    CountImage int
    Name       string
    Key        string
    Url        string
    Images     []string
}

type TopicPersist struct {
    CountReply           string
    CountImage           string
    CountDownloadedImage string
    Status               string
    Name                 string
    Key                  string
    Url                  string
    FailedImages         string
}
