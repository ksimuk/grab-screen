# grab-screen

`grab-screen` is a lightweight command-line utility written in Go that captures screenshots using the [XDG Desktop Portal](https://flatpak.github.io/xdg-desktop-portal/) API. It is designed to be easily integrated into scripts and window manager configurations.

## Features

- Uses the standard XDG Desktop Portal `Screenshot` interface (works on Wayland and X11 with appropriate portal backends).
- Executes a specified command with the screenshot path as an argument.
- Automatically cleans up the temporary screenshot file after execution (unless `--keep` is used).
- Prints the screenshot path to stdout if no command is provided.

## Prerequisites

- A Linux environment.
- A running XDG Desktop Portal implementation (e.g., `xdg-desktop-portal-gnome`, `xdg-desktop-portal-wlr`, `xdg-desktop-portal-kde`, or `xdg-desktop-portal-hyprland`).
- DBus.

## Installation

### Build from source

```bash
git clone https://github.com/ksimuk/grab-screen.git
cd grab-screen
go build .
```

## Usage

### Basic Usage

Run `grab-screen` without arguments to take a screenshot. The path to the saved image will be printed to stdout, and the file will be immediately deleted (unless you handle it quickly or use `--keep`).

```bash
./grab-screen
# Output: /home/user/Pictures/Screenshots/Screenshot from 2023-10-27 10-00-00.png
```

### Execute a Command

You can pass a command to `grab-screen`. The screenshot path will be appended as the last argument to the command.

**Example: Open the screenshot in an image viewer**

```bash
./grab-screen feh
```

**Example: Copy to clipboard (using wl-copy)**

```bash
./grab-screen wl-copy
```

**Example: Edit with Swappy**

```bash
./grab-screen swappy -f
```

**Example: Upload to a server (fictional script)**

```bash
./grab-screen ./upload-script.sh
```

### Keep the File

By default, `grab-screen` deletes the temporary screenshot file after the command finishes executing. To keep the file, use the `-k` or `--keep` flag.

```bash
./grab-screen --keep echo "Saved to:"
```

## License

[MIT](LICENSE)
