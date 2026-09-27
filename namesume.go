// Copyright (c) 2026 Naawaakrit
// Copyright (c) 2026 Nawakarit
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License v3.0.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// ---------- ตรรกะคำนวณ (เหมือนเวอร์ชัน CLI) ----------

func englishValue(r rune) (int, bool) {
	switch {
	case r >= 'a' && r <= 'z':
		return int(r-'a') + 1, true
	case r >= 'A' && r <= 'Z':
		return int(r-'A') + 1, true
	}
	return 0, false
}

var thaiConsonants = []rune("กขฃคฅฆงจฉชซฌญฎฏฐฑฒณดตถทธนบปผฝพฟภมยรฤลฦวศษสหฬอฮ")

func thaiValue(r rune) (int, bool) {
	for i, c := range thaiConsonants {
		if c == r {
			return i + 1, true
		}
	}
	return 0, false
}

type LetterScore struct {
	Char  string
	Value int
	Lang  string
}

type Result struct {
	Name   string
	Scores []LetterScore
	SumEN  int
	SumTH  int
	Total  int
}

func calcName(name string) Result {
	res := Result{Name: name}
	for _, r := range name {
		if v, ok := englishValue(r); ok {
			res.Scores = append(res.Scores, LetterScore{string(r), v, "EN"})
			res.SumEN += v
			res.Total += v
		} else if v, ok := thaiValue(r); ok {
			res.Scores = append(res.Scores, LetterScore{string(r), v, "TH"})
			res.SumTH += v
			res.Total += v
		}
	}
	return res
}

// ---------- หน้าเว็บ (GUI) ----------

const pageTpl = `
<!DOCTYPE html>
<html lang="th">
<head>
<meta charset="UTF-8">
<title>คำนวณผลรวมตัวอักษรชื่อ</title>
<style>
  body { font-family: "TH Sarabun New", "Noto Sans Thai", sans-serif; background:#f4f6f8; margin:0; padding:2rem; }
  .card { max-width:560px; margin:0 auto; background:#fff; border-radius:12px; padding:2rem; box-shadow:0 2px 10px rgba(0,0,0,.08); }
  h1 { font-size:1.4rem; margin-top:0; }
  input[type=text] { width:100%; padding:.6rem; font-size:1.1rem; border:1px solid #ccc; border-radius:8px; box-sizing:border-box; }
  button { margin-top:1rem; padding:.6rem 1.4rem; font-size:1rem; border:none; border-radius:8px; background:#2563eb; color:#fff; cursor:pointer; }
  button:hover { background:#1d4ed8; }
  table { width:100%; border-collapse:collapse; margin-top:1.5rem; }
  th, td { padding:.4rem .6rem; text-align:center; border-bottom:1px solid #eee; }
  th { color:#666; font-weight:600; font-size:.9rem; }
  .totals { margin-top:1rem; font-size:1.05rem; }
  .totals span { display:block; margin:.2rem 0; }
  .grand { font-weight:700; color:#2563eb; }
</style>
</head>
<body>
  <div class="card">
    <h1>ผลรวมตัวอักษรชื่อ (ไทย/อังกฤษ)</h1>
    <form method="POST" action="/">
      <input type="text" name="name" placeholder="พิมพ์ชื่อ เช่น สมชาย John" value="{{.Name}}" autofocus>
      <button type="submit">คำนวณ</button>
    </form>

    {{if .Scores}}
    <table>
      <tr><th>ตัวอักษร</th><th>ภาษา</th><th>ค่า</th></tr>
      {{range .Scores}}
      <tr><td>{{.Char}}</td><td>{{.Lang}}</td><td>{{.Value}}</td></tr>
      {{end}}
    </table>
    <div class="totals">
      <span>ผลรวมอังกฤษ: {{.SumEN}}</span>
      <span>ผลรวมไทย: {{.SumTH}}</span>
      <span class="grand">ผลรวมทั้งหมด: {{.Total}}</span>
    </div>
    {{end}}
  </div>
</body>
</html>
`

var tpl = template.Must(template.New("page").Parse(pageTpl))

func handler(w http.ResponseWriter, r *http.Request) {
	var res Result
	if r.Method == http.MethodPost {
		name := strings.TrimSpace(r.FormValue("name"))
		if name != "" {
			res = calcName(name)
		}
	}
	if err := tpl.Execute(w, res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	addr := "localhost:8080"
	fmt.Printf("เปิดเบราว์เซอร์ไปที่ http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
