package data_structure

type SimpleSet struct {
	dict map[string]struct{}
}

func NewSimpleSet() *SimpleSet {
	return &SimpleSet{
		dict: make(map[string]struct{}),
	}
}

func (s *SimpleSet) Add(members ...string) int {
	added := 0
	for _ , mem := range members {
		if _ , exist := s.dict[mem]; !exist {
			s.dict[mem] = struct{}{}
			added++
		}
	}
	return added
}

func (s *SimpleSet) Rem(members ...string) int {
	removed := 0
	for _ , mem := range members {
		if _ , exist := s.dict[mem] ; exist {
			delete(s.dict , mem)
			removed++
		}
	}
	return removed
}

func (s *SimpleSet) IsMember(member string) int {
	_ , exist := s.dict[member]
	if exist {
		return 1
	}
	return 0
}

func (s *SimpleSet) Members() []string {
	mems := make([]string , 0)
	for k , _ := range s.dict {
		mems = append(mems, k)
	}
	return mems
}