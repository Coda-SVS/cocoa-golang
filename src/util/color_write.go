package util

// https://github.com/xiegeo/coloredgoroutine 의 코드를 수정하여 사용

import (
	"sort"

	"github.com/fatih/color"
)

type colorDetail struct {
	color.Color
	fg color.Attribute
	bg color.Attribute
}

var (
	colors       []colorDetail
	bannedColors []colorDetail
)

func init() {
	add := func(f, b color.Attribute) {
		if banned(f, b) {
			bannedColors = append(bannedColors, newColorDetail(f, b))
		} else {
			colors = append(colors, newColorDetail(f, b))
		}
	}
	hi := color.FgHiBlack - color.FgBlack
	for f := color.FgBlack; f <= color.FgWhite; f++ {
		for b := color.BgBlack; b <= color.BgWhite; b++ {
			add(f, b)
			add(f+hi, b)
			add(f, b+hi)
			add(f+hi, b+hi)
		}
	}
	shuffle(colors)
	order(bannedColors)
}

func GetColorForID(id int) *color.Color {
	id = id % len(colors)
	if id < 0 {
		id += len(colors)
	}
	return &colors[id].Color
}

func newColorDetail(fg, bg color.Attribute) colorDetail {
	return colorDetail{
		Color: *color.New(fg, bg),
		fg:    fg,
		bg:    bg,
	}
}

func banned(f, b color.Attribute) bool {
	same := color.BgBlack - color.FgBlack
	if b-f == same {
		return true
	}

	b = b - same
	if b > f {
		b, f = f, b
	}

	switch b {
	case color.FgGreen:
		return f == color.FgYellow || f == color.FgHiBlue || f == color.FgHiMagenta || f == color.FgHiBlack
	case color.FgYellow:
		return f == color.FgWhite || f == color.FgHiGreen || f == color.FgHiMagenta || f == color.FgHiCyan
	case color.FgBlue:
		return f == color.FgCyan
	case color.FgMagenta:
		return f == color.FgHiBlack || f == color.FgHiRed || f == color.FgCyan
	case color.FgCyan:
		return f == color.FgHiBlack || f == color.FgHiBlue
	case color.FgWhite:
		return f == color.FgHiGreen || f == color.FgHiYellow || f == color.FgHiCyan
	case color.FgHiBlack:
		return f == color.FgHiBlue || f == color.FgHiMagenta || f == color.FgHiRed
	case color.FgHiGreen:
		return f == color.FgHiYellow || f == color.FgHiCyan
	case color.FgHiYellow:
		return f == color.FgHiCyan
	case color.FgHiBlue:
		return f == color.FgHiRed || f == color.FgHiMagenta
	}

	return false
}

func shuffle(c []colorDetail) {
	shuffleKey := 11
	for i := 0; i < len(c); i++ {
		t := (i * shuffleKey) % len(c)
		c[i], c[t] = c[t], c[i]
	}
}

func order(c []colorDetail) {
	same := color.BgBlack - color.FgBlack
	sort.Slice(c, func(i, j int) bool {
		ai, bi := c[i].fg, c[i].bg-same
		if ai > bi {
			ai, bi = bi, ai
		}
		aj, bj := c[j].fg, c[j].bg-same
		if aj > bj {
			aj, bj = bj, aj
		}
		if ai < aj {
			return true
		}
		if ai > aj {
			return false
		}
		return bi < bj
	})
}
