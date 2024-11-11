package set

type Set[E comparable] struct {
  m map[E]struct{}
}

func New[E comparable]() *Set[E] {
  return &Set[E]{m: make(map[E]struct{})}
}

func (s *Set[E]) Add(v E) {
  s.m[v] = struct{}{}
}

func (s *Set[E]) Contains(v E) bool {
  _, ok := s.m[v]
  return ok
}

func Union[E comparable](s1, s2 *Set[E]) *Set[E] {
  r := New[E]()

  for v := range s1.m {
    r.Add(v)
  }
  for v := range s2.m {
    r.Add(v)
  }

  return r
}

func (s *Set[E]) Push(f func(E) bool) {
  for v := range s.m {
    if !f(v) {
      return
    }
  }
}

func (s *Set[E]) Pull() (func() (E, bool), func()) {
  ch := make(chan E)
  stopCh := make(chan bool)

  go func() {
    defer close(ch)
    for v := range s.m {
      select {
      case ch <- v:
      case <- stopCh:
        return
      }
    }
  }()

  next := func() (E, bool) {
    v, ok := <-ch
    return v, ok
  }

  stop := func() {
    close(stopCh)
  }

  return next, stop
}
