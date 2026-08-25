package fixture

import (
	"strings"
	"testing"
)

func TestCoverageTraversalProducesFullPaths(t *testing.T) {
	references := Coverage(minimalFixture())
	want := []RoleReference{
		{Role: RoleTerminalBackground, Path: "scenes[0].regions[0].background", Level: LevelRegion},
		{Role: RoleTerminalForeground, Path: "scenes[0].regions[0].foreground", Level: LevelRegion},
		{Role: RoleSuccess, Path: "scenes[0].regions[0].blocks[0].foreground", Level: LevelBlock},
	}
	if len(references) != len(want) {
		t.Fatalf("got %d references, want %d: %v", len(references), len(want), references)
	}
	for index := range want {
		if references[index] != want[index] {
			t.Fatalf("reference %d = %+v, want %+v", index, references[index], want[index])
		}
	}
}

func TestRequiredRoleCoverageAcrossAllFamilies(t *testing.T) {
	if err := RequiredRoleCoverage(allFixtureSets()...); err != nil {
		t.Fatalf("combined coverage: %v", err)
	}
}

func TestRequiredRoleCoverageReportsMissingRole(t *testing.T) {
	sets := allFixtureSets()
	stripRunRole(&sets[0], RoleANSI0)

	err := RequiredRoleCoverage(sets...)
	if err == nil {
		t.Fatal("missing role was not reported")
	}
	if !strings.Contains(err.Error(), string(RoleANSI0)) {
		t.Fatalf("error %q does not name the missing role %q", err, RoleANSI0)
	}
}
