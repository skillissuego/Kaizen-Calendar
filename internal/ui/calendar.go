package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"kaizen-calendar/internal/storage"
)

type calendarCell struct {
	widget.BaseWidget
	date     time.Time
	progress storage.DayProgress
	selected bool
	onTapped func(time.Time)
}

func newCalendarCell(date time.Time, progress storage.DayProgress, selected bool, onTapped func(time.Time)) *calendarCell {
	cell := &calendarCell{
		date:     date,
		progress: progress,
		selected: selected,
		onTapped: onTapped,
	}
	cell.ExtendBaseWidget(cell)
	return cell
}

func (c *calendarCell) Tapped(*fyne.PointEvent) {
	if c.onTapped != nil {
		c.onTapped(c.date)
	}
}

func (c *calendarCell) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(calendarFill(c.progress.Percent, c.progress.Total == 0))
	text := canvas.NewText(c.date.Format("2"), calendarTextColor())
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle = fyne.TextStyle{Bold: c.selected}

	objects := []fyne.CanvasObject{background, text}
	return &calendarCellRenderer{cell: c, background: background, text: text, objects: objects}
}

type calendarCellRenderer struct {
	cell       *calendarCell
	background *canvas.Rectangle
	text       *canvas.Text
	objects    []fyne.CanvasObject
}

func (r *calendarCellRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.text.Resize(fyne.NewSize(size.Width, theme.TextSize()))
	r.text.Move(fyne.NewPos(0, (size.Height-theme.TextSize())/2))
}

func (r *calendarCellRenderer) MinSize() fyne.Size {
	return fyne.NewSize(48, 48)
}

func (r *calendarCellRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *calendarCellRenderer) Refresh() {
	r.background.FillColor = calendarFill(r.cell.progress.Percent, r.cell.progress.Total == 0)
	r.text.Color = calendarTextColor()
	r.text.TextStyle = fyne.TextStyle{Bold: r.cell.selected}
	r.background.Refresh()
	r.text.Refresh()
}

func (r *calendarCellRenderer) Destroy() {}

// calendarFill blends the theme background into the theme accent color
// so every Kaizen theme gets a readable progress scale.
func calendarFill(percent int, empty bool) color.Color {
	settings := fyne.CurrentApp().Settings()
	current := settings.Theme()
	variant := settings.ThemeVariant()
	background := current.Color(theme.ColorNameBackground, variant)
	if empty || percent <= 0 {
		return background
	}
	if percent >= 100 {
		return current.Color(theme.ColorNamePrimary, variant)
	}
	return mixColors(background, current.Color(theme.ColorNamePrimary, variant), percent)
}

func calendarTextColor() color.Color {
	settings := fyne.CurrentApp().Settings()
	return settings.Theme().Color(theme.ColorNameForeground, settings.ThemeVariant())
}

func mixColors(from color.Color, to color.Color, percent int) color.Color {
	fromR, fromG, fromB, fromA := from.RGBA()
	toR, toG, toB, toA := to.RGBA()
	weight := int64(percent)
	return color.NRGBA64{
		R: uint16((int64(fromR) + (int64(toR)-int64(fromR))*weight/100)),
		G: uint16((int64(fromG) + (int64(toG)-int64(fromG))*weight/100)),
		B: uint16((int64(fromB) + (int64(toB)-int64(fromB))*weight/100)),
		A: uint16((int64(fromA) + (int64(toA)-int64(fromA))*weight/100)),
	}
}
