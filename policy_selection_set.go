package timepolicy

type PolicySelectionSet struct {
	specs     []*PolicySelection
	evictions EntryWriter

	claims map[*Entry]int
}

type policySelectionSetEvictionWriter struct {
	pss *PolicySelectionSet
}

var _ EntryWriter = &policySelectionSetEvictionWriter{}

func (w *policySelectionSetEvictionWriter) EntriesWritten() int64 {
	panic("should not be used")
}

func (w *policySelectionSetEvictionWriter) WriteEntry(e *Entry) error {
	if w.pss.claims[e] == 1 {
		delete(w.pss.claims, e)
		w.pss.evictions.WriteEntry(e)
	} else {
		w.pss.claims[e]--
	}

	return nil
}

func NewPolicySelectionSet(specs []*PolicySpec, evictions EntryWriter) *PolicySelectionSet {
	pss := &PolicySelectionSet{
		evictions: evictions,
		claims:    map[*Entry]int{},
	}

	evictionsAggregator := &policySelectionSetEvictionWriter{pss: pss}

	for _, spec := range specs {
		pss.specs = append(pss.specs, NewPolicySelection(spec, evictionsAggregator))
	}

	return pss
}

func (p *PolicySelectionSet) Entries() []*Entry {
	var entries []*Entry
	uniqEntries := map[*Entry]struct{}{}

	for _, policySelection := range p.specs {
		for _, e := range policySelection.Entries() {
			if _, known := uniqEntries[e]; known {
				continue
			}

			uniqEntries[e] = struct{}{}
			entries = append(entries, e)
		}
	}

	return entries
}

func (p *PolicySelectionSet) EvaluateEntry(e *Entry) (bool, error) {
	var claims int

	for _, policySelection := range p.specs {
		policyClaimed, err := policySelection.EvaluateEntry(e)
		if err != nil {
			return false, err
		} else if policyClaimed {
			claims++
		}
	}

	if claims == 0 {
		return false, nil
	}

	p.claims[e] = claims

	return true, nil
}
