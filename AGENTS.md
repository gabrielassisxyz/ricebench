# RiceBench Agent Briefing

> Read before every interaction. Living spec: short, imperative. On every gotcha or decision, append one line here.

> **What it is:** a local visual preference elicitation tool that discovers a user's desktop color preferences through contextual comparisons and turns them into a reusable palette for Linux ricing. **Calibration:** Tier 2 · Phase: work. The stakes are personal rather than external: the tool runs locally, serves one user and holds nobody else's data, but it is public from its first commit and is built to grow. Update the phase as the project moves work, then right, then fast; an agent reads this line to decide how much rigor a change deserves. **Review gate:** standard. One independent external opinion over the whole branch diff, exactly once, pre-push. No per-commit or mid-development reviews. A review round whose findings are all non-behavioral and none blocking converges as approved; any blocking finding forces a new round. Dosage changes for one task are announced, never silent.

## Stack and commands

- **Stack:** Go standard `net/http` server, React and TypeScript frontend built by Vite, the built assets embedded with `go:embed` and served by the same process. No web framework, no global state or data-fetching library until a concrete need appears.
- **Run:** `go run ./cmd/ricebench`, then open the printed URL. Loopback by default.
- **Frontend dev server:** `cd web && npm run dev`
- **Build the embedded frontend:** `bin/build-web`
- **Build the binary:** `go build ./cmd/ricebench` (after `bin/build-web`, or it starts and reports that its assets are missing)
- **Test:** `go test ./...`
- **Full gate:** `bin/ci`
- **Hooks, once after clone:** `bin/install-hooks`
- **Worktrees:** `bin/worktree new <type>/<description>`. Never work in the main checkout.
- **Planning:** `br` for issues, `bv --robot-*` for graph triage; never run bare `bv`. The public planning surface is `ROADMAP.md`, and issue IDs never appear in it.

## Scope (current)

- **Current scope:** the elicitation instrument itself. Versioned contextual scenes, neutrally labelled candidates, triadic construct elicitation, strength-aware pairwise comparison, controlled one-variable variants, an evidence-linked preference profile, and a consumer-independent source palette. Don't expand beyond it without a present need; if a change drifts past it, STOP and flag it.
- **Out of scope on purpose:** desktop adapters of any kind, live desktop mutation, dotfile install or rollback, a universal ricing schema or plugin system, typography and layout taste, accounts, telemetry, sync, a hosted service, a theme library or marketplace. Adapting a finished palette to a real desktop happens outside this repository, and the first user running one particular desktop is not a reason to add an adapter here.

## Method invariants

The experiment is a measurement instrument. These are what make its output mean anything, and none of them is checkable by a gate.

- **Candidate labels stay neutral.** Participant-facing surfaces show `P1` through `P8` and nothing else. A theme name, a mood word or a generator parameter creates an association before the image is judged, and the judgment then measures the label.
- **Fixtures hold every non-color variable constant.** Content, typography, layout, spacing, icons, dimensions, wallpaper and application state are identical across candidates, and all candidates render from one scene definition rather than hand-maintained screenshots. Otherwise the experiment measures a bundle of changes.
- **Raw judgments are append-only and are never rewritten by an interpretation.** A correction is a new event. Derived files are reproducible from the event log and record which position of it they cover.
- **The user's own words survive.** A construct keeps the participant's terms; the technical reading is an attached, revisable layer, never a replacement.
- **Validity gates precede taste.** Contrast, semantic distinguishability, cursor, selection and muted-text checks filter candidates before preference chooses among the survivors. Passing a ratio does not prove comfort, and failing a required floor is a blocker.
- **Protocol counts are first-experiment guesses.** Eight candidates, four to six triads, roughly twelve to sixteen comparisons: these are pragmatic starting values, not constants established by the cited research. Record what the pilot measures and revise them from it.
- **Reporting failure is a valid outcome.** The instrument must be able to conclude that no candidate was acceptable, that the constructs were too ambiguous, or that judgments were unstable. It must never manufacture a palette to complete the flow.

