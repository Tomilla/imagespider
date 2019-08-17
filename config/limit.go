package config

type Limit struct {
    ImagePathLenLimit int `json:"imagePathLenLimit"`
    PostNameLenLimit  int `json:"postNameLenLimit"`
}
