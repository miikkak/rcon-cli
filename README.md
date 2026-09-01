# rcon-cli

A CLI for attaching to an RCON-enabled game server, such as Minecraft.

This is an independent implementation, not a fork — inspired by
[itzg/rcon-cli](https://github.com/itzg/rcon-cli), same use case and CLI
shape, but built to correctly handle Minecraft's extended truecolor
formatting codes (`§x§R§R§G§G§B§B`, used by e.g. BlueMap) and built on the
actively maintained [gorcon/rcon](https://github.com/gorcon/rcon) client
library instead of the unmaintained `james4k/rcon`.

## About this project

This was built with heavy Claude Code assistance — most of the implementation
is AI-generated, with the design and review driven by me. It has unit test
coverage (see Testing below) and runs as the RCON client in my own production
Minecraft server stack ([mc-server-container](https://github.com/miikkak/mc-server-container)),
so it sees real day-to-day use, not just its own test suite. Read the source
and file issues if something looks off.

## Usage

Without any additional arguments, the CLI starts an interactive session
with the RCON server. Send the keyword `exit` to close the session.

```shell
rcon-cli --host mc1 --port 25575 --password secret
```

If arguments are passed, they're joined into a single command, sent to the
server, the response is printed, and the CLI exits.

```shell
rcon-cli --host mc1 --port 25575 --password secret stop
```

### Flags

| Flag         | Default     | Description            |
| ------------ | ----------- | ---------------------- |
| `--host`     | `localhost` | RCON server's hostname |
| `--port`     | `25575`     | RCON server's port     |
| `--password` | (none)      | RCON server's password |
| `--config`   | (none)      | Path to a config file  |

All flags can also be set via environment variables with an `RCON_` prefix
(e.g. `RCON_PORT`, `RCON_PASSWORD`), or via a `.rcon-cli.yaml` config file
in `$HOME`, `/data`, or `/server`.

## Color handling

Both legacy single-character Minecraft formatting codes (`§a`, `§c`, `§r`,
etc.) and truecolor codes (`§x§R§R§G§G§B§B`) are translated into ANSI
escape sequences for terminal display. Color translation is a no-op on
Windows.

## Building

```shell
make build
```

Produces a `rcon-cli` binary in the repo root. `make install` installs it
to `/usr/local/lib/helpers` (override with `PREFIX`/`DESTDIR`).

## Testing

```shell
go test ./... -race -cover
```
