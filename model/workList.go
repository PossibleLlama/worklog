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

// RemoveOldRevisions of worklogs
// Returns a list of the latest revisions for each ID
func (wList WorkList) RemoveOldRevisions() []*Work {
	deDuplicated := []*Work{}
	uniqueIDWls := make(map[string][]*Work)
	for _, element := range wList {
		uniqueIDWls[element.ID] = append(uniqueIDWls[element.ID], element)
	}
	for _, wls := range uniqueIDWls {
		highestRevision := -1
		for _, element := range wls {
			if element.Revision > highestRevision {
				highestRevision = element.Revision
			}
		}
		for _, element := range wls {
			if element.Revision == highestRevision {
				deDuplicated = append(deDuplicated, element)
				break
			}
		}
	}
	return deDuplicated
}
