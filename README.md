# license

Command-line license generator written in Go.

````
Command-line license generator.

Usage:
    license [-y <year>] [-n <name>] [-o <filename>] <license-name>

Examples:
    license mit
    license -o LICENSE.txt mit
    license -y 2013 -n Alice isc

Additional commands:
    ls            list locally available license names
    ls-remote     list remote license names
    update        update local licenses to latest remote versions
    help          show help information
    version       print current version

Run "license ls" to see list of available license names.
````
<video id="sampleMovie" src="https://zippy.gfycat.com/JoyfulBlandGermanshorthairedpointer.webm" autoplay muted loop></video>

# Contents

* Install
* Get Started
* Options
* Contributing
* License

# Install

To install license, run:

````
$ go get github.com/nishanths/license
````

Otherwise, 

# Get Started

#### Generating a license

To generate a license, simply run `license` followed by the license name. The following command generates the MIT license:

````bash
$ license mit
````

#### Creating a license file

Use the `-o` option to save the license to a file. The following command creates the file `LICENSE.txt` with the contents of the ISC license:

````
$ license -o LICENSE.txt isc
```` 

More options and commands are described section below.

# Usage