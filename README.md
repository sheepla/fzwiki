<div align="right">

![CI](https://github.com/sheepla/fzwiki/actions/workflows/ci.yml/badge.svg)
![Relase](https://github.com/sheepla/fzwiki/actions/workflows/release.yml/badge.svg)

<a href="https://github.com/sheepla/fzwiki/releases/latest">

![Latest Release](https://img.shields.io/github/v/release/sheepla/fzwiki?style=flat-square)

</a>

</div>

<div align="center"><h1>fzwiki</h1></div>

<div align="center">

A command line tool with fzf-like UI to search Wikipedia articles and open it in your browser quickly.

![](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)

</div>

<div align="center">
<img src="https://user-images.githubusercontent.com/62412884/148652520-0dcfafa3-f7e4-4a3d-b9e7-ed93ae74bab4.gif" />
</div>

## Usage

1. Run the command by specifying a search query.
2. Press the `<Tab>` key to select item(s), then press `<Enter>` key to confirm the selection.
3. The pages URL of the selected items will output. If you specify the `-o`, `--open` flag, it will open the page in your default browser.

### Help message

```
fzwiki [OPTIONS] QUERY...

Application Options:
  -V, --version  Show version
  -o, --open     Open URL in your web browser
  -l, --lang=    Language for wikipedia.org such as "en", "ja", ...

Help Options:
  -h, --help     Show this help message
```

### Key bindings

|Key              |Description           |
|-----------------|----------------------|
|type some text   |narrow down candidates|
|`<C-j>` / `<C-n>`|move focus down       |
|`<C-k>` / `<C-p>`|move focus up         |
|`<Tab>`          |select the item       |
|`<Enter>`        |confirm the selection |

## Installation

### Build from source

```bash
git clone https://github.com/sheepla/fzwiki.git
cd fzwiki
go install
```

### Download executable binary

You can download executable binaries from the release page.

> [Latest Release](https://github.com/sheepla/fzwiki/releases/latest)

### Use GitHub release installer tools

These tools make it easy to install executable binaries from GitHub Release.

with [ghg](https://github.com/Songmu/ghg):

```bash
ghg get sheepla/fzwiki  # Install
ls -l $(ghg bin)/fzwiki # It will exists executable
```

with [relma](https://github.com/jiro4989/relma):

Copy download link URL from [Latest Release](https://github.com/sheepla/fzwiki/releases/latest) page, then run below.


```bash
relma init                           # Setup
relma install {{DOWNLOAD_LINK_URL}}  # Install
ls -l $(ghg bin)/fzwiki              # It will exists executable
```

with [gh-install](https://github.com/redraw/gh-install)

```bash
gh install sheepla/fzwiki # Install
ls -l ~/.local/bin/fzwiki # It will exists executable
```

## Configuration

To change the default language for Wikipedia, set a value in the environment variable `FZWIKI_LANG` .

```bash
FZWIKI_LANG="ja" fzwiki ... # --> search from ja.wikipedia.org instead of en.wikipedia.org
```

To make the setting persistent, add the following line to your rc file of the shell.


- **bash** (`~/.bashrc`) or **zsh** (`~/.zshrc`):

```bash
export FZWIKI_LANG="ja"
```

- **fish** (`~/.config/fish/config.fish`):

```fish
set -Ux FZWIKI_LANG ja
```

## LICENSE

[MIT](./LICENSE)

## Contributing

Welcome! ????

## Author

[Sheepla](https://github.com/sheepla)
