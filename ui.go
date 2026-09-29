package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

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
