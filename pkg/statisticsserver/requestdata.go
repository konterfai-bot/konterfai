package statisticsserver

type requestData struct {
	Identifier          string
	Count               int
	Size                string
	IsRobotsTxtViolator string
}

type requestDataSlice []*requestData

// Len is part of sort.Interface.
func (rd requestDataSlice) Len() int {
	return len(rd)
}

// Swap is part of sort.Interface.
func (rd requestDataSlice) Swap(i, j int) {
	rd[i], rd[j] = rd[j], rd[i]
}
