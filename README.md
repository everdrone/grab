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

- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
  - [Downloading Assets](#downloading-assets)
  - [Indexing Data](#indexing-data)
  - [Subdirectories](#subdirectories)
  - [Network Options](#network-options)
- [Options](#command-options)
- [Next steps](#next-steps)

# Installation

Download and install the [latest release](https://github.com/everdrone/grab/releases/latest)

# Usage

Let's start fresh. Run the following command to generate a new configuration file in the current directory

```
grab config generate
```

The file `grab.hcl` should be located in your home directory, or in any parent directory from where you will call the command.

The language of the file is [Hashicorp Configuration Language](https://github.com/hashicorp/hcl)

Read more about the configuration options [here](#configuration)

Once you're happy with your configuration, you can check if everything is ok by running

```
grab config check
```

or, if your file is not located in a parent directory from your current working directory, you can always specify its path with the `--config` option.

```
grab config check -c /var/grab.hcl
```

Now you can start using grab.
To scrape and download assets use the `grab` command and pass at least one url or a file containing a list of urls.

> **Note**
> The list of urls can contain comments, like the `ini` format, all lines starting with `#` and `;` will be ignored

```ini
# single URL
grab get https://url.to/scrape/files?from

# list of URLs
grab get urls.ini

# at least one of each
grab get https://my.url/and urls.ini list.ini
```

# Configuration

Take this example configuration:

```hcl
global {
    location = "/home/user/Downloads/grab"
}

site "example" {
    test = "example\\.com"

    asset "image" {
        pattern  = "<img src=\"([^\"]+)\""
        find_all = true
        capture  = 1
    }

    info "editor" {
        pattern = "editor:\\s@(\\w+)"
        capture = 1
    }

    subdirectory {
        pattern = "gallery\\/(?<id>\\d+)\/"
        capture = "id"
        from    = url
    }
}
```

Let's pass our hypothetical url to `grab get`

```
grab get https://example.com/gallery/1337/overview
```

### Downloading assets

The program will check if our url matches with any `site` block using the `test` pattern. If the pattern matches, the program will fetch the page body to scrape its contents.

```hcl
asset "video" {
    pattern  = "<video src=\"(?P<videourl>[^\"]+)\""
    capture  = "videourl"
    find_all = true  # optional
}
```

> **Note**
> To escape double quotes, you must use one backslash: `\"`
> To escape common regex expressions like `\d` you should escape twice: `\\d`

For each `asset` block, grab will search for matches using the `pattern` regex and then extract the `capture` group from the matches.
By default only the first match will be extracted, if you wish to extract multiple urls from the same page, you can set `find_all` to `true`. Finally all the files will be downloaded from the extracted urls.

> Example of extracted urls with the onfiguration above:
>
> ```
> https://cdn.example.com/img/image1.jpg
> https://cdn.example.com/img/image2.jpg
> https://cdn.example.com/img/image3.jpg
> ```

### Indexing data

After the assets, all the `info` blocks are evaluated and information is extracted from the page and will be stored in a `_info.json` file.

```hcl
info "phone" {
    pattern  = "tel:(\d+)"
    capture  = 1
}
```

Inside the info file, two additional properties will be set by default: `url` and `timestamp`, representing the page url where the information has been scraped from, and the current time.

> Example `_info.json` output:
>
> ```json
> {
>   "url": "https://example.com/gallery/1337/overview",
>   "timestamp": "2022-08-17T13:51:58.7265822Z",
>   "editor": "everdrone"
> }
> ```

By default, grab creates a subdirectory with the site name (in this case `example`) to store the information downloaded from this site.
If you want to create separate subdirectories under `example` you can specify a `subdirectory` block.

### Subdirectories

The `subdirectory` block will extract a string using `pattern` and `capture` just like other blocks, but you can specify the `from` attribute to tell grab to search inside the `url` or inside the `body`

```hcl
subdirectory {
    pattern = "href=\"\\/\\@(?P<user>[^\"]+)"
    capture = "user"
    from    = body  # defaults to url
}
```

The final path of the assets will be `<global.location>/<site.name>/<subdirectory>/<filename>`

> Example of destinations from the configuration above:
>
> ```
> /home/user/Downloads/grab/example/1337/image1.jpg
> ```
>
> Similarly, the `_info.json` file will be saved to `/home/user/Downloads/grab/example/1337/_info.json`

If no `subdirectory` block is specified, the asset destination will conform to: `<global.location>/<site.name>/<filename>`

> **Warning**
> If the `pattern` attribute contains named groups, you must set the `capture` attribute to get the named capture.
>
> Use an integer `capture` groups only if your `pattern` does not contain named groups.
> To learn more about Go's regexp syntax see the [official documentation](https://pkg.go.dev/regexp/syntax).

### Network options

If a site requires specific headers to be set, or a number of retries, you can add optional `network` blocks to your configuration file.

```hcl
network {
    # all attributes are optional
    retries = 3
    timeout = 10000  # in milliseconds
    headers = {
        "User-Agent" = "Mozilla/5.0 ..."
    }
}
```

`network` blocks can be located in the `global` block, inside `site` blocks and even `asset` blocks.

By default, the `global.network` configuration will be inherited to all sites and all site assets. To avoid inheriting the network configuration of a parent block, you can set `inherit = false` like so:

```hcl
site "example" {
    # ...

    network {
        inherit = false
    }

    # ...
}
```

To learn more about advanced configuration patterns, see [Advanced Configuration](/docs/advanced.md)

# Command Options

### `get`

#### Arguments

Accepts both urls or path to lists of urls. Both can be provided at the same time.

```sh
# grab get <url|file> [url|file...] [options]

grab get https://example.com/gallery/1 \
         https://example.com/gallery/2 \
         path/to/list.ini \
         other/file.ini -n
```

#### Options

| Long       | Short | Default | Description                                                                                                                |
| ---------- | ----- | ------- | -------------------------------------------------------------------------------------------------------------------------- |
| `force`    | `f`   | `false` | Overwrites already existing files                                                                                          |
| `config`   | `c`   | `nil`   | Specify the path to a configuration file                                                                                   |
| `strict`   | `s`   | `false` | Will stop the program at the first encountered error                                                                       |
| `dry-run`  | `n`   | `false` | Will send requests without writing to the disk                                                                             |
| `progress` | `p`   | `false` | Show a progress bar                                                                                                        |
| `quiet`    | `q`   | `false` | Suppress all output to `stdout` (errors will still be printed to `stderr`)<br/>This option takes precedence over `verbose` |
| `verbose`  | `v`   | `1`     | Set the verbosity level.<br/>`-v` is 1, `-vv` is 2 and so on...<br/>`quiet` overrides this option.                         |

## Next steps

- [x] Retries & Timeout
- [x] Network options inheritance
- [x] URL manipulation
- [x] Destination manipulation
- [ ] Display progress bar
- [ ] Better logging
- [ ] Add HCL eval context functions
- [ ] Distribute via various package managers:
  - [ ] Homebrew
  - [ ] Apt
  - [ ] Chocolatey
  - [ ] Scoop
- [ ] Scripting language integration
- [ ] Plugins?
- [ ] Sequential jobs (like GitHub workflows)

## License

Distributed under the [MIT License](/LICENSE).
