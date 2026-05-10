package main

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	argparse "github.com/rsa17826/go-arg-lib"
	input "github.com/rsa17826/go-input-lib"
)

// ── colours (matches the AHK dark theme) ─────────────────────────────────────
var (
	colBG      = color.NRGBA{0x2A, 0x2A, 0x2E, 0xFF}
	colKey     = color.NRGBA{0x01, 0x04, 0x09, 0xFF}
	colText    = color.NRGBA{0x8B, 0x94, 0x9E, 0xFF}
	colPressed = color.NRGBA{0x55, 0x3B, 0x6B, 0xFF}
	colToggled = color.NRGBA{0x55, 0x3B, 0x6A, 0xFF}
)

// ── layout constants ──────────────────────────────────────────────────────────
const (
	keyUnit  = float32(22)
	keyGap   = float32(2)
	keyH     = float32(20)
	rowGap   = float32(2)
	fontSize = float32(7)
)

// ── linux key codes ───────────────────────────────────────────────────────────
const (
	KEY_ESC        = 1
	KEY_1          = 2
	KEY_2          = 3
	KEY_3          = 4
	KEY_4          = 5
	KEY_5          = 6
	KEY_6          = 7
	KEY_7          = 8
	KEY_8          = 9
	KEY_9          = 10
	KEY_0          = 11
	KEY_MINUS      = 12
	KEY_EQUAL      = 13
	KEY_BACKSPACE  = 14
	KEY_TAB        = 15
	KEY_Q          = 16
	KEY_W          = 17
	KEY_E          = 18
	KEY_R          = 19
	KEY_T          = 20
	KEY_Y          = 21
	KEY_U          = 22
	KEY_I          = 23
	KEY_O          = 24
	KEY_P          = 25
	KEY_LEFTBRACE  = 26
	KEY_RIGHTBRACE = 27
	KEY_ENTER      = 28
	KEY_LEFTCTRL   = 29
	KEY_A          = 30
	KEY_S          = 31
	KEY_D          = 32
	KEY_F          = 33
	KEY_G          = 34
	KEY_H          = 35
	KEY_J          = 36
	KEY_K          = 37
	KEY_L          = 38
	KEY_SEMICOLON  = 39
	KEY_APOSTROPHE = 40
	KEY_GRAVE      = 41
	KEY_LEFTSHIFT  = 42
	KEY_BACKSLASH  = 43
	KEY_Z          = 44
	KEY_X          = 45
	KEY_C          = 46
	KEY_V          = 47
	KEY_B          = 48
	KEY_N          = 49
	KEY_M          = 50
	KEY_COMMA      = 51
	KEY_DOT        = 52
	KEY_SLASH      = 53
	KEY_RIGHTSHIFT = 54
	KEY_KPASTERISK = 55
	KEY_LEFTALT    = 56
	KEY_SPACE      = 57
	KEY_CAPSLOCK   = 58
	KEY_F1         = 59
	KEY_F2         = 60
	KEY_F3         = 61
	KEY_F4         = 62
	KEY_F5         = 63
	KEY_F6         = 64
	KEY_F7         = 65
	KEY_F8         = 66
	KEY_F9         = 67
	KEY_F10        = 68
	KEY_NUMLOCK    = 69
	KEY_SCROLLLOCK = 70
	KEY_KP7        = 71
	KEY_KP8        = 72
	KEY_KP9        = 73
	KEY_KPMINUS    = 74
	KEY_KP4        = 75
	KEY_KP5        = 76
	KEY_KP6        = 77
	KEY_KPPLUS     = 78
	KEY_KP1        = 79
	KEY_KP2        = 80
	KEY_KP3        = 81
	KEY_KP0        = 82
	KEY_KPDOT      = 83
	KEY_F11        = 87
	KEY_F12        = 88
	KEY_KPENTER    = 96
	KEY_RIGHTCTRL  = 97
	KEY_KPSLASH    = 98
	KEY_SYSRQ      = 99
	KEY_RIGHTALT   = 100
	KEY_HOME       = 102
	KEY_UP         = 103
	KEY_PAGEUP     = 104
	KEY_LEFT       = 105
	KEY_RIGHT      = 106
	KEY_END        = 107
	KEY_DOWN       = 108
	KEY_PAGEDOWN   = 109
	KEY_INSERT     = 110
	KEY_DELETE     = 111
	KEY_PAUSE      = 119
	KEY_LEFTMETA   = 125
	KEY_RIGHTMETA  = 126
	KEY_COMPOSE    = 127
)

// toggle-lock keys (highlight while state is on, not while held)
var toggleKeys = map[uint16]bool{
	KEY_CAPSLOCK:   true,
	KEY_NUMLOCK:    true,
	KEY_SCROLLLOCK: true,
}

// ── layout definition ─────────────────────────────────────────────────────────
// Each entry: {keycode, label, widthUnits, extraGapBefore}
type KeyDef struct {
	Code  uint16
	Label string
	W     float32 // width in key units (default 1)
	Gap   float32 // extra x-gap before this key in key units
}

