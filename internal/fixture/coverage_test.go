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

func TestRoleCoverageGroupsByRoleAndResolvesAliases(t *testing.T) {
	byRole := RoleCoverage(minimalFixture())

	if got := len(byRole[RoleTerminalBackground]); got != 1 {
		t.Fatalf("terminal-background references = %d, want 1", got)
	}
	if got := len(byRole[RoleBackground]); got != 1 {
		t.Fatalf("background references = %d, want 1 (resolved from terminal-background)", got)
	}
	if byRole[RoleBackground][0].Path != "scenes[0].regions[0].background" {
		t.Fatalf("background reference path = %q", byRole[RoleBackground][0].Path)
	}

	success := byRole[RoleSuccess]
	if len(success) != 1 || success[0].Level != LevelBlock {
		t.Fatalf("success references = %+v, want one block-level reference", success)
	}
}

func TestCoverageViewListsEveryRequiredRoleInOrder(t *testing.T) {
	entries := CoverageView(allFixtureSets()...)

	if len(entries) != len(requiredCoverageRoles()) {
		t.Fatalf("coverage view has %d entries, want %d", len(entries), len(requiredCoverageRoles()))
	}
	for index, entry := range entries {
		if entry.Role != requiredCoverageRoles()[index] {
			t.Fatalf("entry %d role = %q, want %q", index, entry.Role, requiredCoverageRoles()[index])
		}
		if len(entry.References) == 0 {
			t.Fatalf("required role %q has no references in the coverage view", entry.Role)
		}
	}
}