<!-- BEGIN universal-principles v3 -->
## Working principles

- **The human defines the WHAT; the agent decides the HOW.** Don't wait for line-by-line dictation. Plan first for non-trivial tasks: show the plan + to-do list, wait for approval.
- **Think before coding — don't assume, don't hide confusion.** State assumptions explicitly; if multiple interpretations exist, present them — don't pick silently. If a simpler approach exists, say so and push back. If a task is impossible under the stated constraints, or info is missing, say so — don't guess. (For trivial tasks, use judgment; this is bias, not ritual.)
- **Surgical changes — touch only what you must.** Every changed line traces to the task. Don't "improve" adjacent code, reformat, or refactor what isn't broken; match existing style even if you'd do it differently. Flag unrelated dead code — don't delete it. Remove only the imports / variables / functions your own change orphaned.
- **Chesterton's Fence — find the problem before undoing the decision.** A config, a flag, a workaround that looks arbitrary is a **fence**: someone put it there, probably to fix something that is invisible to you *because the fence is working*. You arrive with no history, so absence of a visible reason is evidence of your ignorance, not of its uselessness. When your fresh measurement contradicts what the human vaguely remembers ("I changed this once, because of some problem"), **your measurement is the suspect first** — it may be measuring the case that *isn't* failing. Go find the original problem, then decide. *(A CIFS share was benchmarked with a big sequential `dd`, looked fast, and the local-disk download dir was "fixed" away — while the actual failure was random writes: par2, unrar, torrent piece-writes. Two wrong commits.)*
- **Goal-driven execution — define the success check, then loop to it.** Turn the task into something verifiable before coding: "add validation" → write tests for invalid inputs, then pass them; "fix the bug" → write a failing repro test, then pass it; "refactor X" → tests green before and after. For multi-step work, state a brief plan with a verify step each.
- **"Flaky" is not a diagnosis — test in the environment the thing actually runs in.** A component that fails *consistently* under automation is being **mis-invoked**, not being unreliable; "it works when I run it by hand" is not evidence that it works. The shell you test in has a TTY, a `$HOME`, an `ssh-agent`, an interactive stdin — the systemd unit, the CI job and the scripted harness have none of those, so a passing manual run can be testing a different program. Reproduce it *there* (start the unit, `env -u SSH_AUTH_SOCK`, `</dev/null`, `--dry-run` to print the real command line) before accepting "unstable" as a cause. **When a fix doesn't change the symptom, stop fixing and go look at what is actually being executed.** *(An interactive-mode flag with no TTY made one harness fail every review panel for weeks, written off as "flaky"; it was the wrong flag.)*
- **KISS — don't solve a problem you don't have yet.** Simplicity isn't "write less code"; it's not building for a need that doesn't exist. Let structure emerge from the code.
- **YAGNI & flat.** No preventive abstractions, no single-use interfaces. Interfaces for real boundaries only. Architecture is *extracted* once a pattern proves itself in real use — never designed up front for a user who doesn't exist yet. Need pulls architecture.
- **Development cost is not your cost — don't let it pick the design.** Choosing between technical options, weight quality, simplicity, robustness and long-term maintainability; don't weight how long the work takes. The estimate comes out in human units — days, weeks — because that is what the training data measured, and the cheaper option then wins on a cost the agent does not pay. This is **not** licence to over-build: KISS and YAGNI decide *whether* a thing is needed, and this decides *how well* it is built once it is. "That would take a week" is not an argument here; "nothing needs this yet" is.
- **Order: make it work → make it right → make it fast** (Kent Beck), in that order. Most over-engineering is doing "right"/"fast" before a working thing exists to justify it.
- **Flag scope creep — a standing duty, not a suggestion.** When a solo tool starts being framed as a public / multi-user / multi-tenant / plugin-system / configurable-N-backends platform before a real, present need exists, STOP and ask: "Is this needed now?" Justify future-proofing against a need that exists *today*.
- **No silent decisions (comprehension debt).** Never make a silent architectural or design call — state it and record the rationale, so the reasoning is recoverable later.
- **Real decisions are presented in the chat, in isolation — never via popup.** When a design/architecture/scope/trade-off decision arises, surface it on its own: the options, what each means, pros/cons/trade-offs, and a recommendation — then decide together. Don't bury it mid-text or bundle it with other topics, and don't compress it into a quick-pick widget (e.g. AskUserQuestion) — the widget skips the reasoning and overlays the explanation. Widgets are for trivial short-answer picks only.
- **Long answers are written to be scanned, not read twice.** For recaps, status reports, batch reviews, plans, and any comparison of options: lead with the outcome in one line, then break the body into bullets and **bold** the load-bearing terms. Options are always a list — one bullet per option, the recommended one marked — never a paragraph the reader has to parse to find the choices. Reserve unbroken prose for short arguments; a wall of paragraphs costs more in re-reading than the structure would have cost in words.