var rows = [][]KeyDef{
	{ // F-row
		{KEY_ESC, "Esc", 1, 0},
		{KEY_F1, "F1", 1, 0.5}, {KEY_F2, "F2", 1, 0}, {KEY_F3, "F3", 1, 0}, {KEY_F4, "F4", 1, 0},
		{KEY_F5, "F5", 1, 0.5}, {KEY_F6, "F6", 1, 0}, {KEY_F7, "F7", 1, 0}, {KEY_F8, "F8", 1, 0},
		{KEY_F9, "F9", 1, 0.5}, {KEY_F10, "F10", 1, 0}, {KEY_F11, "F11", 1, 0}, {KEY_F12, "F12", 1, 0},
		{KEY_SYSRQ, "Prt\nScr", 1, 1}, {KEY_SCROLLLOCK, "Scr\nLk", 1, 0}, {KEY_PAUSE, "Pause", 1, 0},
	},
	{ // number row
		{KEY_GRAVE, "`", 1, 0},
		{KEY_1, "1", 1, 0}, {KEY_2, "2", 1, 0}, {KEY_3, "3", 1, 0}, {KEY_4, "4", 1, 0},
		{KEY_5, "5", 1, 0}, {KEY_6, "6", 1, 0}, {KEY_7, "7", 1, 0}, {KEY_8, "8", 1, 0},
		{KEY_9, "9", 1, 0}, {KEY_0, "0", 1, 0}, {KEY_MINUS, "-", 1, 0}, {KEY_EQUAL, "=", 1, 0},
		{KEY_BACKSPACE, "BS", 2, 0},
		{KEY_INSERT, "Ins", 1, 0.5}, {KEY_HOME, "Home", 1, 0}, {KEY_PAGEUP, "PgUp", 1, 0},
		// numpad
		{KEY_NUMLOCK, "Num\nLk", 1, 0.5}, {KEY_KPSLASH, "/", 1, 0}, {KEY_KPASTERISK, "*", 1, 0}, {KEY_KPMINUS, "-", 1, 0},
	},
	{ // tab row
		{KEY_TAB, "Tab", 1.5, 0},
		{KEY_Q, "Q", 1, 0}, {KEY_W, "W", 1, 0}, {KEY_E, "E", 1, 0}, {KEY_R, "R", 1, 0},
		{KEY_T, "T", 1, 0}, {KEY_Y, "Y", 1, 0}, {KEY_U, "U", 1, 0}, {KEY_I, "I", 1, 0},
		{KEY_O, "O", 1, 0}, {KEY_P, "P", 1, 0}, {KEY_LEFTBRACE, "[", 1, 0}, {KEY_RIGHTBRACE, "]", 1, 0},
		{KEY_BACKSLASH, "\\", 1.5, 0},
		{KEY_DELETE, "Del", 1, 0.5}, {KEY_END, "End", 1, 0}, {KEY_PAGEDOWN, "PgDn", 1, 0},
		// numpad
		{KEY_KP7, "7", 1, 0.5}, {KEY_KP8, "8", 1, 0}, {KEY_KP9, "9", 1, 0}, {KEY_KPPLUS, "+", 1, 0},
	},
	{ // caps row
		{KEY_CAPSLOCK, "Caps", 1.75, 0},
		{KEY_A, "A", 1, 0}, {KEY_S, "S", 1, 0}, {KEY_D, "D", 1, 0}, {KEY_F, "F", 1, 0},
		{KEY_G, "G", 1, 0}, {KEY_H, "H", 1, 0}, {KEY_J, "J", 1, 0}, {KEY_K, "K", 1, 0},
		{KEY_L, "L", 1, 0}, {KEY_SEMICOLON, ";", 1, 0}, {KEY_APOSTROPHE, "'", 1, 0},
		{KEY_ENTER, "Enter", 2.25, 0},
		// (nav cluster empty this row)
		// numpad
		{KEY_KP4, "4", 1, 4}, {KEY_KP5, "5", 1, 0}, {KEY_KP6, "6", 1, 0}, {0, " ", 1, 0}, // +tall is on row above
	},
	{ // shift row
		{KEY_LEFTSHIFT, "Shift", 2.25, 0},
		{KEY_Z, "Z", 1, 0}, {KEY_X, "X", 1, 0}, {KEY_C, "C", 1, 0}, {KEY_V, "V", 1, 0},
		{KEY_B, "B", 1, 0}, {KEY_N, "N", 1, 0}, {KEY_M, "M", 1, 0},
		{KEY_COMMA, ",", 1, 0}, {KEY_DOT, ".", 1, 0}, {KEY_SLASH, "/", 1, 0},
		{KEY_RIGHTSHIFT, "Shift", 2.75, 0},
		{KEY_UP, "↑", 1, 1.5},
		// numpad
		{KEY_KP1, "1", 1, 1.5}, {KEY_KP2, "2", 1, 0}, {KEY_KP3, "3", 1, 0}, {KEY_KPENTER, "Ent", 1, 0},
	},
	{ // ctrl row
		{KEY_LEFTCTRL, "Ctrl", 1.5, 0}, {KEY_LEFTMETA, "Win", 1, 0}, {KEY_LEFTALT, "Alt", 1.5, 0},
		{KEY_SPACE, "Space", 6.00, 0},
		{KEY_RIGHTALT, "Alt", 1.5, 0}, {KEY_RIGHTMETA, "Win", 1, 0}, {KEY_COMPOSE, "App", 1, 0}, {KEY_RIGHTCTRL, "Ctrl", 1.5, 0},
		{KEY_LEFT, "←", 1, 0.5}, {KEY_DOWN, "↓", 1, 0}, {KEY_RIGHT, "→", 1, 0},
		// numpad
		{KEY_KP0, "0", 2, 0.5}, {KEY_KPDOT, ".", 1, 0}, {0, " ", 1, 0}, // enter tall handled above
	},
}

