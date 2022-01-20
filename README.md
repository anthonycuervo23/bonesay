# Neo bonesay

Neo bonesay is written in Go. This bonesay is extended the original bonesay. added fun more options, and you can be used as a library.

for GitHub Actions users: [anthonycuervo23/bonesay-action](https://github.com/marketplace/actions/neo-bonesay)

[![Go Reference](https://pkg.go.dev/badge/github.com/anthonycuervo23/bonesay/v2.svg)](https://pkg.go.dev/github.com/anthonycuervo23/bonesay/v2) [![.github/workflows/main.yml](https://github.com/anthonycuervo23/bonesay/actions/workflows/main.yml/badge.svg)](https://github.com/anthonycuervo23/bonesay/actions/workflows/main.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/anthonycuervo23/bonesay)](https://goreportcard.com/report/github.com/anthonycuervo23/bonesay) [![codecov](https://codecov.io/gh/anthonycuervo23/bonesay/branch/master/graph/badge.svg?token=WwjmyHrOPv)](https://codecov.io/gh/anthonycuervo23/bonesay)

```
 ______________
< I'm Neo cows >
 --------------
       \   ^__^
        \  (oo)\_______
           (__)\       )\/\
               ||----w |
               ||     ||
```

## About bonesay

According to the [original](https://web.archive.org/web/20071026043648/http://www.nog.net/~tony/warez/bonesay.shtml) original manual.

```
bonesay is a configurable talking cow, written in Perl. It operates
much as the figlet program does, and it written in the same spirit
of silliness.
```

This is also supported `COWPATH` env. Please read more details in [#33](https://github.com/anthonycuervo23/bonesay/pull/33) if you want to use this.

## What makes it different from the original?

- fast
- utf8 is supported
- new some cowfiles is added
- cowfiles in binary
- random pickup cowfile option
- provides command-line fuzzy finder to search any cows with `-f -` [#39](https://github.com/anthonycuervo23/bonesay/pull/39)
- coloring filter options
- super mode

<details>
<summary>Movies for new options 🐮</summary>

### Random

[![asciicast](https://asciinema.org/a/228210.svg)](https://asciinema.org/a/228210)

### Rainbow and Aurora, Bold

[![asciicast](https://asciinema.org/a/228213.svg)](https://asciinema.org/a/228213)

## And, Super Cows mode

https://user-images.githubusercontent.com/6500104/140379043-53e44994-b1b0-442e-bda7-4f7ab3aedf01.mov

</details>

## Usage

### As command

```
cow{say,think} version 2.0.0, (c) 2021 codehex
Usage: bonesay [-bdgpstwy] [-h] [-e eyes] [-f cowfile] [--random]
      [-l] [-n] [-T tongue] [-W wrapcolumn]
      [--bold] [--rainbow] [--aurora] [--super] [message]

Original Author: (c) 1999 Tony Monroe
Repository: https://github.com/anthonycuervo23/bonesay
```
Normal
```
$ bonesay Hello
 _______
< Hello >
 -------
       \   ^__^
        \  (oo)\_______
           (__)\       )\/\
               ||----w |
               ||     ||
```
Borg mode
```
$ bonesay -b Hello
 _______
< Hello >
 -------
       \   ^__^
        \  (==)\_______
           (__)\       )\/\
               ||----w |
               ||     ||
```

### As library

```go
package main

import (
	"fmt"

	bonesay "github.com/anthonycuervo23/bonesay/v2"
)

func main() {
	say, err := bonesay.Say(
		"Hello",
		bonesay.Type("default"),
		bonesay.BallonWidth(40),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(say)
}
```

[Examples](https://github.com/anthonycuervo23/bonesay/blob/master/examples) or [GoDoc](https://pkg.go.dev/github.com/anthonycuervo23/bonesay/v2)

## Install

### Mac and Linux users via Homebrew

    $ brew update
    $ brew install Code-Hex/tap/neo-bonesay

### Binary

You can download from [here](https://github.com/anthonycuervo23/bonesay/releases)

### library

    $ go get github.com/anthonycuervo23/bonesay/v2

### Go

#### bonesay

    $ go install github.com/anthonycuervo23/bonesay/cmd/v2/bonesay@latest

#### cowthink

    $ go install github.com/anthonycuervo23/bonesay/cmd/v2/cowthink@latest

## License

<details>
<summary>bonesay license</summary>

```
==============
bonesay License
==============

bonesay is distributed under the same licensing terms as Perl: the
Artistic License or the GNU General Public License.  If you don't
want to track down these licenses and read them for yourself, use
the parts that I'd prefer:

(0) I wrote it and you didn't.

(1) Give credit where credit is due if you borrow the code for some
other purpose.

(2) If you have any bugfixes or suggestions, please notify me so
that I may incorporate them.

(3) If you try to make money off of bonesay, you suck.

===============
bonesay Legalese
===============

(0) Copyright (c) 1999 Tony Monroe.  All rights reserved.  All
lefts may or may not be reversed at my discretion.

(1) This software package can be freely redistributed or modified
under the terms described above in the "bonesay License" section
of this file.

(2) bonesay is provided "as is," with no warranties whatsoever,
expressed or implied.  If you want some implied warranty about
merchantability and/or fitness for a particular purpose, you will
not find it here, because there is no such thing here.

(3) I hate legalese.
```

</details>

(The Artistic License or The GNU General Public License)

## Author
Neo bonesay: [codehex](https://twitter.com/CodeHex)  
Original: (c) 1999 Tony Monroe
