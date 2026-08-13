# Roadmap

What exists, what is missing, and what is deliberately out of scope. Direction rather than schedule: evidence from each step decides whether the next one is worth taking.

## What exists

- A Go server that serves the embedded frontend on a configurable bind address, loopback by default, with an explicit startup warning when the listener is exposed.
- The frontend build pipeline: React and TypeScript through Vite, embedded into the binary.
- The deterministic gates: secret scan on commit, `bin/ci` running format, vet, lint, tests and dependency audits, the same script running in CI.
- The release path: a tagged build producing one Linux `amd64` executable with checksums.

## What is missing

The whole instrument. In dependency order:

1. **Fixtures.** Versioned scene definitions for four families: terminal and agent TUI, code and diff, desktop shell, reading and monitoring. Every candidate renders from one structural definition, so a comparison can only differ in color. Automated checks prove the non-color output is identical across candidates, and deliberately breaking a role makes a fixture fail visibly.
2. **Candidate generation.** Coherent seed palettes transformed deterministically along documented axes: polarity, neutral temperature, background lightness, foreground contrast, chroma, accent strategy, surface separation. A pool larger than the eight shown, filtered by the validity gates, selected for axis coverage rather than for an aesthetic score. Seeds, parameters and every curatorial intervention are preserved; only neutral IDs are shown.
3. **Validity gates.** Contrast floors, primary against secondary against muted text, cursor and selection visibility, focused against selected states, semantic states under color-vision deficiency simulation, terminal normal and bright pairs. Numeric checks, simulations and human observations reported as distinct kinds of evidence.
4. **The session and event record.** An append-only event log as the source of truth, a manifest fixing identity and setup versions, a single serialized writer, a lock excluding a second writer, and startup validation that reports a truncated record instead of accepting it.
5. **The elicitation flows.** Broad contextual review, triadic construct elicitation in the participant's own words, strength-aware pairwise comparison with rationales and reversed repeats, and controlled variants that move one property at a time.
6. **Derivation.** A transparent weighted method producing an evidence-linked preference profile, surfacing conflicts and low confidence rather than smoothing them, deterministic against a fixed record, with corrections appended as events.
7. **Source palette export.** The semantic core and ANSI terminal extension, OKLCH and resolved sRGB per role, provenance and validity metadata attached, aliases naming their target rather than copying a value.
8. **Real-use observation import.** Somewhere to put what an external adaptation learns, so a palette can be revised with the evidence that revised it.

## Deliberately out of scope

Not forgotten, and not a backlog:

- Adapters for any desktop environment, window manager, bar, launcher, terminal, editor, browser or notes application.
- Changing a live desktop: installation, dotfile management, backup, rollback, synchronization.
- A universal ricing schema, a plugin architecture or a theme interoperability standard.
- Typography, spacing, iconography, layout, motion, or general design taste.
- Accounts, telemetry, cloud sync, collaboration or any hosted service.
- A public theme library, marketplace or community feature.
- Claims about emotion or personality inferred from color choices.
- A complex preference-ranking model before a simple transparent one is shown insufficient.

Adding any of these before the first experiment succeeds would obscure whether contextual preference elicitation creates value at all, which is the only question the first release is built to answer.

## Further out

If the instrument proves itself, the direction is a visual workbench for Linux ricing: previewing and composing a desktop beyond color alone, then stable output contracts that separately maintained adapters can target, then a library where reusable palettes carry their source description and adaptation provenance instead of only a screenshot. Each step has to be earned by the one before it.
