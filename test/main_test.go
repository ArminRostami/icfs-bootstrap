// Package test _
package test

import (
	"mime"
	"testing"
)

func Test1(t *testing.T) {
	ext := ".png"
	tp := mime.TypeByExtension(ext)
	// tp = strings.Split(tp, "/")[0]
	t.Log(tp)
}
