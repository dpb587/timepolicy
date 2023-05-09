package timepolicy

import (
	"fmt"
	"sort"
)

type PolicySelection struct {
	spec      *PolicySpec
	evictions EntryWriter

	buckets map[string][]*Entry
}

func NewPolicySelection(spec *PolicySpec, evictions EntryWriter) *PolicySelection {
	return &PolicySelection{
		spec:      spec,
		evictions: evictions,
		buckets:   map[string][]*Entry{},
	}
}

func (p *PolicySelection) Entries() []*Entry {
	res := []*Entry{}

	for _, bucketEntries := range p.buckets {
		res = append(res, bucketEntries...)
	}

	return res
}

func (p *PolicySelection) EvaluateEntry(e *Entry) (bool, error) {
	match, err := p.spec.MatchEntry(e)
	if err != nil {
		return false, err
	} else if !match {
		return false, nil
	}

	var bucketKey string

	if p.spec.bucket != nil {
		bucketKey, err = p.spec.bucket(e)
		if err != nil {
			return false, fmt.Errorf("bucket: %v", err)
		}
	}

	if p.buckets[bucketKey] == nil {
		p.buckets[bucketKey] = []*Entry{e}

		return true, nil
	}

	bucketEntries := p.buckets[bucketKey]

	if p.spec.max == -1 || len(bucketEntries) < p.spec.max {
		nextBucketEntries := append(bucketEntries, e)
		sort.Slice(nextBucketEntries, func(i, j int) bool {
			return nextBucketEntries[i].Time.Before(nextBucketEntries[j].Time)
		})

		p.buckets[bucketKey] = nextBucketEntries

		return true, nil
	} else if p.spec.oldest {
		if e.Time.After(bucketEntries[p.spec.max-1].Time) {
			return false, nil
		}
	} else {
		if e.Time.Before(bucketEntries[0].Time) {
			return false, nil
		}
	}

	nextBucketEntries := append(bucketEntries, e)
	nextBucketEntriesLen := len(nextBucketEntries)

	sort.Slice(nextBucketEntries, func(i, j int) bool {
		return nextBucketEntries[i].Time.Before(nextBucketEntries[j].Time)
	})

	if p.spec.oldest {
		p.evictions.WriteEntry(nextBucketEntries[nextBucketEntriesLen-1])
		bucketEntries = nextBucketEntries[0:p.spec.max]
	} else {
		p.evictions.WriteEntry(nextBucketEntries[0])
		bucketEntries = nextBucketEntries[1:nextBucketEntriesLen]
	}

	p.buckets[bucketKey] = bucketEntries

	return true, nil
}
