`license` is a command line tool to create LICENSE files.

It provides good defaults for name and year on license (customizable if needed), and
it supports all license types listed on the [GitHub Licenses API](https://developer.github.com/v3/licenses/) and a few more. The license templates used by this program were copied from the GitHub Licenses API, when available.
```
agpl-3.0      (GNU Affero General Public License v3.0)
apache-2.0    (Apache License 2.0)
bsd-2-clause  (BSD 2-Clause "Simplified" License)
bsd-3-clause  (BSD 3-Clause "New" or "Revised" License)
cc0-1.0       (Creative Commons Zero v1.0 Universal)
epl-2.0       (Eclipse Public License 2.0)
free-art-1.3  (Free Art License 1.3)
gpl-2.0       (GNU General Public License v2.0)
gpl-3.0       (GNU General Public License v3.0)
lgpl-2.1      (GNU Lesser General Public License v2.1)
lgpl-3.0      (GNU Lesser General Public License v3.0)
mit           (MIT License)
mpl-2.0       (Mozilla Public License 2.0)
unlicense     (The Unlicense)
```

## Install

### Building from Source

Outside a project using Go modules, get the latest version by running:

```
go get github.com/nishanths/license
```

Inside a project using Go modules, use:

```
go get github.com/nishanths/license/v5
```

### Via the Arch User Repository (AUR)

This program is available via the AUR under the name [`nishanths-license-git`](https://aur.archlinux.org/packages/nishanths-license-git/). Using yay, you can install it like so:

```bash
yay -S nishanths-license-git
```

## Usage

#### Print license

To print a license to stdout, run the `license` command followed by the license name:

```
$ license mit
```

#### Save to file

Use the `-o` flag to save the license to a file, or use your shell's redirection operator:

```
$ license -o LICENSE.txt mit
$ license mit > LICENSE.txt
```

#### Customize name and year

```
$ license -year 2013 -name "Alice L" isc
```

The current year is used if `-year` is omitted.

To determine the name on the license, the following are used in this order:

```
- command line flags: -name, -n
- environment variable: LICENSE_FULL_NAME
- gitconfig and hgconfig
- "os/user".Current()
- empty string
```

If you have your name set in `$HOME/.gitconfig`, you can almost always omit the `-name` flag.

#### Demo

![Demonstration](demo.gif)

## Contributing

Pull requests for new features, bug fixes, and suggestions are welcome!

## License

[MIT](https://github.com/nishanths/license/blob/master/LICENSE)
