
<div align="right">
    <img src="https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square"/>
    <img src="https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square"/>
</div>

<div align="center"><h1>fzwiki</h1></div>

<div align="center">

A command line tool with fzf-like UI to search Wikipedia articles and open it in your browser quickly.

<img src="https://user-images.githubusercontent.com/62412884/148137551-4d2523e6-3292-48bf-896a-52d09f9d0a3e.png" />
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

### Download Executable Binary

You can download executable binaries from the release page.

> [Latest Release](https://github.com/sheepla/fzwiki/releases/latest)

**NOTE**:

With tools like [ghg](https://github.com/songmu/ghg), you can easily install executable from GitHub release and update version.

```bash
ghg get sheepla/fzwiki  # Install
ls -l $(ghg bin)/fzwiki # It will exists executable
```

## Configuration

You can change the default language for Wikipedia by setting a value in the environment variable `FZWIKI_LANG` .

```bash
FZWIKI_LANG="ja" fzwiki ... # --> search from ja.wikipedia.org instead of en.wikipedia.org
```

If you want to make the setting persistent,
add the following line to your rc file of the shell.


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

Welcome! 💕
