package tempura

import (
	"image/color"

	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

// Align is used to specify text alignment
type Align uint8

const (
	// AlignLeft specified to draw the left edge of text at the specified x coordinate
	AlignLeft Align = iota
	// AlignRight specifies to draw the text horizontally centered at the specified x coordinate
	AlignCenter
	// AlignRight specifies to draw the right edge of text at the specified x coordinate
	AlignRight
)

type Text struct {
	Face    font.Face
	Color   color.Color
	Text    string
	W       int
	H       int
	Advance int
}

func NewText(face font.Face, color color.Color, text string) Text {
	w, h, advance := MeasureText(text, face)
	return Text{
		Face:    face,
		Color:   color,
		Text:    text,
		W:       w,
		H:       h,
		Advance: advance,
	}
}

func NewTextf(face font.Face, color color.Color, text string, args ...interface{}) Text {
	if len(args) == 0 {
		return NewText(face, color, text)
	}
	return NewText(face, color, fmt.Sprintf(text, args...))
}

func (t Text) Draw(image *ebiten.Image, x, y int, align Align) {
	// TODO(explodes): do not draw if text is not in bounds of image
	switch align {
	case AlignCenter:
		x = x - t.W/2
	case AlignRight:
		x = x - t.W
	}
	text.Draw(image, t.Text, t.Face, x, y, t.Color)
}

type Texts []Text

func NewTexts(face font.Face, color color.Color, texts []string) Texts {
	measured := make(Texts, len(texts))
	for i, t := range texts {
		measured[i] = NewText(face, color, t)
	}
	return measured
}

func (ts *Texts) Push(face font.Face, color color.Color, text string) {
	next := append(*ts, NewText(face, color, text))
	*ts = next
}

func (ts *Texts) Pushf(face font.Face, color color.Color, text string, args ...interface{}) {
	next := append(*ts, NewTextf(face, color, text, args...))
	*ts = next
}

func (ts Texts) DrawSingleLine(image *ebiten.Image, x, y int, align Align) {
	switch align {
	case AlignLeft:
		for _, t := range ts {
			t.Draw(image, x, y, AlignLeft)
			x += t.Advance
		}
	case AlignCenter:
		width := ts.SingleLineWidth()
		for _, t := range ts {
			t.Draw(image, x-width/2, y, AlignLeft)
			x += t.Advance
		}
	case AlignRight:
		width := ts.SingleLineWidth()
		for _, t := range ts {
			t.Draw(image, x-width, y, AlignLeft)
			x += t.Advance
		}
	}
}

func (ts Texts) SingleLineWidth() int {
	width := 0
	for _, t := range ts {
		width += t.Advance
	}
	return width
}

func (ts Texts) SingleLineHeight() int {
	height := 0
	for _, t := range ts {
		if t.H > height {
			height = t.H
		}
	}
	return height
}

func (ts Texts) MultiLineWidth() int {
	width := 0
	for _, t := range ts {
		if t.Advance > width {
			width = t.Advance
		}
	}
	return width
}

func (ts Texts) MultiLineHeight(addSpace int) int {
	height := 0
	for _, t := range ts {
		height += t.H + addSpace
	}
	return height
}

func (ts Texts) DrawLines(image *ebiten.Image, addSpace, x, y int, align Align) {
	switch align {
	case AlignLeft:
		for _, t := range ts {
			t.Draw(image, x, y, AlignLeft)
			y += t.H + addSpace
		}
	case AlignCenter:
		for _, t := range ts {
			t.Draw(image, x-t.W/2, y, AlignLeft)
			y += t.H + addSpace
		}
	case AlignRight:
		for _, t := range ts {
			t.Draw(image, x-t.W, y, AlignLeft)
			y += t.H + addSpace
		}
	}
}

// MeasureText measures the size of a string with a given font face.
// This is probably best computed the fewest amount of times
// necessary, use MeasuredText for convenience.
func MeasureText(text string, face font.Face) (width, height, advance int) {
	bounds, a := font.BoundString(face, text)
	width = (bounds.Max.X - bounds.Min.X).Ceil()
	height = (bounds.Max.Y - bounds.Min.Y).Ceil()
	advance = a.Ceil()
	return
}
