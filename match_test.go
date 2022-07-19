package find

import (
	"testing"
)

func TestLongestFixedPart(t *testing.T) {
	for name, tc := range map[string]struct {
		pattern string
		want    string
	}{
		"full":                        {"rfuiiufuier", "rfuiiufuier"},
		"wildcard":                    {"rfu*iiufuier", "iiufuier"},
		"escaped-wildcard":            {"rfu\\*ii\\\\", "rfu*ii\\"},
		"escaped-backslash":           {"rfu\\\\*ii\\\\", "rfu\\"},
		"escaped-backslash-wildcard":  {"rfu\\\\\\*ii\\\\", "rfu\\*ii\\"},
		"escaped-backslash2":          {"rfu\\\\\\\\*ii\\\\", "rfu\\\\"},
		"escaped-backslash2-wildcard": {"rfu\\\\\\\\\\*ii\\\\", "rfu\\\\*ii\\"},
	} {
		t.Run(name, func(t *testing.T) {
			part := LongestFixedPart(tc.pattern)
			if part != tc.want {
				t.Errorf("pattern %v: expected %v, got %v", tc.pattern, tc.want, part)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	for name, tc := range map[string]struct {
		line       string
		oldPattern string
		newPattern string
		want       string
		modified   bool
	}{
		"basic":                           {"hej hej hej", "*hej*", "*hopp*", "hopp hopp hopp", true},
		"none":                            {"hej hej hej", "*tjo*", "*hopp*", "", false},
		"duplicates":                      {"hhhh", "*hh*", "*h*", "hh", true},
		"escaped-wildcard":                {"hh*hh", "*h\\*h*", "*g\\*g*", "hg*gh", true},
		"escaped-backslash":               {"hh\\hh", "*h\\\\h*", "*g\\\\g*", "hg\\gh", true},
		"escaped-backslash-wildcard":      {"hh\\*hh", "*h\\\\\\*h*", "*g\\\\g*", "hg\\gh", true},
		"escaped-backslash-wildcard-fail": {"hh\\hh", "*h\\\\\\*h*", "*g\\\\\\*g*", "", false},
		"escaped-backslash-real-wildcard": {"hh\\kh", "*h\\\\*h*", "*g\\\\*g*", "hg\\kg", true},
		"leading-asterix": {
			`	handleClusterAbort(*fywire.ClusterAbort, factory.NodeID) (lnfactory_pb.Message, error)`,
			`*\*fywire.*`,
			`*\*lnfactory_pb.*`,
			`	handleClusterAbort(*lnfactory_pb.ClusterAbort, factory.NodeID) (lnfactory_pb.Message, error)`,
			true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			newLine, modified := Replace(tc.line, tc.oldPattern, tc.newPattern)
			if modified != tc.modified {
				t.Fatal("line not modified")
			}
			if modified {
				if newLine != tc.want {
					t.Errorf("expected %q, got %q", tc.want, newLine)
				}
			}
		})
	}
}
