package hallucinator

// appendHallucination appends a hallucination to the list of hallucinations.
func (h *Hallucinator) appendHallucination(hallucination Hallucination) {
	h.hallucinationLock.Lock()
	defer h.hallucinationLock.Unlock()
	h.hallucinations = append(h.hallucinations, hallucination)
	h.setHallucinationCount()
}

// cleanHallucinations cleans the list of hallucinations and removes hallucinations with requestCount 0.
func (h *Hallucinator) cleanHallucinations() {
	// This function does not have a lock on the hallucinations list. It is expected that the caller has locked the list.
	// This happens in PopHallucination and PopRandomHallucination.
	if h.GetHallucinationCount() < 1 {
		return
	}
	newHallucinations := []Hallucination{}
	for _, hallucination := range h.hallucinations {
		if hallucination.RequestCount > 0 {
			newHallucinations = append(newHallucinations, hallucination)
		}
	}
	h.hallucinations = newHallucinations
	h.setHallucinationCount()
}

// setHallucinationCount sets the hallucination count from the length of the hallucination slice.
func (h *Hallucinator) setHallucinationCount() {
	// This function does not have a lock on the hallucinations list. It is expected that the caller has locked the list.
	// This happens in AppendHallucination and CleanHallucinations.
	h.hallucinationCountLock.Lock()
	defer h.hallucinationCountLock.Unlock()
	h.hallucinationCount = len(h.hallucinations)
}

// GetHallucinationCount returns the current hallucination count.
func (h *Hallucinator) GetHallucinationCount() int {
	h.hallucinationCountLock.Lock()
	defer h.hallucinationCountLock.Unlock()
	return h.hallucinationCount
}
