# Medal.TV Toolbox

This project started out as a very simple PowerShell Script that served the same purpose this project does but everything was hard coded into the script but I needed to all of my clips to different directories for easy of access for my [YouTube videos](https://youtube.com/@EndLordHD). So I decided to re-write in Python and now in Go because Python was slow and the project was very clutered. Here's the old repo for the Python version of this project: https://github.com/EndCod3r/medal-tools

## How to download & use:

### Git Clone:

```
git clone https://github.com/EndCod3r/medal-toolbox.git
cd medal-toolbox
```

or if you don't have Git installed

### Download Zip:

Go to the [releases](https://github.com/EndCod3r/medal-toolbox/releases/latest), and open the Assets drop down, and download the `Source code (zip).` Extract it once it's done.

## How to use:

Open a terminal in the medal-toolbox directory and run
`medal-cli.exe` or `medal-tui.exe` depending on which you want to use. Use the `-h` argument to see help options for each binary.

## How to build:

Ensure you have [Golang](https://go.dev/) installed.

Open a terminal in the medal-toolbox folder and run

```
go build -o medal-cli.exe ./cmd/cli
```

for the CLI or for the TUI run

```
go build -o medal-tui.exe ./cmd/tui
```
