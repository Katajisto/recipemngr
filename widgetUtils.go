package main

import g "github.com/AllenDang/giu"

func LabelWithSize(text string, size float32) g.Widget {
	return g.Style().SetFontSize(size).To(g.Label(text))
}

func ButtonWithSize(text string, size float32, onClick func()) g.Widget {
	return g.Style().SetFontSize(size).To(g.Button(text).OnClick(onClick))
}

// CTF is a Condition, True, False util function. If the condition is true
// it returns ifTrue. If the condition is false, it returns ifFalse.
func CTF[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}
