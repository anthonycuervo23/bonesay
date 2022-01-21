package bonesay

import (
	"fmt"
	"math/rand"
	"strings"
)

// Bone struct!!
type Bone struct {
	eyes            string
	tongue          string
	typ             *BoneFile
	thoughts        rune
	thinking        bool
	ballonWidth     int
	disableWordWrap bool
	balloonOffset   int

	buf strings.Builder
}

// New returns pointer of Bone struct that made by options
func New(options ...Option) (*Bone, error) {
	bone := &Bone{
		eyes:     "oo",
		tongue:   "  ",
		thoughts: '/',
		typ: &BoneFile{
			Name:         "mobile",
			BasePath:     "bones",
			LocationType: InBinary,
		},
		ballonWidth: 15,
	}
	for _, o := range options {
		if err := o(bone); err != nil {
			return nil, err
		}
	}
	return bone, nil
}

// Say returns string that said by bone
func (bone *Bone) Say(phrase string) (string, error) {
	mow, err := bone.GetBone()
	if err != nil {
		return "", err
	}
	return bone.Balloon(phrase) + mow, nil
}

// Clone returns a copy of bone.
//
// If any options are specified, they will be reflected.
func (bone *Bone) Clone(options ...Option) (*Bone, error) {
	ret := new(Bone)
	*ret = *bone
	ret.buf.Reset()
	for _, o := range options {
		if err := o(ret); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

// Option defined for Options
type Option func(*Bone) error

// Eyes specifies eyes
// The specified string will always be adjusted to be equal to two characters.
func Eyes(s string) Option {
	return func(c *Bone) error {
		c.eyes = adjustTo2Chars(s)
		return nil
	}
}

// Tongue specifies tongue
// The specified string will always be adjusted to be less than or equal to two characters.
func Tongue(s string) Option {
	return func(c *Bone) error {
		c.tongue = adjustTo2Chars(s)
		return nil
	}
}

func adjustTo2Chars(s string) string {
	if len(s) >= 2 {
		return s[:2]
	}
	if len(s) == 1 {
		return s + " "
	}
	return "  "
}

func containBones(target string) (*BoneFile, error) {
	bonePaths, err := Bones()
	if err != nil {
		return nil, err
	}
	for _, bonePath := range bonePaths {
		bonefile, ok := bonePath.Lookup(target)
		if ok {
			return bonefile, nil
		}
	}
	return nil, nil
}

// NotFound is indicated not found the bonefile.
type NotFound struct {
	Bonefile string
}

var _ error = (*NotFound)(nil)

func (n *NotFound) Error() string {
	return fmt.Sprintf("not found %q bonefile", n.Bonefile)
}

// Type specify name of the bonefile
func Type(s string) Option {
	if s == "" {
		s = "default"
	}
	return func(c *Bone) error {
		bonefile, err := containBones(s)
		if err != nil {
			return err
		}
		if bonefile != nil {
			c.typ = bonefile
			return nil
		}
		return &NotFound{Bonefile: s}
	}
}

// Thinking enables thinking mode
func Thinking() Option {
	return func(c *Bone) error {
		c.thinking = true
		return nil
	}
}

// Thoughts Thoughts allows you to specify
// the rune that will be drawn between
// the speech bubbles and the bone
func Thoughts(thoughts rune) Option {
	return func(c *Bone) error {
		c.thoughts = thoughts
		return nil
	}
}

// Random specifies something .bone from bones directory
func Random() Option {
	pick, err := pickBone()
	return func(c *Bone) error {
		if err != nil {
			return err
		}
		c.typ = pick
		return nil
	}
}

func pickBone() (*BoneFile, error) {
	bonePaths, err := Bones()
	if err != nil {
		return nil, err
	}
	bonePath := bonePaths[rand.Intn(len(bonePaths))]

	n := len(bonePath.BoneFiles)
	bonefile := bonePath.BoneFiles[rand.Intn(n)]
	return &BoneFile{
		Name:         bonefile,
		BasePath:     bonePath.Name,
		LocationType: bonePath.LocationType,
	}, nil
}

// BallonWidth specifies ballon size
func BallonWidth(size uint) Option {
	return func(c *Bone) error {
		c.ballonWidth = int(size)
		return nil
	}
}

// DisableWordWrap disables word wrap.
// Ignoring width of the ballon.
func DisableWordWrap() Option {
	return func(c *Bone) error {
		c.disableWordWrap = true
		return nil
	}
}
