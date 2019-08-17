package splay

import (
    "fmt"
)

type (
    Any interface{}
    LessFunc func(interface{}, interface{}) bool
    VisitFunc func(interface{}) bool

    node struct {
        value               Any
        parent, left, right *node
    }
    nodei struct {
        step int
        node *node
        prev *nodei
    }

    SplayTree struct {
        length int
        root   *node
        less   LessFunc
    }
)

// Create a new splay tree, using the less function to determine the order.
func New(less LessFunc) *SplayTree {
    return &SplayTree{0, nil, less}
}

// Get the first value from the collection. Returns nil if empty.
func (s *SplayTree) First() Any {
    if s.length == 0 {
        return nil
    }

    n := s.root
    for n.left != nil {
        n = n.left
    }
    return n.value
}

// Get the last value from the collection. Returns nil if empty.
func (s *SplayTree) Last() Any {
    if s.length == 0 {
        return nil
    }
    n := s.root
    for n.right != nil {
        n = n.right
    }
    return n.value
}

// Get an item from the splay tree
func (s *SplayTree) Get(item Any) Any {
    if s.length == 0 {
        return nil
    }

    n := s.root
    for n != nil {
        if s.less(item, n.value) {
            n = n.left
            continue
        }

        if s.less(n.value, item) {
            n = n.right
            continue
        }

        s.splay(n)
        return n.value
    }
    return nil
}
func (s *SplayTree) Has(value Any) bool {
    return s.Get(value) != nil
}
func (s *SplayTree) Init() {
    s.length = 0
    s.root = nil
}
func (s *SplayTree) Add(value Any) {
    if s.length == 0 {
        s.root = &node{value, nil, nil, nil}
        s.length = 1
        return
    }

    n := s.root
    for {
        if s.less(value, n.value) {
            if n.left == nil {
                n.left = &node{value, n, nil, nil}
                s.length++
                n = n.left
                break
            }
            n = n.left
            continue
        }

        if s.less(n.value, value) {
            if n.right == nil {
                n.right = &node{value, n, nil, nil}
                s.length++
                n = n.right
                break
            }
            n = n.right
            continue
        }

        n.value = value
        break
    }
    s.splay(n)
}
func (s *SplayTree) PreOrder(visit VisitFunc) {
    if s.length == 1 {
        return
    }
    i := &nodei{0, s.root, nil}
    for i != nil {
        switch i.step {
        // Value
        case 0:
            i.step++
            if !visit(i.node.value) {
                break
            }
        // Left
        case 1:
            i.step++
            if i.node.left != nil {
                i = &nodei{0, i.node.left, i}
            }
        // Right
        case 2:
            i.step++
            if i.node.right != nil {
                i = &nodei{0, i.node.right, i}
            }
        default:
            i = i.prev
        }
    }
}
func (s *SplayTree) InOrder(visit VisitFunc) {
    if s.length == 1 {
        return
    }
    i := &nodei{0, s.root, nil}
    for i != nil {
        switch i.step {
        // Left
        case 0:
            i.step++
            if i.node.left != nil {
                i = &nodei{0, i.node.left, i}
            }
        // Value
        case 1:
            i.step++
            if !visit(i.node.value) {
                break
            }
        // Right
        case 2:
            i.step++
            if i.node.right != nil {
                i = &nodei{0, i.node.right, i}
            }
        default:
            i = i.prev
        }
    }
}
func (s *SplayTree) PostOrder(visit VisitFunc) {
    if s.length == 1 {
        return
    }
    i := &nodei{0, s.root, nil}
    for i != nil {
        switch i.step {
        // Left
        case 0:
            i.step++
            if i.node.left != nil {
                i = &nodei{0, i.node.left, i}
            }
        // Right
        case 1:
            i.step++
            if i.node.right != nil {
                i = &nodei{0, i.node.right, i}
            }
        // Value
        case 2:
            i.step++
            if !visit(i.node.value) {
                break
            }
        default:
            i = i.prev
        }
    }
}
func (s *SplayTree) Do(visit VisitFunc) {
    s.InOrder(visit)
}
func (s *SplayTree) Len() int {
    return s.length
}
func (s *SplayTree) Remove(value Any) {
    if s.length == 0 {
        return
    }

    n := s.root
    for n != nil {
        if s.less(value, n.value) {
            n = n.left
            continue
        }
        if s.less(n.value, value) {
            n = n.right
            continue
        }

        // First splay the parent node
        if n.parent != nil {
            s.splay(n.parent)
        }

        // No children
        if n.left == nil && n.right == nil {
            // guess we're the root node
            if n.parent == nil {
                s.root = nil
                break
            }
            if n.parent.left == n {
                n.parent.left = nil
            } else {
                n.parent.right = nil
            }
        } else if n.left == nil {
            // root node
            if n.parent == nil {
                s.root = n.right
                break
            }
            if n.parent.left == n {
                n.parent.left = n.right
            } else {
                n.parent.right = n.right
            }
        } else if n.right == nil {
            // root node
            if n.parent == nil {
                s.root = n.left
                break
            }
            if n.parent.left == n {
                n.parent.left = n.left
            } else {
                n.parent.right = n.left
            }
        } else {
            // find the successor
            s := n.right
            for s.left != nil {
                s = s.left
            }

            np := n.parent
            nl := n.left
            nr := n.right

            sp := s.parent
            sr := s.right

            // Update parent
            s.parent = np
            if np == nil {
                s.parent = s
                // s.root = s
            } else {
                if np.left == n {
                    np.left = s
                } else {
                    np.right = s
                }
            }

            // Update left
            s.left = nl
            s.left.parent = s

            // Update right
            if nr != s {
                s.right = nr
                s.right.parent = s
            }

            // Update successor parent
            if sp.left == s {
                sp.left = sr
            } else {
                sp.right = sr
            }
        }

        break
    }

    if n != nil {
        s.length--
    }
}
func (s *SplayTree) String() string {
    if s.length == 0 {
        return "{}"
    }
    return s.root.String()
}

