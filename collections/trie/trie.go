package trie

import (
    "fmt"
)

type (
    Trie struct {
        root *node
        size int
    }
    node struct {
        key   interface{}
        value interface{}
        next  [256]*node
    }
    iterator struct {
        step int
        node *node
        prev *iterator
    }
)

func toBytes(obj interface{}) []byte {
    switch o := obj.(type) {
    case []byte:
        return o
    case string:
        return []byte(o)
    }
    return []byte(fmt.Sprint(obj))
}

func New() *Trie {
    return &Trie{nil, 0}
}
func (t *Trie) Do(handler func(interface{}, interface{}) bool) {
    if t.size > 0 {
        t.root.do(handler)
    }
}
func (t *Trie) Get(key interface{}) interface{} {
    if t.size == 0 {
        return nil
    }

    bs := toBytes(key)
    cur := t.root
    for i := 0; i < len(bs); i++ {
        if cur.next[bs[i]] != nil {
            cur = cur.next[bs[i]]
        } else {
            return nil
        }
    }
    return cur.value
}
func (t *Trie) Has(key interface{}) bool {
    return t.Get(key) != nil
}
func (t *Trie) Init() {
    t.root = nil
    t.size = 0
}
func (t *Trie) Insert(key interface{}, value interface{}) {
    if t.size == 0 {
        t.root = newNode()
    }

    bs := toBytes(key)
    cur := t.root
    for i := 0; i < len(bs); i++ {
        if cur.next[bs[i]] != nil {
            cur = cur.next[bs[i]]
        } else {
            cur.next[bs[i]] = newNode()
            cur = cur.next[bs[i]]
        }
    }
    if cur.key == nil {
        t.size++
    }
    cur.key = key
    cur.value = value
}
func (t *Trie) Len() int {
    return t.size
}
func (t *Trie) Remove(key interface{}) interface{} {
    if t.size == 0 {
        return nil
    }
    bs := toBytes(key)
    cur := t.root

    for i := 0; i < len(bs); i++ {
        if cur.next[bs[i]] != nil {
            cur = cur.next[bs[i]]
        } else {
            return nil
        }
    }

    // TODO: cleanup dead nodes

    val := cur.value

    if cur.value != nil {
        t.size--
        cur.value = nil
        cur.key = nil
    }
    return val
}
func (t *Trie) String() string {
    str := "{"
    i := 0
    t.Do(func(k, v interface{}) bool {
        if i > 0 {
            str += ", "
        }
        str += fmt.Sprint(k, ":", v)
        i++
        return true
    })
    str += "}"
    return str
}

func newNode() *node {
    var next [256]*node
    return &node{nil, nil, next}
}
func (this *node) do(handler func(interface{}, interface{}) bool) bool {
    for i := 0; i < 256; i++ {
        if this.next[i] != nil {
            if this.next[i].key != nil {
                if !handler(this.next[i].key, this.next[i].value) {
                    return false
                }
            }
            if !this.next[i].do(handler) {
                return false
            }
        }
    }
    return true
}
