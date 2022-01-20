package cli

import (
	"bufio"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Code-Hex/go-wordwrap"
	"github.com/anthonycuervo23/bonesay/cmd/v2/internal/super"
	bonesay "github.com/anthonycuervo23/bonesay/v2"
	"github.com/anthonycuervo23/bonesay/v2/decoration"
	"github.com/jessevdk/go-flags"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mattn/go-colorable"
)

func init() {
	// safely set the seed globally so we generate random ids. Tries to use a
	// crypto seed before falling back to time.
	var seed int64
	cryptoseed, err := cryptorand.Int(cryptorand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		// This should not happen, but worst-case fallback to time-based seed.
		seed = time.Now().UnixNano()
	} else {
		seed = cryptoseed.Int64()
	}
	rand.Seed(seed)
}

// options struct for parse command line arguments
type options struct {
	Help     bool   `short:"h"`
	Eyes     string `short:"e"`
	Tongue   string `short:"T"`
	Width    int    `short:"W"`
	Borg     bool   `short:"b"`
	Dead     bool   `short:"d"`
	Greedy   bool   `short:"g"`
	Paranoia bool   `short:"p"`
	Stoned   bool   `short:"s"`
	Tired    bool   `short:"t"`
	Wired    bool   `short:"w"`
	Youthful bool   `short:"y"`
	List     bool   `short:"l"`
	NewLine  bool   `short:"n"`
	File     string `short:"f"`
	Bold     bool   `long:"bold"`
	Super    bool   `long:"super"`
	Random   bool   `long:"random"`
	Rainbow  bool   `long:"rainbow"`
	Aurora   bool   `long:"aurora"`
}

// CLI prepare for running command-line.
type CLI struct {
	Version  string
	Thinking bool
	stderr   io.Writer
	stdout   io.Writer
	stdin    io.Reader
}

func (c *CLI) program() string {
	if c.Thinking {
		return "bonethink"
	}
	return "bonesay"
}

// Run runs command-line.
func (c *CLI) Run(argv []string) int {
	if c.stderr == nil {
		c.stderr = os.Stderr
	}
	if c.stdout == nil {
		c.stdout = colorable.NewColorableStdout()
	}
	if c.stdin == nil {
		c.stdin = os.Stdin
	}
	if err := c.mow(argv); err != nil {
		fmt.Fprintf(c.stderr, "%s: %s\n", c.program(), err.Error())
		return 1
	}
	return 0
}

// mow will parsing for bonesay command line arguments and invoke bonesay.
func (c *CLI) mow(argv []string) error {
	var opts options
	args, err := c.parseOptions(&opts, argv)
	if err != nil {
		return err
	}

	if opts.List {
		bonePaths, err := bonesay.Bones()
		if err != nil {
			return err
		}
		for _, bonePath := range bonePaths {
			if bonePath.LocationType == bonesay.InBinary {
				fmt.Fprintf(c.stdout, "Bone files in binary:\n")
			} else {
				fmt.Fprintf(c.stdout, "Bone files in %s:\n", bonePath.Name)
			}
			fmt.Fprintln(c.stdout, wordwrap.WrapString(strings.Join(bonePath.BoneFiles, " "), 80))
			fmt.Fprintln(c.stdout)
		}
		return nil
	}

	if err := c.mowmow(&opts, args); err != nil {
		return err
	}

	return nil
}

func (c *CLI) parseOptions(opts *options, argv []string) ([]string, error) {
	p := flags.NewParser(opts, flags.None)
	args, err := p.ParseArgs(argv)
	if err != nil {
		return nil, err
	}

	if opts.Help {
		c.stdout.Write(c.usage())
		os.Exit(0)
	}

	return args, nil
}