// Splay a node in the tree (send it to the top)
func (s *SplayTree) splay(n *node) {
    // Already root, nothing to do
    if n.parent == nil {
        s.root = n
        return
    }

    p := n.parent
    g := p.parent

    // Zig
    if p == s.root {
        if n == p.left {
            p.rotateRight()
        } else {
            p.rotateLeft()
        }
    } else {
        // Zig-zig
        if n == p.left && p == g.left {
            g.rotateRight()
            p.rotateRight()
        } else if n == p.right && p == g.right {
            g.rotateLeft()
            p.rotateLeft()
            // Zig-zag
        } else if n == p.right && p == g.left {
            p.rotateLeft()
            g.rotateRight()
        } else if n == p.left && p == g.right {
            p.rotateRight()
            g.rotateLeft()
        }
    }
    s.splay(n)
}

// Swap two nodes in the tree
func (s *SplayTree) swap(n1, n2 *node) {
    p1 := n1.parent
    l1 := n1.left
    r1 := n1.right

    p2 := n2.parent
    l2 := n2.left
    r2 := n2.right

    // Update node links
    n1.parent = p2
    n1.left = l2
    n1.right = r2

    n2.parent = p1
    n2.left = l1
    n2.right = r1

    // Update parent links
    if p1 != nil {
        if p1.left == n1 {
            p1.left = n2
        } else {
            p1.right = n2
        }
    }
    if p2 != nil {
        if p2.left == n2 {
            p2.left = n1
        } else {
            p2.right = n1
        }
    }

    if n1 == s.root {
        s.root = n2
    } else if n2 == s.root {
        s.root = n1
    }
}

// Node methods
func (this *node) String() string {
    str := "{" + fmt.Sprint(this.value) + "|"
    if this.left != nil {
        str += this.left.String()
    }
    str += "|"
    if this.right != nil {
        str += this.right.String()
    }
    str += "}"
    return str
}
func (this *node) rotateLeft() {
    parent := this.parent
    pivot := this.right
    child := pivot.left

    if pivot == nil {
        return
    }

    // Update the parent
    if parent != nil {
        if parent.left == this {
            parent.left = pivot
        } else {
            parent.right = pivot
        }
    }

    // Update the pivot
    pivot.parent = parent
    pivot.left = this

    // Update the child
    if child != nil {
        child.parent = this
    }

    // Update this
    this.parent = pivot
    this.right = child
}
func (this *node) rotateRight() {
    parent := this.parent
    pivot := this.left
    child := pivot.right

    if pivot == nil {
        return
    }

    // Update the parent
    if parent != nil {
        if parent.left == this {
            parent.left = pivot
        } else {
            parent.right = pivot
        }
    }

    // Update the pivot
    pivot.parent = parent
    pivot.right = this

    if child != nil {
        child.parent = this
    }

    // Update this
    this.parent = pivot
    this.left = child
}
