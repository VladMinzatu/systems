# tui-test

A minimal [Bubble Tea](https://github.com/charmbracelet/bubbletea) (v2) TUI: type a line in the bottom panel, press Enter, and it gets appended to the top panel.

| Key            | Action                  |
| -------------- | ----------------------- |
| `enter`        | append the current line |
| `pgup`/`pgdown`| scroll the output panel |
| `esc`/`ctrl+c` | quit                    |


## Run in Docker

```sh
make docker-build
make docker-run
```

`docker-run` uses `docker run --rm -it`: `-i` keeps stdin attached so keystrokes reach the app, and `-t` allocates a pseudo-TTY so Bubble Tea can enter raw mode, use the alternate screen and receive window-resize events. Without `-t` the TUI won't render properly. `TERM` and `COLORTERM` are forwarded from the host so colors match your terminal.

The image is a multi-stage build: a static binary compiled in `golang:1.27.1-alpine`, shipped on `gcr.io/distroless/static-debian12:nonroot`.

Override the image name or tag with `make docker-build IMAGE=foo TAG=dev`.

## CPU panel

The right-hand panel graphs CPU usage once per second using [ntcharts](https://github.com/NimbleMarkets/ntcharts). Usage is computed from the aggregate `cpu` line of `/proc/stat`, so the app expects to run on Linux (i.e. in the container). `/proc/stat` is not namespaced, so inside Docker it reports the whole host (on macOS, the Docker Desktop VM), not just this container.
