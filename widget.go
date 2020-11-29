package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type widgetInfo struct {
	name   string
	create func() fyne.CanvasObject
	edit   func(fyne.CanvasObject, map[string]string) []*widget.FormItem
}

var widgets = map[string]widgetInfo{
	"*widget.Button": {
		name: "Button",
		create: func() fyne.CanvasObject {
			return widget.NewButton("Button", func() {})
		},
		edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
			b := obj.(*widget.Button)
			entry := widget.NewEntry()
			entry.SetText(b.Text)
			entry.OnChanged = func(text string) {
				b.SetText(text)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry),
				widget.NewFormItem("Icon", widget.NewSelect(iconNames, func(selected string) {
					b.SetIcon(icons[selected])
				}))}
		},
	},
	"*widget.Card": {
		name: "Card",
		create: func() fyne.CanvasObject {
			return widget.NewCard("Title", "Subtitle", widget.NewLabel("Content here"))
		},
		edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
			c := obj.(*widget.Card)
			title := widget.NewEntry()
			title.SetText(c.Title)
			title.OnChanged = func(text string) {
				c.SetTitle(text)
			}
			subtitle := widget.NewEntry()
			subtitle.SetText(c.Subtitle)
			subtitle.OnChanged = func(text string) {
				c.SetSubTitle(text)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Title", title),
				widget.NewFormItem("Title", subtitle)}
		},
	},
	"*widget.Entry": {
		name: "Entry",
		create: func() fyne.CanvasObject {
			e := widget.NewEntry()
			e.SetPlaceHolder("Entry")
			return e
		},
		edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
			l := obj.(*widget.Entry)
			entry1 := widget.NewEntry()
			entry1.SetText(l.Text)
			entry1.OnChanged = func(text string) {
				l.SetText(text)
			}
			entry2 := widget.NewEntry()
			entry2.SetText(l.PlaceHolder)
			entry2.OnChanged = func(text string) {
				l.SetPlaceHolder(text)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry1),
				widget.NewFormItem("PlaceHolder", entry2)}
		},
	},
	"*widget.Icon": {
		name: "Icon",
		create: func() fyne.CanvasObject {
			return widget.NewIcon(theme.HelpIcon())
		},
		edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
			i := obj.(*widget.Icon)
			return []*widget.FormItem{
				widget.NewFormItem("Icon", widget.NewSelect(iconNames, func(selected string) {
					i.SetResource(icons[selected])
				}))}
		},
	},
	"*widget.Label": {
		name: "Label",
		create: func() fyne.CanvasObject {
			return widget.NewLabel("Label")
		},
		edit: func(obj fyne.CanvasObject, _ map[string]string) []*widget.FormItem {
			l := obj.(*widget.Label)
			entry := widget.NewEntry()
			entry.SetText(l.Text)
			entry.OnChanged = func(text string) {
				l.SetText(text)
			}
			return []*widget.FormItem{
				widget.NewFormItem("Text", entry)}
		},
	},
	"*fyne.Container": {
		name: "Container",
		create: func() fyne.CanvasObject {
			return container.NewMax()
		},
		edit: func(obj fyne.CanvasObject, prop map[string]string) []*widget.FormItem {
			c := obj.(*fyne.Container)
			var items []*widget.FormItem
			var choose *widget.FormItem
			// TODO figure out how to work Border...
			choose = widget.NewFormItem("Layout", widget.NewSelect(layoutNames, func(l string) {
				lay := layouts[l]
				prop["layout"] = l
				c.Layout = lay.create()
				c.Refresh()
				choose.Widget.Hide()

				edit := lay.edit
				items = []*widget.FormItem{choose}
				if edit != nil {
					items = append(items, edit()...)
				}

				editForm = widget.NewForm(items...)
				paletteList.Objects = []fyne.CanvasObject{editForm}
				choose.Widget.Show()
				paletteList.Refresh()
			}))
			choose.Widget.(*widget.Select).SetSelected(prop["layout"])
			return items
		},
	},
}

// widgetNames is an array with the list of names of all the widgets
var widgetNames = extractWidgetNames()

// extractWidgetNames returns all the list of names of all the widgets from our data
func extractWidgetNames() []string {
	var widgetNamesFromData = make([]string, len(widgets))
	i := 0
	for k := range widgets {
		widgetNamesFromData[i] = k
		i++
	}
	return widgetNamesFromData
}
