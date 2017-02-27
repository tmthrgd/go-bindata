package bindata

import "testing"

func TestSafeFunctionName(t *testing.T) {
	var knownFuncs = make(map[string]int)
	t.Run(`foo/bar <-> foo_bar`, func(t *testing.T) {
		name1 := safeFunctionName("foo/bar", knownFuncs)
		name2 := safeFunctionName("foo_bar", knownFuncs)
		if name1 == name2 {
			t.Errorf("name collision")
		}
	})

	t.Run(`reserved words`, func(t *testing.T) {
		name1 := safeFunctionName("json/foo.json", knownFuncs)
		// TODO check that there's no Json or json
		t.Logf("%s", name1)
	})
}
