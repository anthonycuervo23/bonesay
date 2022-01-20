package bonesay

import (
	"fmt"
	"strings"

	wordwrap "github.com/Code-Hex/go-wordwrap"
	runewidth "github.com/mattn/go-runewidth"
)

type border struct {
	first  [2]rune
	middle [2]rune
	last   [2]rune
	only   [2]rune
}

func (bone *Bone) borderType() border {
	if bone.thinking {
		return border{
			first:  [2]rune{'(', ')'},
			middle: [2]rune{'(', ')'},
			last:   [2]rune{'(', ')'},
			only:   [2]rune{'(', ')'},
		}
	}

	return border{
		first:  [2]rune{'/', '\\'},
		middle: [2]rune{'|', '|'},
		last:   [2]rune{'\\', '/'},
		only:   [2]rune{'<', '>'},
	}
}

type line struct {
	text      string
	runeWidth int
}

type lines []*line

func (bone *Bone) maxLineWidth(lines []*line) int {
	maxWidth := 0
	for _, line := range lines {
		if line.runeWidth > maxWidth {
			maxWidth = line.runeWidth
		}
		if !bone.disableWordWrap && maxWidth > bone.ballonWidth {
			return bone.ballonWidth
		}
	}
	return maxWidth
}

func (bone *Bone) getLines(phrase string) []*line {
	text := bone.canonicalizePhrase(phrase)
	lineTexts := strings.Split(text, "\n")
	lines := make([]*line, 0, len(lineTexts))
	for _, lineText := range lineTexts {
		lines = append(lines, &line{
			text:      lineText,
			runeWidth: runewidth.StringWidth(lineText),
		})
	}
	return lines
}

func (bone *Bone) canonicalizePhrase(phrase string) string {
	// Replace tab to 8 spaces
	phrase = strings.Replace(phrase, "\t", "       ", -1)

	if bone.disableWordWrap {
		return phrase
	}
	width := bone.ballonWidth
	return wordwrap.WrapString(phrase, uint(width))
}

// Balloon to get the balloon and the string entered in the balloon.
func (bone *Bone) Balloon(phrase string) string {
	defer bone.buf.Reset()

	lines := bone.getLines(phrase)
	maxWidth := bone.maxLineWidth(lines)

	bone.writeBallon(lines, maxWidth)

	return bone.buf.String()
}

func (bone *Bone) writeBallon(lines []*line, maxWidth int) {
	top := make([]byte, 0)
	bottom := make([]byte, 0)

	for i := 0; i < bone.balloonOffset; i++ {
		top = append(top, ' ')
		bottom = append(bottom, ' ')
	}

	for i := 0; i < maxWidth+2; i++ {
		top = append(top, '_')
		bottom = append(bottom, '-')
	}

	borderType := bone.borderType()

	bone.buf.Write(top)
	bone.buf.Write([]byte{' ', '\n'})
	defer func() {
		bone.buf.Write(bottom)
		bone.buf.Write([]byte{' ', '\n'})
	}()

	l := len(lines)
	if l == 1 {
		border := borderType.only
		for i := 0; i < (bone.balloonOffset - 1); i++ {
			bone.buf.WriteRune(' ')
		}
		bone.buf.WriteRune(border[0])
		bone.buf.WriteRune(' ')
		bone.buf.WriteString(lines[0].text)
		bone.buf.WriteRune(' ')
		bone.buf.WriteRune(border[1])
		bone.buf.WriteRune('\n')
		return
	}

	var border [2]rune
	for i := 0; i < l; i++ {
		switch i {
		case 0:
			border = borderType.first
		case l - 1:
			border = borderType.last
		default:
			border = borderType.middle
		}
		for i := 0; i < (bone.balloonOffset - 1); i++ {
			bone.buf.WriteRune(' ')
		}
		bone.buf.WriteRune(border[0])
		bone.buf.WriteRune(' ')
		bone.padding(lines[i], maxWidth)
		bone.buf.WriteRune(' ')
		bone.buf.WriteRune(border[1])
		bone.buf.WriteRune('\n')
	}
}

func (bone *Bone) flush(text, top, bottom fmt.Stringer) string {
	return fmt.Sprintf(
		"%s\n%s%s\n",
		top.String(),
		text.String(),
		bottom.String(),
	)
}

func (bone *Bone) padding(line *line, maxWidth int) {
	if maxWidth <= line.runeWidth {
		bone.buf.WriteString(line.text)
		return
	}

	bone.buf.WriteString(line.text)
	l := maxWidth - line.runeWidth
	for i := 0; i < l; i++ {
		bone.buf.WriteRune(' ')
	}
}
