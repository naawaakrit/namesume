// Copyright (c) 2026 Naawaakrit
// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("ผลรวมตัวอักษรชื่อ")
	w.Resize(fyne.NewSize(460, 680))

	title := widget.NewLabelWithStyle(
		"คำนวณผลรวมตัวอักษรชื่อ-นามสกุล (ไทย/อังกฤษ)",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	firstEntry := widget.NewEntry()
	firstEntry.SetPlaceHolder("เช่น สมชาย หรือ John")
	lastEntry := widget.NewEntry()
	lastEntry.SetPlaceHolder("เช่น ใจดี หรือ Smith")

	form := widget.NewForm(
		widget.NewFormItem("ชื่อ", firstEntry),
		widget.NewFormItem("นามสกุล", lastEntry),
	)

	resultsBox := container.NewVBox()
	scroll := container.NewVScroll(resultsBox)
	scroll.SetMinSize(fyne.NewSize(420, 320))

	summary := widget.NewLabel("")
	summary.Wrapping = fyne.TextWrapWord
	summaryScroll := container.NewVScroll(summary)
	summaryScroll.SetMinSize(fyne.NewSize(420, 180))

	runCalc := func() {
		first := strings.TrimSpace(firstEntry.Text)
		last := strings.TrimSpace(lastEntry.Text)

		resultsBox.Objects = nil // ล้างผลลัพธ์เก่า
		if first == "" && last == "" {
			summary.SetText("")
			resultsBox.Refresh()
			return
		}

		resFirst := calcName(first)
		resLast := calcName(last)

		renderSection(resultsBox, "ชื่อ", resFirst)
		resultsBox.Add(widget.NewSeparator())
		renderSection(resultsBox, "นามสกุล", resLast)
		resultsBox.Refresh()

		summary.SetText(strings.Join([]string{
			summaryLine("ชื่อ", resFirst),
			summaryLine("นามสกุล", resLast),
			combinedLine(resFirst.Total, resLast.Total),
		}, "\n"))
	}

	firstEntry.OnSubmitted = func(string) { runCalc() }
	lastEntry.OnSubmitted = func(string) { runCalc() }
	calcBtn := widget.NewButton("คำนวณ", runCalc)

	content := container.NewVBox(
		title,
		form,
		calcBtn,
		widget.NewSeparator(),
		scroll,
		widget.NewSeparator(),
		summaryScroll,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
