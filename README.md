# license [![wercker status](https://app.wercker.com/status/1407b8c71c720358bf15eeb5815f99bd/s "wercker status")](https://app.wercker.com/project/bykey/1407b8c71c720358bf15eeb5815f99bd)

## Install

```
go get github.com/nishanths/license
``` 

Create LICENSE files from the command-line.

* Supports all the [licenses available on GitHub](https://developer.github.com/v3/licenses/)
* Good defaults for name and year on license; easy to customize when needed

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
$ license -o=LICENSE.txt mit
$ license mit > LICENSE.txt
```

#### Customize name and year

```
$ license -year=2013 -name=Alice isc
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

### List of licenses

```
$ license -list

```

### Help

Help is available by runnning `license -help`

## Contributing

Pull requests for new features, bug fixes, and suggestions are welcome!

## License

[MIT](https://github.com/nishanths/license/blob/master/LICENSE)
