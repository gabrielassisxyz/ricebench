package fixture

import "fmt"

// ResolvedColors returns the resolved sRGB value for every role in the palette,
// following aliases to their authored target. A role that resolves to a target with no
// authored value is an error rather than a silent empty color.
func ResolvedColors(palette Palette) (map[RoleID]string, error) {
	roles := paletteRoleMap(palette)
	resolved := make(map[RoleID]string, len(roles))
	for id := range roles {
		target, err := resolveRole(id, roles, nil)
		if err != nil {
			return nil, err
		}
		targetRole := roles[target]
		if targetRole.Value == nil {
			return nil, fmt.Errorf("role %q resolves to %q which has no authored value", id, target)
		}
		resolved[id] = targetRole.Value.SRGB
	}
	return resolved, nil
}
