package timepolicy

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func mustParseRFC3339(v string) time.Time {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		panic(err)
	}

	return t
}

func stubNow() time.Time {
	return mustParseRFC3339("2023-01-01T01:02:03Z")
}

func TestBasicNewest(t *testing.T) {
	Now = stubNow

	spec, err := ParsePolicySpecString("test", "7d;by=day;max=2")
	if err != nil {
		t.Fatalf("expected `nil` but got: %v", err)
	}

	deferredEvictions := bytes.NewBuffer(nil)
	deferredEvictionsWriter := NewEntryWriter(deferredEvictions)

	ps := NewPolicySelection(spec, deferredEvictionsWriter)

	{ // match; new; accept
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-0",
			Time: mustParseRFC3339("2022-12-31T18:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := true, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	{ // match; new; accept
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-1",
			Time: mustParseRFC3339("2022-12-31T09:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := true, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	{ // match; evict; accept
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-2",
			Time: mustParseRFC3339("2022-12-31T12:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := true, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	{ // match; stale; skip
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-3",
			Time: mustParseRFC3339("2022-12-31T06:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := false, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	{ // no-match; skip
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-4",
			Time: mustParseRFC3339("2020-12-31T06:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := false, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	if _e, _a := true, strings.Contains(deferredEvictions.String(), "entry-1\n"); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}

	entries := ps.Entries()
	if _e, _a := 2, len(entries); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "entry-2", entries[0].Raw; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "entry-0", entries[1].Raw; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}

func TestBasicOldest(t *testing.T) {
	Now = stubNow

	spec, err := ParsePolicySpecString("test", "7d;by=day;max=2;oldest")
	if err != nil {
		t.Fatalf("expected `nil` but got: %v", err)
	}

	deferredEvictions := bytes.NewBuffer(nil)
	deferredEvictionsWriter := NewEntryWriter(deferredEvictions)

	ps := NewPolicySelection(spec, deferredEvictionsWriter)

	{ // match; new; accept
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-0",
			Time: mustParseRFC3339("2022-12-31T18:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := true, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	{ // match; new; accept
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-1",
			Time: mustParseRFC3339("2022-12-31T09:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := true, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	{ // match; evict; accept
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-2",
			Time: mustParseRFC3339("2022-12-31T12:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := true, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	{ // match; evict; accept
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-3",
			Time: mustParseRFC3339("2022-12-31T06:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := true, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	{ // no-match; skip
		selected, err := ps.EvaluateEntry(&Entry{
			Raw:  "entry-4",
			Time: mustParseRFC3339("2020-12-31T06:00:00Z"),
		})
		if err != nil {
			t.Fatalf("expected `nil` but got: %v", err)
		} else if _e, _a := false, selected; _e != _a {
			t.Fatalf("expected `%v` but got: %v", _e, _a)
		}
	}

	if _e, _a := true, strings.Contains(deferredEvictions.String(), "entry-0\n"); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := true, strings.Contains(deferredEvictions.String(), "entry-2\n"); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}

	entries := ps.Entries()
	if _e, _a := 2, len(entries); _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "entry-3", entries[0].Raw; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	} else if _e, _a := "entry-1", entries[1].Raw; _e != _a {
		t.Fatalf("expected `%v` but got: %v", _e, _a)
	}
}
