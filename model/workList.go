package model

// WorkList an array of work that can be sorted by when it was created
type WorkList []*Work

func (wList WorkList) Len() int {
	return len(wList)
}

func (wList WorkList) Less(i, j int) bool {
	return wList[i].When.Before(wList[j].When)
}

func (wList WorkList) Swap(i, j int) {
	wList[i], wList[j] = wList[j], wList[i]
}
