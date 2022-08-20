<div align="center">
    <h1>
        <img width="750" src="https://raw.githubusercontent.com/everdrone/grab/main/.github/media/Dark@2x.png#gh-light-mode-only" alt="GRAB" />
        <img width="750" src="https://raw.githubusercontent.com/everdrone/grab/main/.github/media/Light@2x.png#gh-dark-mode-only" alt="GRAB" />
    </h1>
    <h3>Greedy, Regex-Aware Binary Downloader</h3>
</div>

<p align="center">
<a href="https://github.com/everdrone/grab/stargazers">
    <img src="https://img.shields.io/github/stars/everdrone/grab?color=8bd5ca&logo=github&logoColor=d9e0ee&labelColor=1e1d2f&style=for-the-badge" alt="Stargazers">
</a>
<!-- <img src="https://img.shields.io/static/v1?label=Reference&message=GO&color=7dc4e4&logoColor=d9e0ee&labelColor=1e1d2f&style=for-the-badge" alt="Go Package Reference"> -->
<a href="https://github.com/everdrone/grab">
    <img src="https://img.shields.io/tokei/lines/github/everdrone/grab?color=7dc4e4&logoColor=d9e0ee&labelColor=1e1d2f&style=for-the-badge&label=Lines&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxNiIgaGVpZ2h0PSIxNiI+PHBhdGggZmlsbD0iI2Q5ZTBlZSIgZD0iTTUgOC4yNWEuNzUuNzUgMCAwIDEgLjc1LS43NWg0YS43NS43NSAwIDAgMSAwIDEuNWgtNEEuNzUuNzUgMCAwIDEgNSA4LjI1ek00IDEwLjVBLjc1Ljc1IDAgMCAwIDQgMTJoNGEuNzUuNzUgMCAwIDAgMC0xLjVINHoiLz48cGF0aCBmaWxsPSIjZDllMGVlIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiIGQ9Ik0xMyAwSDNhMyAzIDAgMCAwLTMgM2MwIC42Ny4yMiAxLjI1LjYgMS43Mi41My42NSAxLjMzLjc4IDEuOTEuNzhoMS4xOEEzMy43IDMzLjcgMCAwIDEgMi40IDcuNTVsLS42Mi45NEMuODkgOS44NyAwIDExLjQyIDAgMTNhMyAzIDAgMCAwIDMgM2gxMGEzIDMgMCAwIDAgMS42Ny01LjUuNzUuNzUgMCAwIDAtLjg0IDEuMjVBMS41IDEuNSAwIDEgMSAxMS41IDEzYTUgNSAwIDAgMSAuNjItMi4xNGMuNC0uNzguOTQtMS41OSAxLjUtMi40NGwuMDEtLjAyYy41Ni0uODMgMS4xNC0xLjcgMS41OS0yLjU5LjQ0LS44OC43OC0xLjgzLjc4LTIuODFhMyAzIDAgMCAwLTMtM3pNMyAxLjVBMS41IDEuNSAwIDAgMCAxLjUgM2MwIC4zMi4xLjU2LjI3Ljc3LjEuMTIuMzMuMjMuNzQuMjNoNy42N0EyLjc0IDIuNzQgMCAwIDEgMTAgM2MwLS41NS4xNS0xLjA2LjQtMS41SDN6bTEwIDBjLjgzIDAgMS41LjY2IDEuNSAxLjVhNSA1IDAgMCAxLS42MiAyLjE0Yy0uNC43OS0uOTQgMS42LTEuNSAyLjQ1bC0uMDIuMDJjLS41NS44My0xLjEzIDEuNy0xLjU4IDIuNThBNi4zMyA2LjMzIDAgMCAwIDEwIDEzYzAgLjU1LjE1IDEuMDYuNCAxLjVIM0ExLjUgMS41IDAgMCAxIDEuNSAxM2MwLTEuMDguNjMtMi4yOSAxLjU0LTMuN2wuNTUtLjg0QzQuMjMgNy41MSA0LjkgNi41IDUuMzggNS41aDYuMzhhLjc1Ljc1IDAgMCAwIC40Mi0xLjM3Yy0uNDUtLjMtLjY3LS42Ni0uNjctMS4xNEExLjUgMS41IDAgMCAxIDEzIDEuNXoiLz48L3N2Zz4=" alt="Lines of code">
