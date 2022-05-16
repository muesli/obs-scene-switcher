# obs-scene-switcher

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/obs-scene-switcher)
[![Go ReportCard](http://goreportcard.com/badge/muesli/obs-scene-switcher)](http://goreportcard.com/report/muesli/obs-scene-switcher)

obs-scene-switcher is a command-line remote control for OBS. It requires the
[obs-websocket](https://github.com/Palakis/obs-websocket) plugin to be installed on your system.

## Installation

### Packages & Binaries

On Arch Linux you can simply install the package from the AUR:

    yay -S obs-scene-switcher

Or download a binary from the [releases](https://github.com/muesli/obs-scene-switcher/releases)
page. Linux (including ARM) binaries are available, as well as Debian and RPM
packages.

### Build From Source

Alternatively you can also build `obs-scene-switcher` from source. Make sure you
have a working Go environment (Go 1.12 or higher is required). See the
[install instructions](http://golang.org/doc/install.html).

To build obs-scene-switcher, simply run:

    go get github.com/muesli/obs-scene-switcher

## Configuration

Edit scenes.toml and define which scenes you want to be connected to which
windows, e.g.:

```
[[scenes]]
    scene_name = "IDE"
    window_class = "code-oss"

[[scenes]]
    scene_name = "Terminal"
    window_name = "Konsole"

[[scenes]]
    scene_name = "Browser"
    window_class = "Chromium"

[[away_scenes]]
    scene_name = "Be Right Back"
```

In plain english, this means that whenever you focus your `VS Code` window, OBS
will be asked to switch to the scene called `IDE`. If you focus your `Konsole`
window it switches to the scene `Terminal`, and so on.

The `away_scenes` define scenes which, when currently active, temporarily stop
automatic scene switching. This is useful for keeping special scenes active,
like a "Be Right Back" mode you only want to manually disable again.

## Usage

Start obs-scene-switcher:

```bash
obs-scene-switcher -config scenes.toml
```
