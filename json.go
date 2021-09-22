package main

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
)

type canvObj struct {
	Type   string
	Struct fyne.CanvasObject `json:",omitempty"`
}

type cont struct {
	canvObj
	Layout string `json:",omitempty"`
	Objects []interface{}
}

func encodeObj(obj fyne.CanvasObject) interface{} {
	if c, ok := obj.(*fyne.Container); ok { // the content of an overlayWidget container
		return encodeObj(c.Objects[0]) // 0 is the widget, 1 is the overlayWidget
	} else if c, ok := obj.(*overlayContainer); ok {
		var node cont
		node.Type = "*fyne.Container"
		node.Layout = strings.Split(reflect.TypeOf(c.c.Layout).String(), ".")[1]
		node.Layout = strings.ToTitle(node.Layout[0:1]) + node.Layout[1:]
		p := strings.Index(node.Layout, "Layout")
		if p > 0 {
			node.Layout = node.Layout[:p]
		}
		if node.Layout == "Box" {
			node.Layout = "VBox" // TODO remove this hack with layoutProps
		}
		for _, o := range c.c.Objects {
			node.Objects = append(node.Objects, encodeObj(o))
		}
		return &node
	}

	return encodeWidget(obj)
}

func encodeWidget(obj fyne.CanvasObject) interface{} {
	return &canvObj{Type: reflect.TypeOf(obj).String(), Struct: obj}
}

func DecodeJSON(r io.Reader) fyne.CanvasObject {
	var data interface{}
	_ = json.NewDecoder(r).Decode(&data)

	return decodeMap(data.(map[string]interface{}))
}

func decodeTextStyle(m map[string]interface{}) (s fyne.TextStyle) {
	if m["Bold"] == true {
		s.Bold = true
	}
	if m["Italic"] == true {
		s.Italic = true
	}
	if m["Monospace"] == true {
		s.Monospace = true
	}

	if m["TabWidth"] != 0 {
		s.TabWidth = int(m["TabWidth"].(float64))
	}
	return
}

func decodeMap(m map[string]interface{}) fyne.CanvasObject {
	if m["Type"] == "*fyne.Container" {
		obj := &fyne.Container{}
		obj.Layout = layouts[m["Layout"].(string)].create(nil)
		for _, o := range m["Objects"].([]interface{}) {
			obj.Objects = append(obj.Objects, decodeMap(o.(map[string]interface{})))
		}
		return obj
	}

	obj := widgets[m["Type"].(string)].create()
	e := reflect.ValueOf(obj).Elem()
	for k, v := range m["Struct"].(map[string]interface{}) {
		f := e.FieldByName(k)

		if f.Type().String() == "fyne.TextAlign" || f.Type().String() == "fyne.TextWrap" ||
			f.Type().String() == "widget.ButtonAlign" || f.Type().String() == "widget.ButtonImportance" || f.Type().String() == "widget.ButtonIconPlacement" {
			f.SetInt(int64(reflect.ValueOf(v).Float()))
		} else if f.Type().String() == "fyne.TextStyle" {
			f.Set(reflect.ValueOf(decodeTextStyle(reflect.ValueOf(v).Interface().(map[string]interface{}))))
		} else if f.Type().String() == "fyne.Resource" {
			res := icons[reflect.ValueOf(v).String()]
			if res != nil {
				f.Set(reflect.ValueOf(res))
			}
		} else {
			f.Set(reflect.ValueOf(v))
		}
	}

	return obj
}

func EncodeJSON(obj fyne.CanvasObject, w io.Writer) {
	tree := encodeObj(obj)

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	_ = e.Encode(tree)
}
