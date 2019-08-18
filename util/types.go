package util

import (
    "crypto/sha1"
    "encoding/hex"

    "github.com/Tomilla/imagespider/config"
)

func Hash(s string) string {

    // The pattern for generating a hash is `sha1.New()`,
    // `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
    // Here we start with a new hash.
    h := sha1.New()

    // `Write` expects bytes. If you have a string `s`,
    // use `[]byte(s)` to coerce it to bytes.
    h.Write([]byte(s))

    // This gets the finalized hash result as a byte
    // slice. The argument to `Sum` can be used to append
    // to an existing byte slice: it usually isn't needed.
    bs := h.Sum(nil)

    // SHA1 values are often printed in hex, for example
    // in git commits. Use the `%x` format verb to convert
    // a hash results to a hex string.
    // fmt.Println(s)
    return hex.EncodeToString(bs)

}

type OrderQueryMap struct {
    M    map[string]string
    keys []string
}

func NewOrderQueryMap() *OrderQueryMap {
    return &OrderQueryMap{
        M:    map[string]string{},
        keys: []string{},
    }
}
func (o *OrderQueryMap) Init(anotherM [][]string) {
    for _, item := range anotherM {
        key, val := item[0], item[1]
        o.Set(key, val)
    }
}

func (o *OrderQueryMap) Set(k string, v string) {
    // this will not update insertion order for key k if it was already in the map
    _, present := o.M[k]
    o.M[k] = v
    if !present {
        o.keys = append(o.keys, k)
    }
}
func (o OrderQueryMap) Iterate(f func(string, string)) {
    config.L.Debugf("keys(in insert order): %v", o.keys)
    for _, key := range o.keys {
        val := o.M[key]
        f(key, val)
    }
}

