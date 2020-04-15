# license [![wercker status](https://app.wercker.com/status/1407b8c71c720358bf15eeb5815f99bd/s "wercker status")](https://app.wercker.com/project/bykey/1407b8c71c720358bf15eeb5815f99bd)

## Install

```
go get github.com/nishanths/license
``` 

Create LICENSE files from the command-line. 

* Good defaults for name and year on license; easy to customize when needed
* Supports all license types available on [GitHub](https://developer.github.com/v3/licenses/).
```
agpl-3.0      (GNU Affero General Public License v3.0)
apache-2.0    (Apache License 2.0)
bsd-2-clause  (BSD 2-Clause "Simplified" License)
bsd-3-clause  (BSD 3-Clause "New" or "Revised" License)
cc0-1.0       (Creative Commons Zero v1.0 Universal)
epl-2.0       (Eclipse Public License 2.0)
gpl-2.0       (GNU General Public License v2.0)
gpl-3.0       (GNU General Public License v3.0)
lgpl-2.1      (GNU Lesser General Public License v2.1)
lgpl-3.0      (GNU Lesser General Public License v3.0)
mit           (MIT License)
mpl-2.0       (Mozilla Public License 2.0)
unlicense     (The Unlicense)
```

<br>
<img src="https://zippy.gfycat.com/JoyfulBlandGermanshorthairedpointer.gif" width="700px"/>
<br>

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
$ license -year 2013 -name "Alice G" isc
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

### Help

Help is available by runnning `license -help`

## Contributing

Pull requests for new features, bug fixes, and suggestions are welcome!

## License

[MIT](https://github.com/nishanths/license/blob/master/LICENSE)
