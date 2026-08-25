package fixture

// DesktopShellFixtureSet returns the desktop-shell scene family. These scenes model the
// non-terminal system surfaces where focus, urgency, overlays and adjacent backgrounds
// fail differently from text surfaces, in this project's own vocabulary rather than any
// consumer configuration format.
func DesktopShellFixtureSet() FixtureSet {
	return FixtureSet{
		SchemaVersion: FixtureSchemaVersion,
		ID:            "fixture-desktop-shell",
		Scenes: []Scene{
			desktopShellWorkspaceScene(),
			desktopShellLauncherScene(),
			desktopShellNotificationScene(),
		},
	}
}

func desktopShellWorkspaceScene() Scene {
	return Scene{
		ID:     "desktop-shell-workspace",
		Family: FamilyDesktopShell,
		Regions: []Region{
			{
				ID:         "wallpaper",
				Kind:       RegionSurface,
				State:      StateDefault,
				Background: RoleBackground,
			},
			{
				ID:         "status-bar",
				Kind:       RegionStatus,
				State:      StateDefault,
				Background: RoleSurface,
				Foreground: RoleForeground,
				Blocks: []ContentBlock{
					{
						ID:         "workspace-active",
						Kind:       ContentStatusItem,
						State:      StateActive,
						Text:       "1",
						Background: RoleAccent,
						Foreground: RoleAccentForeground,
					},
					{
						ID:         "workspace-inactive",
						Kind:       ContentStatusItem,
						State:      StateInactive,
						Text:       "2",
						Background: RoleSurface,
						Foreground: RoleTextMuted,
					},
					{
						ID:         "module-primary",
						Kind:       ContentStatusItem,
						State:      StateDefault,
						Text:       "primary",
						Foreground: RoleForeground,
					},
					{
						ID:         "module-muted",
						Kind:       ContentStatusItem,
						State:      StateMuted,
						Text:       "muted",
						Foreground: RoleTextMuted,
					},
				},
			},
			{
				ID:         "window-focused",
				Kind:       RegionFrame,
				State:      StateFocused,
				Background: RoleSurfaceElevated,
				Foreground: RoleForeground,
				Border:     RoleFocus,
				Blocks: []ContentBlock{
					{
						ID:     "focus-ring",
						Kind:   ContentFocusRing,
						State:  StateFocused,
						Border: RoleFocus,
					},
					{
						ID:         "window-focused-title",
						Kind:       ContentText,
						State:      StateFocused,
						Text:       "Focused window",
						Foreground: RoleForeground,
					},
				},
			},
			{
				ID:         "window-unfocused",
				Kind:       RegionFrame,
				State:      StateInactive,
				Background: RoleSurface,
				Foreground: RoleTextSecondary,
				Border:     RoleSurface,
				Blocks: []ContentBlock{
					{
						ID:         "window-unfocused-title",
						Kind:       ContentText,
						State:      StateInactive,
						Text:       "Unfocused window",
						Foreground: RoleTextSecondary,
					},
				},
			},
		},
	}
}

func desktopShellLauncherScene() Scene {
	return Scene{
		ID:     "desktop-shell-launcher",
		Family: FamilyDesktopShell,
		Regions: []Region{
			{
				ID:         "launcher-overlay",
				Kind:       RegionOverlay,
				State:      StateDefault,
				Background: RoleSurfaceElevated,
				Foreground: RoleForeground,
				Border:     RoleSurface,
				Blocks: []ContentBlock{
					{
						ID:         "launcher-query",
						Kind:       ContentText,
						State:      StateDefault,
						Text:       "query",
						Foreground: RoleForeground,
					},
				},
			},
			{
				ID:         "launcher-results",
				Kind:       RegionList,
				State:      StateDefault,
				Background: RoleSurfaceElevated,
				Foreground: RoleForeground,
				Blocks: []ContentBlock{
					{
						ID:         "result-neighbor-above",
						Kind:       ContentListItem,
						State:      StateDefault,
						Text:       "Neighboring result",
						Foreground: RoleTextSecondary,
					},
					{
						ID:         "result-selected",
						Kind:       ContentListItem,
						State:      StateSelected,
						Text:       "Selected result",
						Background: RoleSelectionBackground,
						Foreground: RoleSelectionForeground,
					},
					{
						ID:         "result-secondary-metadata",
						Kind:       ContentText,
						State:      StateMuted,
						Text:       "Secondary metadata",
						Foreground: RoleTextMuted,
					},
					{
						ID:         "result-neighbor-below",
						Kind:       ContentListItem,
						State:      StateDefault,
						Text:       "Another result",
						Foreground: RoleTextSecondary,
					},
				},
			},
		},
	}
}

func desktopShellNotificationScene() Scene {
	return Scene{
		ID:     "desktop-shell-notification",
		Family: FamilyDesktopShell,
		Regions: []Region{
			{
				ID:         "notification",
				Kind:       RegionSurface,
				State:      StateDefault,
				Background: RoleSurfaceElevated,
				Foreground: RoleForeground,
				Border:     RoleSurface,
				Blocks: []ContentBlock{
					{
						ID:         "notification-title",
						Kind:       ContentText,
						State:      StateDefault,
						Text:       "Notification",
						Foreground: RoleForeground,
					},
					{
						ID:         "notification-body",
						Kind:       ContentText,
						State:      StateDefault,
						Text:       "Body",
						Foreground: RoleTextSecondary,
					},
				},
			},
			{
				ID:         "notification-urgent",
				Kind:       RegionSurface,
				State:      StateUrgent,
				Background: RoleSurfaceElevated,
				Foreground: RoleForeground,
				Border:     RoleWarning,
				Blocks: []ContentBlock{
					{
						ID:         "notification-urgent-title",
						Kind:       ContentText,
						State:      StateUrgent,
						Text:       "Urgent",
						Foreground: RoleWarning,
					},
				},
			},
			{
				ID:         "transient-overlay",
				Kind:       RegionOverlay,
				State:      StateDefault,
				Background: RoleSurfaceElevated,
				Foreground: RoleForeground,
				Border:     RoleSurface,
				Blocks: []ContentBlock{
					{
						ID:         "overlay-content",
						Kind:       ContentText,
						State:      StateDefault,
						Text:       "Overlay",
						Foreground: RoleForeground,
					},
				},
			},
		},
	}
}
