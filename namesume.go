// Copyright (c) 2026 Naawaakrit
// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	"fmt"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ---------- ตารางค่าตัวอักษร (เลขศาสตร์) ----------

var english = map[rune]int{
	'a': 1, 'i': 1, 'j': 1, 'q': 1, 'y': 1,
	'b': 2, 'k': 2, 'r': 2,
	'c': 3, 'g': 3, 'l': 3, 's': 3,
	'd': 4, 'm': 4, 't': 4,
	'e': 5, 'h': 5, 'n': 5, 'x': 5,
	'u': 6, 'v': 6, 'w': 6,
	'o': 7, 'z': 7,
	'f': 8, 'p': 8,
}

var thaiConsonant = map[rune]int{
	'ก': 1, 'ด': 1, 'ท': 1, 'ถ': 1, 'ภ': 1, 'ฤ': 1,
	'ข': 2, 'ช': 2, 'บ': 2, 'ป': 2, 'ง': 2,
	'ฆ': 3, 'ฑ': 3, 'ฒ': 3, 'ต': 3, 'ฃ': 3,
	'ค': 4, 'ธ': 4, 'ร': 4, 'ญ': 4, 'ษ': 4,
	'ฉ': 5, 'ณ': 5, 'ฌ': 5, 'น': 5, 'ม': 5, 'ห': 5, 'ฮ': 5, 'ฎ': 5, 'ฬ': 5,
	'จ': 6, 'ล': 6, 'ว': 6, 'อ': 6,
	'ศ': 7, 'ส': 7, 'ซ': 7,
	'ย': 8, 'พ': 8, 'ฟ': 8, 'ผ': 8, 'ฝ': 8,
	'ฏ': 9, 'ฐ': 9,
}

var thaiVowel = map[rune]int{
	'า': 1, 'ุ': 1, 'ำ': 1,
	'เ': 2, 'แ': 2, 'ู': 2,
	'โ': 4, 'ะ': 4, 'ิ': 4,
	'ึ': 5,
	'ใ': 6,
	'ี': 7, 'ื': 7,
	'็': 8,
	'ไ': 9, '์': 9,
}

var thaiTone = map[rune]int{
	'่': 1,
	'้': 2,
	'๋': 3,
	'ั': 4,
	'๊': 7,
}

func classify(r rune) (value int, category string, ok bool) {
	lower := unicode.ToLower(r)
	if v, found := english[lower]; found {
		return v, "อังกฤษ", true
	}
	if v, found := thaiConsonant[r]; found {
		return v, "พยัญชนะไทย", true
	}
	if v, found := thaiVowel[r]; found {
		return v, "สระ", true
	}
	if v, found := thaiTone[r]; found {
		return v, "วรรณยุกต์", true
	}
	return 0, "", false
}

type LetterScore struct {
	Char     string
	Category string
	Value    int
}

type Result struct {
	Name         string
	Scores       []LetterScore
	SumEnglish   int
	SumConsonant int
	SumVowel     int
	SumTone      int
	Total        int
}

func calcName(name string) Result {
	res := Result{Name: name}
	for _, r := range name {
		v, cat, ok := classify(r)
		if !ok {
			continue
		}
		res.Scores = append(res.Scores, LetterScore{string(r), cat, v})
		res.Total += v
		switch cat {
		case "อังกฤษ":
			res.SumEnglish += v
		case "พยัญชนะไทย":
			res.SumConsonant += v
		case "สระ":
			res.SumVowel += v
		case "วรรณยุกต์":
			res.SumTone += v
		}
	}
	return res
}

// ---------- ความหมายของเลขผลรวม ----------

// meanings เก็บความหมายของเลขแต่ละตัว แก้ไข/เพิ่มเลขได้ตรงนี้เลย
// รูปแบบ:  เลข: "ความหมาย",
var meanings = map[int]string{
	// 15: "ใส่ความหมายของเลข 15 ที่นี่",
	// 20: "ใส่ความหมายของเลข 20 ที่นี่",
	// 35: "ใส่ความหมายของเลข 35 ที่นี่",
}

func meaningFor(n int) string {
	if m, ok := meanings[n]; ok && m != "" {
		return m
	}
	return fmt.Sprintf("(ยังไม่มีความหมายของเลข %d ในตาราง meanings)", n)
}

// ---------- Fyne GUI ----------

// renderSection เพิ่มหัวข้อ + ตารางตัวอักษรของชื่อ/นามสกุลลงใน box
func renderSection(box *fyne.Container, label string, res Result) {
	box.Add(widget.NewLabelWithStyle(
		fmt.Sprintf("%s: %s", label, res.Name),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	))
	if len(res.Scores) == 0 {
		box.Add(widget.NewLabel("(ไม่มีตัวอักษรที่นับได้)"))
		return
	}
	for _, s := range res.Scores {
		box.Add(widget.NewLabel(fmt.Sprintf("%-6s %-14s %d", s.Char, s.Category, s.Value)))
	}
}

// summaryLine สรุปผลรวมของชื่อ/นามสกุลหนึ่งส่วน พร้อมความหมายของเลข
func summaryLine(label string, res Result) string {
	if res.Total == 0 {
		return fmt.Sprintf("%s — (ไม่มีตัวอักษรที่นับได้)", label)
	}
	return fmt.Sprintf(
		"%s — รวม %d (อังกฤษ %d, พยัญชนะไทย %d, สระ %d, วรรณยุกต์ %d)\nความหมายเลข %d ของ%s: %s",
		label, res.Total, res.SumEnglish, res.SumConsonant, res.SumVowel, res.SumTone,
		res.Total, label, meaningFor(res.Total),
	)
}

// combinedLine สรุปผลรวมของชื่อ + นามสกุล พร้อมความหมายของเลขรวม
func combinedLine(a, b int) string {
	t := a + b
	if t == 0 {
		return ""
	}
	return fmt.Sprintf(
		"ผลรวมทั้งหมด (ชื่อ + นามสกุล): %d + %d = %d\nความหมายเลข %d ของชื่อ + นามสกุล: %s",
		a, b, t, t, meaningFor(t),
	)
}

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
