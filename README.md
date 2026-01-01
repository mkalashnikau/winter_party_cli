# Winter Party CLI

Winter Party CLI renders a festive snowfall animation right in your terminal using the Go `tcell` library. Snowflakes drift across the screen at 60 FPS, and you can trigger bursts of flakes for extra sparkle.

## Prerequisites

- Go 1.21 or newer
- A terminal that supports truecolor (most modern terminals do)

## Install Dependencies

The project relies on Go modules. Fetch dependencies with:

```bash
go mod tidy
```

## Run Locally

You can run the animation without building an executable:

```bash
go run ./...
```

Alternatively, build and run a binary:

```bash
go build -o winter-party
./winter-party
```

## Controls

- `Esc` — spawn an extra burst of snowflakes
- `q` or `Ctrl+C` — exit the program

Enjoy the terminal snowfall!
