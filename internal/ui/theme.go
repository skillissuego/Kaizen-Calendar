package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	ThemeTurquoiseDark  = "turquoise-dark"
	ThemeGraphite       = "graphite"
	ThemePurpleNight    = "purple-night"
	ThemeTurquoiseLight = "turquoise-light"
)

type ThemeOption struct {
	Name  string
	Title string
}

var themeOptions = []ThemeOption{
	{Name: ThemeTurquoiseDark, Title: "Тёмная бирюза"},
	{Name: ThemeGraphite, Title: "Графит"},
	{Name: ThemePurpleNight, Title: "Фиолетовая ночь"},
	{Name: ThemeTurquoiseLight, Title: "Светлая бирюза"},
}

func ThemeOptions() []ThemeOption {
	return append([]ThemeOption(nil), themeOptions...)
}

func ThemeTitle(name string) string {
	for _, option := range themeOptions {
		if option.Name == name {
			return option.Title
		}
	}
	return ThemeTitle(ThemeTurquoiseDark)
}

// KaizenTheme provides several deliberately restrained palettes for the app.
type KaizenTheme struct {
	base    fyne.Theme
	palette themePalette
}

type themePalette struct {
	background        color.NRGBA
	foreground        color.NRGBA
	button            color.NRGBA
	hover             color.NRGBA
	focus             color.NRGBA
	border            color.NRGBA
	disabled          color.NRGBA
	inputBackground   color.NRGBA
	overlayBackground color.NRGBA
	foregroundPrimary color.NRGBA
	placeholder       color.NRGBA
	separator         color.NRGBA
	selection         color.NRGBA
	pressed           color.NRGBA
	menuBackground    color.NRGBA
}

func NewKaizenTheme(name string) fyne.Theme {
	return &KaizenTheme{base: theme.DefaultTheme(), palette: paletteFor(name)}
}

func (t *KaizenTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return t.palette.background
	case theme.ColorNameForeground:
		return t.palette.foreground
	case theme.ColorNameButton, theme.ColorNamePrimary:
		return t.palette.button
	case theme.ColorNameHover:
		return t.palette.hover
	case theme.ColorNameFocus:
		return t.palette.focus
	case theme.ColorNameInputBorder:
		return t.palette.border
	case theme.ColorNameDisabled:
		return t.palette.disabled
	case theme.ColorNameInputBackground:
		return t.palette.inputBackground
	case theme.ColorNameOverlayBackground:
		return t.palette.overlayBackground
	case theme.ColorNameForegroundOnPrimary:
		return t.palette.foregroundPrimary
	case theme.ColorNamePlaceHolder:
		return t.palette.placeholder
	case theme.ColorNameSeparator:
		return t.palette.separator
	case theme.ColorNameSelection:
		return t.palette.selection
	case theme.ColorNamePressed:
		return t.palette.pressed
	case theme.ColorNameMenuBackground:
		return t.palette.menuBackground
	default:
		return t.base.Color(name, variant)
	}
}