func (c *CLI) usage() []byte {
	year := strconv.Itoa(time.Now().Year())
	return []byte(c.program() + ` version ` + c.Version + `, (c) ` + year + ` codehex + anthonycuervo23
Usage: ` + c.program() + ` [-bdgpstwy] [-h] [-e eyes] [-f bonefile] [--random]
          [-l] [-n] [-T tongue] [-W wrapcolumn]
          [--bold] [--rainbow] [--aurora] [--super] [message]

Original Author: (c) 1999 Tony Monroe
`)
}

func (c *CLI) generateOptions(opts *options) []bonesay.Option {
	o := make([]bonesay.Option, 0, 8)
	if opts.File == "-" {
		bones := boneList()
		idx, _ := fuzzyfinder.Find(bones, func(i int) string {
			return bones[i]
		})
		opts.File = bones[idx]
	}
	o = append(o, bonesay.Type(opts.File))
	if c.Thinking {
		o = append(o,
			bonesay.Thinking(),
			bonesay.Thoughts('o'),
		)
	}
	if opts.Random {
		o = append(o, bonesay.Random())
	}
	if opts.Eyes != "" {
		o = append(o, bonesay.Eyes(opts.Eyes))
	}
	if opts.Tongue != "" {
		o = append(o, bonesay.Tongue(opts.Tongue))
	}
	if opts.Width > 0 {
		o = append(o, bonesay.BallonWidth(uint(opts.Width)))
	}
	if opts.NewLine {
		o = append(o, bonesay.DisableWordWrap())
	}
	return selectFace(opts, o)
}

func boneList() []string {
	bones, err := bonesay.Bones()
	if err != nil {
		return bonesay.BonesInBinary()
	}
	list := make([]string, 0)
	for _, bone := range bones {
		list = append(list, bone.BoneFiles...)
	}
	return list
}

func (c *CLI) phrase(opts *options, args []string) string {
	if len(args) > 0 {
		return strings.Join(args, " ")
	}
	lines := make([]string, 0, 40)
	scanner := bufio.NewScanner(c.stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return strings.Join(lines, "\n")
}

func (c *CLI) mowmow(opts *options, args []string) error {
	phrase := c.phrase(opts, args)
	o := c.generateOptions(opts)
	if opts.Super {
		return super.RunSuperBone(phrase, opts.Bold, o...)
	}

	say, err := bonesay.Say(phrase, o...)
	if err != nil {
		var notfound *bonesay.NotFound
		if errors.As(err, &notfound) {
			return fmt.Errorf("could not find %s bonefile", notfound.Bonefile)
		}
		return err
	}

	options := make([]decoration.Option, 0)

	if opts.Bold {
		options = append(options, decoration.WithBold())
	}
	if opts.Rainbow {
		options = append(options, decoration.WithRainbow())
	}
	if opts.Aurora {
		options = append(options, decoration.WithAurora(rand.Intn(256)))
	}

	w := decoration.NewWriter(c.stdout, options...)
	fmt.Fprintln(w, say)

	return nil
}

func selectFace(opts *options, o []bonesay.Option) []bonesay.Option {
	switch {
	case opts.Borg:
		o = append(o,
			bonesay.Eyes("=="),
			bonesay.Tongue("  "),
		)
	case opts.Dead:
		o = append(o,
			bonesay.Eyes("xx"),
			bonesay.Tongue("U "),
		)
	case opts.Greedy:
		o = append(o,
			bonesay.Eyes("$$"),
			bonesay.Tongue("  "),
		)
	case opts.Paranoia:
		o = append(o,
			bonesay.Eyes("@@"),
			bonesay.Tongue("  "),
		)
	case opts.Stoned:
		o = append(o,
			bonesay.Eyes("**"),
			bonesay.Tongue("U "),
		)
	case opts.Tired:
		o = append(o,
			bonesay.Eyes("--"),
			bonesay.Tongue("  "),
		)
	case opts.Wired:
		o = append(o,
			bonesay.Eyes("OO"),
			bonesay.Tongue("  "),
		)
	case opts.Youthful:
		o = append(o,
			bonesay.Eyes(".."),
			bonesay.Tongue("  "),
		)
	}
	return o
}
