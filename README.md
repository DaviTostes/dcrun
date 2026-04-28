# dcrun — Don't Care, Run

You have a project. You want it running. You don't remember if it's `npm run dev`, `yarn dev`, `pnpm dev`, `bun dev`, `deno task dev`, `cargo run`, `go run .`, or one of the other 47 incantations the JavaScript ecosystem invented last Tuesday.

`dcrun` doesn't care. Neither should you.

```bash
cd ~/some-project
dcrun
```

That's the whole interface. You're welcome.

## Install

Pick one. They both work. Probably.

```bash
go install github.com/davitostes/dcrun@latest
```

Or, for people who enjoy typing:

```bash
git clone https://github.com/davitostes/dcrun
cd dcrun
go build -o dcrun .
mv dcrun ~/.local/bin/
```

## How it works

It looks at your project files. It guesses. It's usually right.

If `bun.lockb` exists, it's Bun. If `Cargo.toml` exists, it's Rust. If `manage.py` exists, congratulations, you're maintaining a Django app. We're sorry.

## Supported

| Stack | Detected by |
|---|---|
| Bun | `bun.lock` / `bun.lockb` |
| pnpm | `pnpm-lock.yaml` |
| Yarn | `yarn.lock` |
| Deno | `deno.json` / `deno.lock` |
| Node (npm) | `package.json` |
| Go | `go.mod` |
| Rust | `Cargo.toml` |
| Python (uv) | `uv.lock` |
| Python (poetry) | `poetry.lock` |
| Django | `manage.py` |
| Dotnet | `*.csproj` / `*.sln` / `*.fsproj` |
| Lua | `main.lua` / `init.lua` / `*.lua` |
| Make | `Makefile` with a `dev:` target |

Lockfiles win over `package.json`, so the right JS tool gets picked. Make is the last-resort escape hatch — write a `dev` target, dcrun will find it.

## Your stack isn't here?

Open a PR. Or an issue. Or suffer.

## License

[MIT](LICENSE) — do whatever.
