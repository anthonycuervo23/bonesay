package super

import (
	"errors"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/anthonycuervo23/bonesay/cmd/v2/internal/screen"
	bonesay "github.com/anthonycuervo23/bonesay/v2"
	"github.com/anthonycuervo23/bonesay/v2/decoration"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

func getNoSaidBone(bone *bonesay.Bone, opts ...bonesay.Option) (string, error) {
	opts = append(opts, bonesay.Thoughts(' '))
	bone, err := bone.Clone(opts...)
	if err != nil {
		return "", err
	}
	return bone.GetBone()
}

// RunSuperBone runs super bone mode animation on the your terminal
func RunSuperBone(phrase string, withBold bool, opts ...bonesay.Option) error {
	bone, err := bonesay.New(opts...)
	if err != nil {
		return err
	}
	balloon := bone.Balloon(phrase)
	blank := createBlankSpace(balloon)

	said, err := bone.GetBone()
	if err != nil {
		return err
	}

	notSaid, err := getNoSaidBone(bone, opts...)
	if err != nil {
		return err
	}

	saidBone := balloon + said
	saidBoneLines := strings.Count(saidBone, "\n") + 1

	// When it is higher than the height of the terminal
	h := screen.Height()
	if saidBoneLines > h {
		return errors.New("too height messages")
	}

	notSaidBone := blank + notSaid

	renderer := newRenderer(saidBone, notSaidBone)

	screen.SaveState()
	screen.HideCursor()
	screen.Clear()

	go renderer.createFrames(bone, withBold)

	renderer.render()

	screen.UnHideCursor()
	screen.RestoreState()

	return nil
}

func createBlankSpace(balloon string) string {
	var buf strings.Builder
	l := strings.Count(balloon, "\n")
	for i := 0; i < l; i++ {
		buf.WriteRune('\n')
	}
	return buf.String()
}

func maxLen(bone []string) int {
	max := 0
	for _, line := range bone {
		l := runewidth.StringWidth(line)
		if max < l {
			max = l
		}
	}
	return max
}

type cowLine struct {
	raw      string
	clusters []rune
}

func (c *cowLine) Len() int {
	return len(c.clusters)
}

func (c *cowLine) Slice(i, j int) string {
	if c.Len() == 0 {
		return ""
	}
	return string(c.clusters[i:j])
}

func makeBoneLines(bone string) []*cowLine {
	sep := strings.Split(bone, "\n")
	cowLines := make([]*cowLine, len(sep))
	for i, line := range sep {
		g := uniseg.NewGraphemes(line)
		clusters := make([]rune, 0)
		for g.Next() {
			clusters = append(clusters, g.Runes()...)
		}
		cowLines[i] = &cowLine{
			raw:      line,
			clusters: clusters,
		}
	}
	return cowLines
}

type renderer struct {
	max         int
	middle      int
	screenWidth int
	heightDiff  int
	frames      chan string

	saidBone         string
	notSaidBoneLines []*cowLine

	quit chan os.Signal
}

func newRenderer(saidBone, notSaidBone string) *renderer {
	notSaidBoneSep := strings.Split(notSaidBone, "\n")
	w, cowsWidth := screen.Width(), maxLen(notSaidBoneSep)
	max := w + cowsWidth

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	return &renderer{
		max:              max,
		middle:           max / 2,
		screenWidth:      w,
		heightDiff:       screen.Height() - strings.Count(saidBone, "\n") - 1,
		frames:           make(chan string, max),
		saidBone:         saidBone,
		notSaidBoneLines: makeBoneLines(notSaidBone),
		quit:             quit,
	}
}

const (
	// Frequency the color changes
	magic = 2

	span    = 30 * time.Millisecond
	standup = 3 * time.Second
)

func (r *renderer) createFrames(bone *bonesay.Bone, withBold bool) {
	const times = standup / span
	w := r.newWriter(withBold)

	for x, i := 0, 1; i <= r.max; i++ {
		if i == r.middle {
			w.SetPosx(r.posX(i))
			for k := 0; k < int(times); k++ {
				base := x * 70
				// draw colored bone
				w.SetColorSeq(base)
				w.WriteString(r.saidBone)
				r.frames <- w.String()
				if k%magic == 0 {
					x++
				}
			}
		} else {
			base := x * 70
			w.SetPosx(r.posX(i))
			w.SetColorSeq(base)

			for _, line := range r.notSaidBoneLines {
				if i > r.screenWidth {
					// Left side animations
					n := i - r.screenWidth
					if n < line.Len() {
						w.WriteString(line.Slice(n, line.Len()))
					}
				} else if i <= line.Len() {
					// Right side animations
					w.WriteString(line.Slice(0, i-1))
				} else {
					w.WriteString(line.raw)
				}
				w.Write([]byte{'\n'})
			}
			r.frames <- w.String()
		}
		if i%magic == 0 {
			x++
		}
	}
	close(r.frames)
}

func (r *renderer) render() {
	initCh := make(chan struct{}, 1)
	initCh <- struct{}{}

	for view := range r.frames {
		select {
		case <-r.quit:
			screen.Clear()
			return
		case <-initCh:
		case <-time.After(span):
		}
		io.Copy(screen.Stdout, strings.NewReader(view))
	}
}

func (r *renderer) posX(i int) int {
	posx := r.screenWidth - i
	if posx < 1 {
		posx = 1
	}
	return posx
}

// Writer is wrapper which is both screen.MoveWriter and decoration.Writer.
type Writer struct {
	buf *strings.Builder
	mw  *screen.MoveWriter
	dw  *decoration.Writer
}

func (r *renderer) newWriter(withBold bool) *Writer {
	var buf strings.Builder
	mw := screen.NewMoveWriter(&buf, r.posX(0), r.heightDiff)
	options := []decoration.Option{
		decoration.WithAurora(0),
	}
	if withBold {
		options = append(options, decoration.WithBold())
	}
	dw := decoration.NewWriter(mw, options...)
	return &Writer{
		buf: &buf,
		mw:  mw,
		dw:  dw,
	}
}

// WriteString writes string. which is implemented io.StringWriter.
func (w *Writer) WriteString(s string) (int, error) { return w.dw.WriteString(s) }

// Write writes bytes. which is implemented io.Writer.
func (w *Writer) Write(p []byte) (int, error) { return w.dw.Write(p) }

// SetPosx sets posx
func (w *Writer) SetPosx(x int) { w.mw.SetPosx(x) }

// SetColorSeq sets color sequence.
func (w *Writer) SetColorSeq(colorSeq int) { w.dw.SetColorSeq(colorSeq) }

// Reset resets calls some Reset methods.
func (w *Writer) Reset() {
	w.buf.Reset()
	w.mw.Reset()
}

func (w *Writer) String() string {
	defer w.Reset()
	return w.buf.String()
}