// ── key widget ────────────────────────────────────────────────────────────────
type KeyWidget struct {
	rect    *canvas.Rectangle
	toggle  bool // is this a toggle-lock key?
	toggled bool // current toggle state
}

var keyMap = map[uint16]*KeyWidget{}

func setKeyColour(kw *KeyWidget, c color.NRGBA) {
	kw.rect.FillColor = c
	fyne.Do(func() {
		canvas.Refresh(kw.rect)
	})
}

func pressKey(code uint16) {
	kw, ok := keyMap[code]
	if !ok {
		return
	}
	if kw.toggle {
		kw.toggled = !kw.toggled
		if kw.toggled {
			setKeyColour(kw, colToggled)
		} else {
			setKeyColour(kw, colKey)
		}
	} else {
		setKeyColour(kw, colPressed)
	}
}

func releaseKey(code uint16) {
	kw, ok := keyMap[code]
	if !ok {
		return
	}
	if !kw.toggle {
		setKeyColour(kw, colKey)
	}
}

// ── build the keyboard canvas ─────────────────────────────────────────────────
func buildKeyboard() (fyne.CanvasObject, fyne.Size) {
	var rects []fyne.CanvasObject
	var texts []fyne.CanvasObject

	maxX := float32(0)
	y := float32(0)

	for _, row := range rows {
		x := float32(0)
		for _, k := range row {
			x += k.Gap * keyUnit

			w := k.W*keyUnit - keyGap

			rect := canvas.NewRectangle(colKey)
			rect.CornerRadius = 4
			rect.Resize(fyne.NewSize(w, keyH))
			rect.Move(fyne.NewPos(x, y))
			rects = append(rects, rect)

			// text label - split on \n for two-line keys
			lines := strings.Split(k.Label, "\n")
			lineH := keyH / float32(len(lines)+1)
			for li, line := range lines {
				t := canvas.NewText(line, colText)
				t.TextSize = fontSize
				t.Alignment = fyne.TextAlignCenter
				t.Resize(fyne.NewSize(w, lineH))
				t.Move(fyne.NewPos(x, y+lineH*float32(li+1)-lineH*0.6))
				texts = append(texts, t)
			}

			if k.Code != 0 {
				keyMap[k.Code] = &KeyWidget{
					rect:   rect,
					toggle: toggleKeys[k.Code],
				}
			}

			x += k.W * keyUnit
			if x > maxX {
				maxX = x
			}
		}
		y += keyH + rowGap
	}

	// rects first so texts render on top
	all := append(rects, texts...)
	return container.NewWithoutLayout(all...), fyne.NewSize(maxX+10, y+10)
}

// ── input loop ────────────────────────────────────────────────────────────────
func monitorInput(devicePath string) {
	f, err := os.Open(devicePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open device: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var ev input.InputEvent
	for {
		if err := binary.Read(f, binary.NativeEndian, &ev); err != nil {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			return
		}
		if ev.Type != input.EV_KEY {
			continue
		}
		switch ev.Value {
		case 1: // key down
			pressKey(ev.Code)
		case 0: // key up
			releaseKey(ev.Code)
			// value 2 = repeat - ignore
		}
	}
}

// ── dark theme ────────────────────────────────────────────────────────────────
type darkTheme struct{ fyne.Theme }

func (darkTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch n {
	case theme.ColorNameBackground:
		return colBG
	}
	return theme.DefaultTheme().Color(n, theme.VariantDark)
}

// ── main ──────────────────────────────────────────────────────────────────────
func main() {
	var deviceArg string
	var detect bool

	argData := []argparse.ArgumentData{
		{Keys: []string{"device", "d"}, AfterCount: 1, Target: &deviceArg, Description: "the device to display data from"},
		{Keys: []string{"detect"}, AfterCount: 0, Target: &detect, Description: "identify the device to use"},
	}
	argparse.ParseArgs(argData)
	if detect {
		input.GetDeviceToUser()
		return
	}

	if deviceArg == "" {
		argparse.PrintHelp(argData)
	}

	devicePath := input.WaitForDevice(deviceArg)
	fmt.Println("using device:", devicePath)

	a := app.New()
	a.Settings().SetTheme(darkTheme{theme.DefaultTheme()})

	w := a.NewWindow("Key Display")
	w.SetFixedSize(true)

	kb, size := buildKeyboard()
	w.Resize(size)
	w.SetContent(kb)

	go monitorInput(devicePath)

	w.ShowAndRun()
}
