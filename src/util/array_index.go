package util

type Index struct {
	Start int
	End   int
}

func NewIndex(start, end int) *Index {
	return &Index{
		Start: start,
		End:   end,
	}
}

func (a Index) Equal(b Index) bool {
	return a.Start == b.Start && a.End == b.End
}

func (a Index) Size() int {
	return a.End - a.Start
}
