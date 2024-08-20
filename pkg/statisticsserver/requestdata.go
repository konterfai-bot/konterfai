package statisticsserver

type RequestData struct {
	Identifier          string
	Count               int
	Size                string
	IsRobotsTxtViolator string
}

type RequestDataSlice []*RequestData

// Len is part of sort.Interface.
func (rd RequestDataSlice) Len() int {
	return len(rd)
}

// Swap is part of sort.Interface.
func (rd RequestDataSlice) Swap(i, j int) {
	rd[i], rd[j] = rd[j], rd[i]
}

// Less is part of sort.Interface. We use count as the primary sort key.
func (rd RequestDataSlice) Less(i, j int) bool {
	return rd[i].Count < rd[j].Count
}
