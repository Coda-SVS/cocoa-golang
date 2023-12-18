package util

type ArrayIndex struct {
	Start int
	End   int
}

func NewArrayIndex(start, end int) *ArrayIndex {
	return &ArrayIndex{
		Start: start,
		End:   end,
	}
}

func (a ArrayIndex) Equal(b ArrayIndex) bool {
	return a.Start == b.Start && a.End == b.End
}

func (a ArrayIndex) Size() int {
	return a.End - a.Start
}
