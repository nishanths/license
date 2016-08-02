# [license](https://github.com/nishanths/license) 

[![wercker status](https://app.wercker.com/status/1407b8c71c720358bf15eeb5815f99bd/s "wercker status")](https://app.wercker.com/project/bykey/1407b8c71c720358bf15eeb5815f99bd)

Create licenses from the command-line. *Hello, productivity!*

* Supports all the licenses available on GitHub
* Does not need network access (except on first run)
* Good defaults for name and year on the license; easy to customize when needed

<br>
<img src="https://zippy.gfycat.com/JoyfulBlandGermanshorthairedpointer.gif" width="700px"/>
<br>

## Install

Using go:

```
go get -u github.com/nishanths/license
``` 

[more info](https://golang.org/doc/install#install)

## Usage

#### Generate a license

To print a license to stdout, run `license` followed by the license name:

```
license mit
```

#### Save to file

Use the `-o` flag to save the license to a file:

```
license -o LICENSE.txt mit
```

#### Customize name and year

```
license -year=2013 -name=Alice isc
```

If unspecified, the current year is used.

To determine the name, license uses the following in order. Since you likely have your name set in `.gitconfig`, you can always omit the `-name` flags.

```
- command line flags: -name, -n
- environment variable: LICENSE_FULL_NAME
- gitconfig and hgconfig
- "os/user".Current()
- empty string
```

## Authentication

If you receive a `403 Forbidden: API rate limit exceeded` while updating licenses, use your GitHub username and a [personal access token](https://github.com/settings/tokens) (no scopes required).

```
license -auth username:e0a8a01b1f125a785ea3d7ada98eb6a018e2fe4f -update
```

(The token above will not work.)

## More docs

Help is available by runnning `license -help`

View the list of available licenses by running `license -list`


## Contributing

Pull requests for new features, bug fixes, and suggestions are welcome!

## License

[MIT](https://github.com/nishanths/license/blob/master/LICENSE)
