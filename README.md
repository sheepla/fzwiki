
<div align="right">
    <img src="https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square"/>
    <img src="https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square"/>
</div>

<div align="center"><h1>fzwiki</h1></div>

<div align="center">

A command with fzf-like UI to quickly search Wikipedia articles and open it in your browser.

<img src="https://user-images.githubusercontent.com/62412884/147824124-b3e26a5a-752a-4714-9ceb-f552a0648ebb.png" />
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

### Binary

TODO

## LICENSE

[MIT](./LICENSE)

## Contributing

Welcome! 💕