## Git: branches, commits, PRs, comments

- **Ask the repo for its default branch; never assume one.** Repos differ — `master` and `main` are both common, often in the same person's account — and a wrong guess sends a PR to a branch that does not exist, or, worse, has you "fixing" a URL that was right all along. `git symbolic-ref --short refs/remotes/origin/HEAD | sed 's|^origin/||'`, or `gh repo view --json defaultBranchRef -q .defaultBranchRef.name`. Never commit directly to it: branch, then PR.
- **A new repo starts on `main`.** That is the preferred name, and `init.defaultBranch` is set to it, so `git init` produces it without anyone choosing. It settles new repos only: an existing one keeps the branch it has, because renaming breaks open PRs, CI filters, deploy hooks and every permalink into the tree, and buys nothing. The rule above still governs everything already in existence — ask, never assume.
- **Branches** — Conventional Branch (conventionalbranch.org): `<type>/<kebab-description>`, types `feature/`, `bugfix/`, `hotfix/`, `chore/`, `release/`, `docs/`.
- **Commits** — Conventional Commits (conventionalcommits.org): `<type>(scope): <description>`, types `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `ci`, `build`, `perf`, `style`. Breaking change → `!` after the type or a `BREAKING CHANGE:` footer.
- **Atomic commits** — one logical change per commit, each independently green and revertible. Never `git add .` blind; split unrelated changes.
- **Always work in your own worktree — mandatory, not conditional.** Parallel sessions are opened freely and nothing signals their existence to you, so a "check whether another session is here first" step can never be reliable — the honest answer is always "maybe". The only collision-proof arrangement is structural: keep the main working tree on the default branch as a clean reference and **never work in it** — before your first write (commit, branch, rebase, stash; read-only exploration is exempt), create your own worktree and do everything there: `git worktree add ../<repo>-<task> -b <your-branch> <origin>/<default-branch>`. Do this **whether or not** you believe another agent is running — that belief is exactly what you cannot verify. Report which worktree/branch you used; remove it once merged. Only the human can see all the open sessions.
- **Pull requests** — describe **what + why**. *What*: a 1–3 line summary. *Why* (the bulk): decisions, trade-offs, rejected alternatives. The diff shows the what; the PR explains why.
- **Comments** — always **WHY, not WHAT**: explain intent, never restate the obvious mechanics. Keep existing comments; they carry intent.

## Code style (baseline)

- Functions: 4–40 lines, one thing each (SRP). Files: under ~500 lines, split by responsibility.
- Names specific and unique — avoid `data`, `handler`, `Manager`, `util`.
- Explicit types. Early returns over nested ifs; max ~2 levels of indentation.
- Inject dependencies; wrap third-party libs behind a thin interface this project owns.
- No duplication — but don't extract *too early*. Tolerate duplication while the pattern is still forming; extract the abstraction *from* proven, repeated code, never ahead of it.
- **Refactoring is not automatic.** After a large feature, list refactoring candidates (files > ~500 lines, duplicated logic, long functions, hardcoded config) and ask before pruning — the human decides, the tests are the safety net. Consolidate when the thing works and the seams are obvious, not before.
<!-- END universal-principles v3 -->

## Tests (TDD)

- Every feature is born with a test; every bugfix with a regression test.
- Break every new guard on purpose once and watch it fail. A test that cannot fail is worse than no test, because it reports success.
- Tests run with ONE command, no manual setup, no secret credential. If it cannot run headless, it is wrong.
- Go tests are hermetic: no real network, no host filesystem, no exec. Mock at boundaries through interfaces this project owns.
- Before saying "done", run `bin/ci` and report the result.

## Small releases

- Every commit on `main` passes `bin/ci` and is releasable. No broken commit fixed by the next one.
- Closed work is committed before switching tasks; flag it if it has not been.

## Release

- A `v*` tag builds the Linux `amd64` binary through goreleaser, emits `checksums.txt` and publishes a GitHub Release. The frontend build runs first, in the goreleaser `before` hook.
- The release artifact is one executable that runs with no Go, Node.js or npm installed. It installs no service and changes no firewall rule.
- `install.sh` fetches the archive for the current release into `~/.local/bin`.

## Security (habit, not a phase)

- The listener is loopback by default. A non-loopback bind is an explicit opt-in that must print, at startup, that the session is unauthenticated and readable by anything that can reach the port. Never make that quieter.
- Never add automatic firewall changes, and never write outside the configured data directory.
- Raw rationales are free text and can carry personal information. They stay local. Nothing is published automatically.
- When touching user input, network, filesystem or path construction, flag the risk and propose the guard.
- Dependency CVEs are caught by `govulncheck` and `npm audit` in `bin/ci`.

## Prose

- No em-dash. Use a comma, a colon, a semicolon or a full stop. This is checked by `bin/ci`, and it applies to Markdown, source comments, config, commit messages and PR text alike.
- Markdown is soft-wrapped: one paragraph, one line, and the editor wraps it for display. Hard wrapping is a per-file accident that spreads by contact, since every editor and every agent inherits the wrap of whatever it is editing, and it makes a one-word edit show up as a reflowed paragraph in review. `scripts/md-unwrap.py --check` in `bin/ci` enforces it, and `--write` sweeps a file that arrived wrapped.
- Bold marks structure, a bullet lead-in or a table header, never emphasis mid-sentence. Same for italics: a term being introduced, not a word being stressed.
- No process narration anywhere a stranger can read it: no task ids, no phase names, no review rounds, no mention of who or what reviewed a diff, no reference to a session or a conversation. Commit and PR text describe the problem and the change, never how the work was organised.
- No audience in the text. A README says what the software does, not who is going to read it.
- Comment density is low by default: the non-obvious only, the why and not the what. Long reasoning belongs in an ADR, not in a header comment.

## Git and secrets

- Before any commit, show `git status` and `git diff --cached`; confirm no secret is staged. If you spot one, STOP and report it. The gitleaks pre-commit hook is the deterministic backstop; this habit is the probabilistic one.
- Real secrets stay out of git. Only an `.env.example` with fake values is committed.
- Experiment data lives outside the repository and is never committed, not even as a fixture.

## Post-implementation checklist (run before "done")

1. Commits small and well-described.
2. Refactoring candidates listed, if the change was large.
3. Security risks flagged, if you touched a sensitive surface.
4. This spec updated if behavior, setup or release flow changed, and any hurdle it gained is classified rather than just appended.

## Common hurdles

| hurdle | class | gate |
|---|---|---|
| Vite runs with `emptyOutDir` off, because that directory carries the tracked `.gitkeep` that `go:embed` needs on a clone with no frontend build. A bare `npm run build` therefore leaves the previous hashed assets behind and they end up inside the binary. | tripwire | none; use `bin/build-web`, which clears the output first |

**A hurdle promoted to a gate is deleted from this table, not duplicated.** The gate is the instruction; a line here restating it only dilutes the ones still unguarded.