func paletteFor(name string) themePalette {
	switch name {
	case ThemeGraphite:
		return themePalette{
			background: color.NRGBA{R: 25, G: 27, B: 31, A: 255}, foreground: color.NRGBA{R: 235, G: 237, B: 240, A: 255},
			button: color.NRGBA{R: 91, G: 101, B: 116, A: 255}, hover: color.NRGBA{R: 116, G: 128, B: 145, A: 255}, focus: color.NRGBA{R: 166, G: 181, B: 201, A: 255},
			border: color.NRGBA{R: 85, G: 91, B: 103, A: 255}, disabled: color.NRGBA{R: 126, G: 130, B: 138, A: 255},
			inputBackground: color.NRGBA{R: 35, G: 38, B: 44, A: 255}, overlayBackground: color.NRGBA{R: 29, G: 31, B: 36, A: 255}, foregroundPrimary: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			placeholder: color.NRGBA{R: 160, G: 165, B: 174, A: 255}, separator: color.NRGBA{R: 63, G: 68, B: 78, A: 255}, selection: color.NRGBA{R: 70, G: 86, B: 108, A: 255}, pressed: color.NRGBA{R: 75, G: 84, B: 98, A: 255}, menuBackground: color.NRGBA{R: 35, G: 38, B: 44, A: 255},
		}
	case ThemePurpleNight:
		return themePalette{
			background: color.NRGBA{R: 31, G: 24, B: 42, A: 255}, foreground: color.NRGBA{R: 244, G: 237, B: 250, A: 255},
			button: color.NRGBA{R: 126, G: 83, B: 166, A: 255}, hover: color.NRGBA{R: 157, G: 105, B: 202, A: 255}, focus: color.NRGBA{R: 204, G: 153, B: 238, A: 255},
			border: color.NRGBA{R: 105, G: 77, B: 130, A: 255}, disabled: color.NRGBA{R: 143, G: 128, B: 152, A: 255},
			inputBackground: color.NRGBA{R: 47, G: 36, B: 62, A: 255}, overlayBackground: color.NRGBA{R: 38, G: 29, B: 51, A: 255}, foregroundPrimary: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			placeholder: color.NRGBA{R: 177, G: 157, B: 192, A: 255}, separator: color.NRGBA{R: 79, G: 59, B: 101, A: 255}, selection: color.NRGBA{R: 101, G: 70, B: 132, A: 255}, pressed: color.NRGBA{R: 110, G: 75, B: 145, A: 255}, menuBackground: color.NRGBA{R: 47, G: 36, B: 62, A: 255},
		}
	case ThemeTurquoiseLight:
		return themePalette{
			background: color.NRGBA{R: 235, G: 248, B: 247, A: 255}, foreground: color.NRGBA{R: 24, G: 52, B: 54, A: 255},
			button: color.NRGBA{R: 25, G: 145, B: 140, A: 255}, hover: color.NRGBA{R: 32, G: 174, B: 166, A: 255}, focus: color.NRGBA{R: 17, G: 111, B: 108, A: 255},
			border: color.NRGBA{R: 122, G: 178, B: 175, A: 255}, disabled: color.NRGBA{R: 130, G: 157, B: 157, A: 255},
			inputBackground: color.NRGBA{R: 255, G: 255, B: 255, A: 255}, overlayBackground: color.NRGBA{R: 248, G: 253, B: 252, A: 255}, foregroundPrimary: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			placeholder: color.NRGBA{R: 91, G: 126, B: 127, A: 255}, separator: color.NRGBA{R: 185, G: 216, B: 214, A: 255}, selection: color.NRGBA{R: 177, G: 229, B: 224, A: 255}, pressed: color.NRGBA{R: 19, G: 125, B: 120, A: 255}, menuBackground: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		}
	default:
		return themePalette{
			background: color.NRGBA{R: 18, G: 28, B: 32, A: 255}, foreground: color.NRGBA{R: 232, G: 245, B: 245, A: 255},
			button: color.NRGBA{R: 26, G: 148, B: 143, A: 255}, hover: color.NRGBA{R: 35, G: 181, B: 174, A: 255}, focus: color.NRGBA{R: 71, G: 214, B: 204, A: 255},
			border: color.NRGBA{R: 66, G: 101, B: 105, A: 255}, disabled: color.NRGBA{R: 113, G: 133, B: 136, A: 255},
			inputBackground: color.NRGBA{R: 28, G: 42, B: 46, A: 255}, overlayBackground: color.NRGBA{R: 22, G: 34, B: 38, A: 255}, foregroundPrimary: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			placeholder: color.NRGBA{R: 139, G: 169, B: 170, A: 255}, separator: color.NRGBA{R: 54, G: 78, B: 81, A: 255}, selection: color.NRGBA{R: 28, G: 113, B: 111, A: 255}, pressed: color.NRGBA{R: 24, G: 126, B: 121, A: 255}, menuBackground: color.NRGBA{R: 28, G: 42, B: 46, A: 255},
		}
	}
}

func (t *KaizenTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t *KaizenTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t *KaizenTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.base.Size(name)
}
