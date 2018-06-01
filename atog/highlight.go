/*
 * Copyright (c) 2018, João Lucas Nunes e Silva
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL JOÃO LUCAS NUNES E SILVA BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package atog

import (
	"bytes"
	"strings"
	"unicode"
)

var (
	reservedWords = make(map[string]bool)
)

func init() {
	rw := []string{
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"byte",
		"short",
		"long",
		"char",
		"rune",
		"string",
		"void",
		"()",
		"fn",
		"fun",
		"func",
		"function",
		"if",
		"else",
		"elseif",
		"elif",
		"switch",
		"break",
		"fallthrough",
		"go",
		"for",
		"do",
		"while",
		"type",
		"typedef",
		"struct",
		"class",
		"interface",
		"abstract",
		"const",
		"var",
		"def",
		"defn",
		"define",
		"in",
		"into",
		"select",
		"where",
		"from",
		"join",
		"foreach",
		"defer",
		"defn",
		"import",
		"package",
		"ns",
		"require",
	}
	for _, w := range rw {
		reservedWords[w] = true
		reservedWords[strings.ToUpper(w)] = true
		reservedWords[strings.ToTitle(w)] = true
	}
}

type state struct {
	text []rune
	idx  int
}

func (st *state) peek() rune {
	if st.idx < len(st.text) {
		return st.text[st.idx]
	}
	return 0
}

func (st *state) pop() rune {
	r := st.peek()
	st.idx++
	return r
}

// Highlight highlight a string.
func Highlight(s string) string {
	st := &state{text: []rune(s)}
	buf := &bytes.Buffer{}

	for r := st.peek(); r != 0; r = st.peek() {
		switch {
		case isQuote(r):
			str(st, buf)
		case unicode.IsNumber(r):
			number(st, buf)
		case isDelim(r):
			delim(st, buf)
		case unicode.IsLetter(r):
			word(st, buf)
		default:
			buf.WriteRune(st.pop())
		}
	}
	return string(buf.String())
}

func str(st *state, buf *bytes.Buffer) {
	quote := st.pop()
	prev := quote

	buf.WriteString("[#217844]")
	buf.WriteRune(quote)
	for r := st.peek(); r != 0 && r != quote || r == quote && prev == '\\'; r = st.peek() {
		buf.WriteRune(st.pop())
		prev = r
	}
	buf.WriteRune(st.pop())
	buf.WriteString("[white]")
}

func number(st *state, buf *bytes.Buffer) {
	buf.WriteString("[#9b7e00]")
	buf.WriteRune(st.pop())
	for unicode.IsNumber(st.peek()) {
		buf.WriteRune(st.pop())
	}
	buf.WriteString("[white]")
}

func delim(st *state, buf *bytes.Buffer) {
	buf.WriteString("[#c83737]")
	buf.WriteRune(st.pop())
	for isDelim(st.peek()) {
		buf.WriteRune(st.pop())
	}
	buf.WriteString("[white]")
}

func word(st *state, buf *bytes.Buffer) {
	rs := make([]rune, 0, 10)

	for unicode.IsLetter(st.peek()) {
		rs = append(rs, st.pop())
	}

	w := string(rs)
	if isHilightedWord(w) {
		buf.WriteString("[#c83737]")
		buf.WriteString(w)
		buf.WriteString("[white]")
	} else if reservedWords[w] {
		buf.WriteString("[#162d50]")
		buf.WriteString(w)
		buf.WriteString("[white]")
	} else {
		buf.WriteString(w)
	}
}

func isHilightedWord(w string) bool {
	return strings.EqualFold(w, "line") ||
		strings.EqualFold(w, "column") ||
		strings.EqualFold(w, "col") ||
		w == "L" ||
		w == "C"
}

func isQuote(r rune) bool {
	return strings.ContainsRune("'\"`", r)
}

func isDelim(r rune) bool {
	return strings.ContainsRune("!@#$%&*()_-=+[]{}~^;:.>,</?", r)
}
