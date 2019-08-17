package set

type (
    Set struct {
        hash map[interface{}]nothing
    }

    nothing struct{}
)

// Create a new set
func New(initial ...interface{}) *Set {
    s := &Set{make(map[interface{}]nothing)}

    for _, v := range initial {
        s.Insert(v)
    }

    return s
}

// Find the difference between two sets
func (s *Set) Difference(set *Set) *Set {
    n := make(map[interface{}]nothing)

    for k := range s.hash {
        if _, exists := set.hash[k]; !exists {
            n[k] = nothing{}
        }
    }

    return &Set{n}
}

// Call f for each item in the set
func (s *Set) Do(f func(interface{}) bool) {
    for k := range s.hash {
        f(k)
    }
}

// Test to see whether or not the element is in the set
func (s *Set) Has(element interface{}) bool {
    _, exists := s.hash[element]
    return exists
}

// Add an element to the set
func (s *Set) Insert(element interface{}) {
    s.hash[element] = nothing{}
}

// Find the intersection of two sets
func (s *Set) Intersection(set *Set) *Set {
    n := make(map[interface{}]nothing)

    for k := range s.hash {
        if _, exists := set.hash[k]; exists {
            n[k] = nothing{}
        }
    }

    return &Set{n}
}

// Return the number of items in the set
func (s *Set) Len() int {
    return len(s.hash)
}

// Test whether or not this set is a proper subset of "set"
func (s *Set) ProperSubsetOf(set *Set) bool {
    return s.SubsetOf(set) && s.Len() < set.Len()
}

// Remove an element from the set
func (s *Set) Remove(element interface{}) {
    delete(s.hash, element)
}

// Test whether or not this set is a subset of "set"
func (s *Set) SubsetOf(set *Set) bool {
    if s.Len() > set.Len() {
        return false
    }
    for k := range s.hash {
        if _, exists := set.hash[k]; !exists {
            return false
        }
    }
    return true
}

// Find the union of two sets
func (s *Set) Union(set *Set) *Set {
    n := make(map[interface{}]nothing)

    for k := range s.hash {
        n[k] = nothing{}
    }
    for k := range set.hash {
        n[k] = nothing{}
    }

    return &Set{n}
}