</a>
<a href="https://github.com/everdrone/grab/releases/latest">
    <img src="https://img.shields.io/github/v/release/everdrone/grab?color=b7bdf8&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxNiAxNiIgd2lkdGg9IjE2IiBoZWlnaHQ9IjE2Ij48cGF0aCBmaWxsPSIjZDllMGVlIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiIGQ9Ik02LjEyMi4zOTJhMS43NSAxLjc1IDAgMDExLjc1NiAwbDUuMjUgMy4wNDVjLjU0LjMxMy44NzIuODkuODcyIDEuNTE0VjcuMjVhLjc1Ljc1IDAgMDEtMS41IDBWNS42NzdMNy43NSA4LjQzMnY2LjM4NGExIDEgMCAwMS0xLjUwMi44NjVMLjg3MiAxMi41NjNBMS43NSAxLjc1IDAgMDEwIDExLjA0OVY0Ljk1MWMwLS42MjQuMzMyLTEuMi44NzItMS41MTRMNi4xMjIuMzkyek03LjEyNSAxLjY5bDQuNjMgMi42ODVMNyA3LjEzMyAyLjI0NSA0LjM3NWw0LjYzLTIuNjg1YS4yNS4yNSAwIDAxLjI1IDB6TTEuNSAxMS4wNDlWNS42NzdsNC43NSAyLjc1NXY1LjUxNmwtNC42MjUtMi42ODNhLjI1LjI1IDAgMDEtLjEyNS0uMjE2em0xMC44MjggMy42ODRhLjc1Ljc1IDAgMTAxLjA4NyAxLjAzNGwyLjM3OC0yLjVhLjc1Ljc1IDAgMDAwLTEuMDM0bC0yLjM3OC0yLjVhLjc1Ljc1IDAgMDAtMS4wODcgMS4wMzRMMTMuNTAxIDEySDEwLjI1YS43NS43NSAwIDAwMCAxLjVoMy4yNTFsLTEuMTczIDEuMjMzeiI+PC9wYXRoPjwvc3ZnPg==&logoColor=d9e0ee&labelColor=1e1d2f&style=for-the-badge" alt="Latest Release">
</a>
<a href="https://app.codecov.io/gh/everdrone/grab" target="_blank">
    <img src="https://img.shields.io/codecov/c/github/everdrone/grab?color=c6a0f6&logo=codecov&logoColor=d9e0ee&labelColor=1e1d2f&style=for-the-badge&token=NkRjXNdxZI" alt="Codecov">
</a>
<a href="https://github.com/everdrone/grab/issues">
    <img src="https://img.shields.io/github/issues/everdrone/grab?color=f8bd96&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxNiAxNiIgd2lkdGg9IjE2IiBoZWlnaHQ9IjE2Ij48cGF0aCBmaWxsPSIjZDllMGVlIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiIGQ9Ik0xMC41NjEgMS41YS4wMTYuMDE2IDAgMDAtLjAxLjAwNEwzLjI4NiA4LjU3MUEuMjUuMjUgMCAwMDMuNDYyIDlINi43NWEuNzUuNzUgMCAwMS42OTQgMS4wMzRsLTEuNzEzIDQuMTg4IDYuOTgyLTYuNzkzQS4yNS4yNSAwIDAwMTIuNTM4IDdIOS4yNWEuNzUuNzUgMCAwMS0uNjgzLTEuMDZsMi4wMDgtNC40MTguMDAzLS4wMDZhLjAyLjAyIDAgMDAtLjAwNC0uMDA5LjAyLjAyIDAgMDAtLjAwNi0uMDA2TDEwLjU2IDEuNXpNOS41MDQuNDNhMS41MTYgMS41MTYgMCAwMTIuNDM3IDEuNzEzTDEwLjQxNSA1LjVoMi4xMjNjMS41NyAwIDIuMzQ2IDEuOTA5IDEuMjIgMy4wMDRsLTcuMzQgNy4xNDJhMS4yNSAxLjI1IDAgMDEtLjg3MS4zNTRoLS4zMDJhMS4yNSAxLjI1IDAgMDEtMS4xNTctMS43MjNMNS42MzMgMTAuNUgzLjQ2MmMtMS41NyAwLTIuMzQ2LTEuOTA5LTEuMjItMy4wMDRMOS41MDMuNDI5eiI+PC9wYXRoPjwvc3ZnPg==&logoColor=d9e0ee&labelColor=1e1d2f&style=for-the-badge" alt="GitHub issues">
</a>
</p>

# Table of contents

