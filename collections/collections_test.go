package collections

import (
    "testing"

    "github.com/Tomilla/imagespider/collections/set"
)

func Test(t *testing.T) {
    s := set.New()

    s.Insert(5)
    s.Insert(4)
    s.Insert(3)
    s.Insert(2)
    s.Insert(1)
    GetRange(s, 0, 5)
    GetRange(s, 0, 2)

}
