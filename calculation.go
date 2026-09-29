package main

import "unicode"

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
