package code

type Gender int

const (
	Gender_Male   Gender = 1
	Gender_Female Gender = 2
)

type Model struct {
	Name   string
	Age    int
	Gender Gender
}

type ModelsByCondition []Model

func (s ModelsByCondition) Len() int {
	return len(s)
}

func (s ModelsByCondition) Less(i, j int) bool {
	if s[i].Age == s[j].Age {
		return s[i].Gender > s[j].Gender
	} else {
		return s[i].Age > s[j].Age
	}
	return false
}

// Swap swaps the elements with indexes i and j.
func (s ModelsByCondition) Swap(i, j int) {
	temp := s[i]
	s[i] = s[j]
	s[j] = temp
}
