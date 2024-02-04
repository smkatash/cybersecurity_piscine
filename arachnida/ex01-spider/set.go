package main

func (s StringSet) Add(str string) {
    s[str] = struct{}{}
}

func (s StringSet) Remove(str string) {
    delete(s, str)
}

func (s StringSet) Contains(str string) bool {
    _, exists := s[str]
    return exists
}
