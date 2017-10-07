package tui

import (
	"image"
	"strings"

	runewidth "github.com/mattn/go-runewidth"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

type RuneBuffer struct {
	buf []rune
	idx int

	wordwrap bool
}

// Width returns the width of the rune buffer, taking into account for CJK.
func (r *RuneBuffer) Width() int {
	return runewidth.StringWidth(string(r.buf))
}

// Set the buffer and the index at the end of the buffer.
func (b *RuneBuffer) Set(buf []rune) {
	b.SetWithIdx(len(buf), buf)
}

// SetWithIdx set the the buffer with a given index.
func (b *RuneBuffer) SetWithIdx(idx int, buf []rune) {
	b.buf = buf
	b.idx = idx
}

// WriteRune appends a rune to the buffer.
func (r *RuneBuffer) WriteRune(s rune) {
	r.WriteRunes([]rune{s})
}

// WriteRunes appends runes to the buffer.
func (r *RuneBuffer) WriteRunes(s []rune) {
	tail := append(s, r.buf[r.idx:]...)
	r.buf = append(r.buf[:r.idx], tail...)
	r.idx += len(s)
}

// Pos returns the current index in the buffer.
func (r *RuneBuffer) Pos() int {
	return r.idx
}

// Len returns the number of runes in the buffer.
func (r *RuneBuffer) Len() int {
	return len(r.buf)
}

func (r *RuneBuffer) SplitByLine(width int) []string {
	var text string
	if r.wordwrap {
		text = wordwrap.WrapString(r.String(), uint(width))
	} else {
		text = r.String()
	}
	return strings.Split(text, "\n")
}

func getSplitByLine(rs []rune, width int, wrap bool) []string {
	var text string
	if wrap {
		text = wordwrap.WrapString(string(rs), uint(width))
	} else {
		text = string(rs)
	}
	return strings.Split(text, "\n")
}

func (r *RuneBuffer) CursorPos(width int) image.Point {
	if width == 0 {
		return image.ZP
	}

	sp := getSplitByLine(r.buf[:r.idx], width, r.wordwrap)

	return image.Pt(stringWidth(sp[len(sp)-1]), len(sp))
}

func (b *RuneBuffer) String() string {
	return string(b.buf)
}

func (r *RuneBuffer) MoveBackward() {
	if r.idx == 0 {
		return
	}
	r.idx--
}

func (r *RuneBuffer) MoveForward() {
	if r.idx == len(r.buf) {
		return
	}
	r.idx++
}

func (r *RuneBuffer) MoveToLineStart() {
	for i := r.idx; i > 0; i-- {
		if r.buf[i-1] == '\n' {
			r.idx = i
			return
		}
	}
	r.idx = 0
}

func (r *RuneBuffer) MoveToLineEnd() {
	for i := r.idx; i < len(r.buf)-1; i++ {
		if r.buf[i+1] == '\n' {
			r.idx = i
			return
		}
	}
	r.idx = len(r.buf)
}

func (r *RuneBuffer) Backspace() {
	if r.idx == 0 {
		return
	}
	r.idx--
	r.buf = append(r.buf[:r.idx], r.buf[r.idx+1:]...)
}

func (r *RuneBuffer) Delete() {
	if r.idx == len(r.buf) {
		return
	}
	r.buf = append(r.buf[:r.idx], r.buf[r.idx+1:]...)
}

func (r *RuneBuffer) Kill() {
	r.buf = r.buf[:r.idx]
}

func IsWordBreak(i rune) bool {
	switch {
	case i >= 'a' && i <= 'z':
	case i >= 'A' && i <= 'Z':
	case i >= '0' && i <= '9':
	default:
		return true
	}
	return false
}
