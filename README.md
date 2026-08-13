# RiceBench

A visual workbench for Linux ricing.

RiceBench is a local visual preference elicitation tool. It discovers desktop color preferences through contextual comparisons and turns them into a reusable palette for Linux ricing.

## Why it exists

Picking a desktop palette usually happens one of two ways: adopt a popular named theme and live with its assumptions, or read hexadecimal values and swatches outside the applications where those colors will occupy hours of attention. Both invite the same loop. Choose something that looks nearly right, discover fatigue or invisible states in real use, replace it, repeat.

RiceBench changes the order. Candidates are seen in identical representative scenes first, the comparisons are expressed in the user's own words, the technical properties behind those judgments are isolated afterwards, and the result is one palette carrying the evidence that produced it.

## What it produces

One consumer-independent source palette, in the sense that Nord's central palette is one thing that many applications later adapt. It carries semantic roles for surfaces, text hierarchy, accent, focus, cursor, selection, semantic states and diffs, plus an ANSI terminal extension, each role in OKLCH for manipulation and resolved sRGB for compatibility. Alongside it, the raw judgments and an evidence-linked preference profile, so a decision can be traced, audited and revised rather than taken on faith.

It also produces a valid negative result. If no candidate is acceptable, if the elicited constructs are too ambiguous to map, or if the comparisons are unstable across repeats, the experiment reports that instead of manufacturing a palette.

## What it does not do

RiceBench does not change a desktop. It ships no adapter for any desktop environment, terminal, editor or browser, and it installs, backs up or rolls back no configuration. Adapting a finished palette to a real environment happens outside this repository, deliberately: the point of the first release is to find out whether contextual elicitation produces something worth adapting.

## Status

Early. The repository carries the server shell, the embedded asset pipeline and the release path; the elicitation instrument itself is not implemented. `ROADMAP.md` records what exists, what is missing and what is out of scope on purpose.

## Running from source

Requires Go and Node.js.

```sh
bin/build-web                 # build the frontend into the directory the binary embeds
go run ./cmd/ricebench        # serve it, loopback only, on the printed URL
```

For frontend work, `cd web && npm run dev` runs Vite with hot reload.

The server binds to `127.0.0.1:7391` by default. `--addr` accepts any bind address, including a LAN-reachable one for judging from another device. That mode is unauthenticated: every client that can reach the port can read the experiment and submit judgments, the server says so at startup, and it never changes a firewall rule.

## Development

```sh
bin/install-hooks    # once after cloning: versioned git hooks, including the secret scan
bin/ci               # the full gate, and the exact thing CI runs
```

`bin/worktree new <type>/<description>` creates an isolated worktree and branch. `AGENTS.md` carries the working conventions for both people and agents.

## License

MIT. See `LICENSE`.
