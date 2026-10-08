package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	kaizenapp "kaizen-calendar/internal/app"
	"kaizen-calendar/internal/storage"
)

type TodoScreen struct {
	service  *kaizenapp.TodoService
	window   fyne.Window
	input    *widget.Entry
	items    *fyne.Container
	root     *fyne.Container
	updating bool
}

func NewTodoScreen(service *kaizenapp.TodoService, window fyne.Window) fyne.CanvasObject {
	screen := &TodoScreen{
		service: service,
		window:  window,
		input:   widget.NewEntry(),
		items:   container.NewVBox(),
	}
	screen.input.SetPlaceHolder("Новое действие")
	add := widget.NewButton("Добавить", func() {
		if err := screen.service.Create(context.Background(), screen.input.Text); err != nil {
			dialog.ShowError(err, screen.window)
			return
		}
		screen.input.SetText("")
		screen.refreshLater()
	})
	addRow := container.NewBorder(nil, nil, nil, add, screen.input)
	screen.root = container.NewBorder(addRow, nil, nil, nil, container.NewVScroll(screen.items))
	screen.refresh()
	return screen.root
}

func (s *TodoScreen) refresh() {
	s.items.Objects = nil
	items, err := s.service.List(context.Background())
	if err != nil {
		s.items.Add(widget.NewLabel("Ошибка загрузки задач: " + err.Error()))
		s.items.Refresh()
		return
	}
	if len(items) == 0 {
		s.items.Add(widget.NewLabel("Список задач пуст."))
		s.items.Refresh()
		return
	}
	for _, item := range items {
		s.items.Add(s.todoRow(item))
	}
	s.items.Refresh()
}

func (s *TodoScreen) todoRow(item storage.TodoItem) fyne.CanvasObject {
	check := widget.NewCheck(item.Title, func(completed bool) {
		if s.updating {
			return
		}
		if err := s.service.SetCompleted(context.Background(), item.ID, completed); err != nil {
			dialog.ShowError(err, s.window)
			return
		}
		// Do not remove the tapped checkbox from its container inside its callback.
		s.refreshLater()
	})
	s.updating = true
	check.SetChecked(item.Completed)
	s.updating = false
	remove := widget.NewButton("Удалить", func() {
		dialog.ShowConfirm("Удалить задачу?", item.Title, func(ok bool) {
			if !ok {
				return
			}
			if err := s.service.Delete(context.Background(), item.ID); err != nil {
				dialog.ShowError(err, s.window)
				return
			}
			s.refreshLater()
		}, s.window)
	})
	return container.NewBorder(nil, nil, check, remove, nil)
}

func (s *TodoScreen) refreshLater() {
	fyne.Do(func() {
		s.refresh()
	})
}