- [Motivation](#why)
- [Installation](#installation)
- [Usage](#usage)
- [Quickstart](#quickstart)
- [Options](#command-options)
- [Next steps](#next-steps)

# Why

This project helps you automate scraping data and downloading assets from the internet. Based on Go's Regular Expression engine and HCL, for ease of use, performance and flexibility.

# Installation

Download and install the [latest release](https://github.com/everdrone/grab/releases/latest).

# Usage

Run the following command to generate a new configuration file in the current directory.

```
grab config generate
```

> **Note**  
> Grab's configuration file uses [Hashicorp's HCL](https://github.com/hashicorp/hcl).  
> You can always refer to their specification for topics not covered by the documentation in this repo.

Once you're happy with your configuration, you can check if everything is ok by running:

```
grab config check
```

To scrape and download assets, pass one or more URLs to the `get` subcommand:

```ini
# single URL
grab get https://url.to/scrape/files?from

# list of URLs
grab get urls.ini

# at least one of each
grab get https://my.url/and urls.ini list.ini
```

> **Note**  
> The list of URLs can contain comments, like the `ini` format: all lines starting with `#` and `;` will be ignored.

# Quickstart

The default configuration, generated with `grab config generate` already works out of the box.

```hcl
global {
  location = "/home/yourusername/Downloads/grab"
}

site "unsplash" {
  test = "unsplash"

  asset "image" {
    pattern = "contentUrl\":\"([^\"]+)\""
    capture = 1

    transform filename {
      pattern = "(?:.+)photos\\/(.*)"
      replace = "$${1}.jpg"
    }
  }

  info "title" {
    pattern = "meta[^>]+property=\"og:title\"[^>]+content=\"(?P<title>[^\"]+)\""
    capture = "title"
  }

  subdirectory {
    pattern = "\\(@(?P<username>\\w+)\\)"
    capture = "username"
    from    = body
  }
}
```

For demonstration purposes, we can already download pictures from [unsplash](https://unsplash.com) by using the following command:

```
grab get https://unsplash.com/photos/uOi3lg8fGl4
```

> **Warning**  
> Please use this tool responsibly. Don't use this tool for Denial of Service attacks! Don't violate Copyright or intellectual property!

Internally, the program checks checks each URL passed to `get`, if it matches a `test` pattern inside of any `site` block, it will parse find all matches for assets or data defined in `asset` and `info` blocks.
Once all the asset URLs are gathered, the download starts.

After running the above command, you should have a new `grab` directory in your `~/Downloads` folder, containing subdirectories for each site defined in the configuration. Inside each site directories you will find all the assets extracted from the provided URLs.

The configuration syntax is based on a few fundamental blocks:

- `global` block defines the main download directory and global network options.
- `site <name>` blocks group other blocks based on the site URL.
- `asset <name>` blocks define what to look for from each site and how to download it.
- `info <name>` blocks define what strings to extract from the page body.

Additional configuration settings can be specified:

- `network` blocks to pass headers and other network options when making requests.
- `transform url` blocks to replace the asset URL before downloading.
- `transform filename` blocks to replace the asset's destination path.
- `subdirectory` blocks to organize downloads into subdirectories named by strings present in the page body or URL.

For a more in-depth look into Grab's confguration options, check out [the guide](/docs/guide.md).

# Command Options

To get help about any command, use the `help` subcommand or the `--help` flag:

```ini
# to list all available commands:
grab help

# to show instructions for a specific subcommand:
grab help <subcommand>
```

### `get`

#### Arguments

Accepts both URLs or path to lists of URLs. Both can be provided at the same time.

```sh
# grab get <url|file> [url|file...] [options]

grab get https://example.com/gallery/1 \
         https://example.com/gallery/2 \
         path/to/list.ini \
         other/file.ini -n
```

#### Options

| Long       | Short | Default | Description                                                                                                                    |
| ---------- | ----- | ------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `force`    | `f`   | `false` | To overwrite already existing files                                                                                            |
| `config`   | `c`   | `nil`   | To specify the path to a configuration file                                                                                    |
| `strict`   | `s`   | `false` | To stop the program at the first encountered error                                                                             |
| `dry-run`  | `n`   | `false` | To send requests without writing to the disk                                                                                   |
| `progress` | `p`   | `false` | To show a progress bar                                                                                                         |
| `quiet`    | `q`   | `false` | To suppress all output to `stdout` (errors will still be printed to `stderr`).<br/>This option takes precedence over `verbose` |
| `verbose`  | `v`   | `1`     | To set the verbosity level:<br/>`-v` is 1, `-vv` is 2 and so on...<br/>`quiet` overrides this option.                          |

## Next steps

- [x] Retries & Timeout
- [x] Network options with inheritance
- [x] URL manipulation
- [x] Destination manipulation
- [ ] Display a progress bar
- [ ] Improve logging
- [ ] Add HCL eval context functions
- [ ] Distribute via various package managers:
  - [ ] Homebrew
  - [ ] Apt
  - [ ] Chocolatey
  - [ ] Scoop
- [ ] Scripting language integration
- [ ] Plugin system
- [ ] Sequential jobs (like GitHub workflows)

## Credits

- [Catppuccin](https://github.com/catppuccin/) for the color palette
- [Shields.io](https://github.com/badges/shields) for the badges

## License

Distributed under the [MIT License](/LICENSE).
