// heap.go
package main

type MinHeap []*Img

func (mh MinHeap) Len() int { return len(mh) }

func (mh MinHeap) Less(i, j int) bool {
	return mh[i].error > mh[j].error
}

func (mh MinHeap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
}

func (mh *MinHeap) Pop() interface{} {
	old := *mh
	n := len(old)
	img := old[n-1]
	*mh = old[0 : n-1]
	return img
}

func (mh *MinHeap) Push(x interface{}) {
	img := x.(*Img)
	*mh = append(*mh, img)
}
