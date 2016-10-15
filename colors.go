package logging

import (
	"strconv"
	"errors"
	"fmt"
	"strings"
)

func colorResetFrag(ctx *outputContext) {
	ctx.dst.Write([]byte("\033[0m"))
}

var ErrTooManyColors = errors.New("too many colors mentioned")

func makeFixedColorFrag(options string) (fragmentFormatter, error) {
	color, err := parseColor(options)
	if err != nil {
		return nil, err
	}
	return literals(color), nil
}

func makeColorFrag(options string) (fragmentFormatter, error) {
	var palette []string
	if options != "" {
		if options[0] == '=' {
			return makeFixedColorFrag(options[1:])
		}
		
		palette = strings.Split(options, ",")
		if len(palette) > 8 {
			return nil, ErrTooManyColors
		}
	}
	
	var colors [8][]byte
	for i := 0; i < 8; i++ {
		color := defaultColorPalette[i]
		if i < len(palette) {
			var err error
			color, err = parseColor(palette[i])
			if err != nil {
				return nil, err
			}
		}
		colors[i] = color
	}

	return func(ctx *outputContext) {
		ctx.dst.Write(colors[ctx.src.Level])
	}, nil

}

type ErrInvalidColorSpec struct {
	spec string
}

func (err ErrInvalidColorSpec) Error() string {
	return fmt.Sprintf("invalid color spec %q", err.spec)
}

func parseColor(spec string) ([]byte, error) {
	bg := false
	faint := false
	bold := false

	choice := color(0)

	for _, ch := range spec {
		switch ch {
		case '_':
			faint = true
		case '*':
			bold = true
		case '#':
			bg = true
		case 'k':
			choice = colorBlack
		case 'r':
			choice = colorRed
		case 'g':
			choice = colorGreen
		case 'y':
			choice = colorYellow
		case 'b':
			choice = colorBlue
		case 'm':
			choice = colorMagenta
		case 'c':
			choice = colorCyan
		case 'w':
			choice = colorWhite
		case ' ':
		default:
			return nil, ErrInvalidColorSpec{spec}
		}
	}
	if choice == color(0) {
		return nil, ErrInvalidColorSpec{spec}
	}

	c := []byte{'\033', '['}

	if bg {
		choice += 10
	}
	c = strconv.AppendUint(c, uint64(choice), 10)
	if bold {
		c = append(c, ';', '1')
	}
	if faint {
		c = append(c, ';', '2')
	}
	c = append(c, 'm')
	return c, nil
}

type color int

const (
	colorBlack = (iota + 30)
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

func colorSeq(color color) []byte {
	return []byte(fmt.Sprintf("\033[%dm", int(color)))
}

func colorSeqBold(color color) []byte {
	return []byte(fmt.Sprintf("\033[%d;1m", int(color)))
}

func colorSeqBg(color color) []byte {
	return []byte(fmt.Sprintf("\033[%dm", int(10+color)))
}

func colorSeqFaint(color color) []byte {
	return []byte(fmt.Sprintf("\033[%d;2m", int(color)))
}

var defaultColorPalette = [8][]byte{
	EMERGENCY: colorSeqBg(colorMagenta),
	ALERT:     colorSeqBg(colorRed),
	CRITICAL:  colorSeqBold(colorMagenta),
	ERROR:     colorSeqBold(colorRed),
	WARNING:   colorSeqBold(colorYellow),
	NOTICE:    colorSeqBold(colorGreen),
	INFO:      colorSeq(colorGreen),
	DEBUG:     colorSeq(colorWhite),
}

func makeColorResetFrag(_ string) (fragmentFormatter, error) {
	return colorResetFrag, nil
}
