package pagination

type ranger struct {
	minIndex   int
	maxIndex   int
	segmentLen int
	index      int
}

func (r *ranger) Init(minIndex, maxIndex, segmentLen, index int) *ranger {
	if minIndex > maxIndex {
		panic("ranger's minimum index is greater than maximum index!")
	}
	if minIndex > index {
		panic("ranger's minimum index is greater than cursor index!")
	}
	if maxIndex < index {
		panic("ranger's maximum index is less than cursor index!")
	}

	r.minIndex = minIndex
	r.maxIndex = maxIndex
	r.segmentLen = segmentLen
	r.index = index

	return r
}

func (r *ranger) SetMin(segmentLen int) *ranger {
	// TODO::
	return r
}
func (r *ranger) SetMax(segmentLen int) *ranger {
	// TODO::
	return r
}
func (r *ranger) SetLen(segmentLen int) *ranger {
	// TODO::
	return r
}
func (r *ranger) SetIndex(segmentLen int) *ranger {
	// TODO::
	return r
}

func (r *ranger) getSegmentNum() int {

	return (r.maxIndex - r.minIndex) / r.segmentLen
}

func (r *ranger) GetRange() (startIndex, endIndex int) {
	startIndex = ((r.index - r.minIndex) / r.segmentLen) * r.segmentLen
	endIndex = startIndex + r.segmentLen
	if endIndex > r.maxIndex {
		endIndex = r.maxIndex
	}

	return
}

func (r *ranger) GetRangeByIndex() (index int) {
	// TODO::

	return
}
func (r *ranger) GetRangeBySegNum() (segmentNum int) {
	// TODO::

	return
}
