// bulletlist - A simple tool to generate an ordered list in text
// Copyright (C) 2016 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"fmt"
	"github.com/StefanSchroeder/Golang-Roman"
	flag "maunium.net/go/mauflag"
	"os"
	"strconv"
	"strings"
)

// Node ...
type Node struct {
	Number   int
	Length   int
	Type     string
	Closure  string
	AltPad   bool
	Children []Node
}

// Possible number types
const (
	TypeNumber     = "number"
	TypeRomanBig   = "romanbig"
	TypeRoman      = "roman"
	TypeRomanSmall = "romansmall"
	TypeAlpha      = "alpha"
	TypeAlphaBig   = "alphabig"
	TypeAlphaSmall = "alphasmall"
)

func main() {
	flag.Parse()
	master := Node{Number: -1, Length: -1}
	for _, str := range flag.Args() {
		master.Children = append(master.Children, Parse(str))
	}
	master.Print("", "   ")
}

// Parse parses a Node from a string
func Parse(str string) Node {
	parts := strings.Split(str, ";")

	primParts := strings.Split(parts[0], ":")
	number, _ := strconv.Atoi(primParts[0])
	length, _ := strconv.Atoi(primParts[1])

	node := Node{Number: number, Length: length, Type: TypeNumber, Closure: "."}

	for _, opt := range primParts[2:] {
		optParts := strings.SplitN(opt, "=", 2)
		key := optParts[0]
		val := ""
		if len(optParts) > 1 {
			val = optParts[1]
		}

		switch key {
		case "type":
			node.Type = val
		case "closure":
			node.Closure = val
		case "leftpad":
			node.AltPad = true
		}
	}

	if len(parts) > 0 {
		node.Children = ParseChildren(strings.Join(parts[1:], ";"))
	} else {
		node.Children = []Node{}
	}

	return node
}

// ParseChildren parses children from a string.
func ParseChildren(children string) []Node {
	var childNodes []Node
	openings := 0
	startIndex := -1
	previousWasEscape := false
	for i, char := range children {
		if previousWasEscape {
			previousWasEscape = false
			continue
		} else if char == '\\' {
			previousWasEscape = true
		} else if char == ';' && openings == 0 {
			startIndex = -1
		} else if char == '[' || char == '{' {
			openings++
			if startIndex == -1 {
				startIndex = i
			}
		} else if char == ']' || char == '}' {
			openings--
			if openings == 0 {
				childNodes = append(childNodes, Parse(children[startIndex+1:i]))
			} else if openings < 0 {
				fmt.Println("Syntax Error: Unopened brackets")
				os.Exit(1)
			}
		}
	}

	if openings > 0 {
		fmt.Println("Syntax Error: Unclosed brackets!")
		os.Exit(1)
	}

	return childNodes
}

// Print prints the node in a pretty way.
func (node Node) Print(indent, indentIncr string) {
	if node.Number == -1 {
		for _, actualMaster := range node.Children {
			actualMaster.Print(indent, indentIncr)
		}
		return
	}

	if len(indent) == 0 {
		fmt.Printf("%s%s ", node.FormatN(node.Number), node.Closure)
		indent += indentIncr
	}

	maxLen := node.LongestN(node.Length)
	for i := 1; i <= node.Length; i++ {
		if i != 1 {
			fmt.Print(indent)
		}

		formatted := node.FormatN(i)
		if node.AltPad {
			fmt.Print(node.NSpaces(maxLen - len(formatted)))
		}
		fmt.Print(formatted)
		fmt.Print(node.Closure, " ")
		if !node.AltPad {
			fmt.Print(node.NSpaces(maxLen - len(formatted)))
		}

		childPrinted := false
		for _, child := range node.Children {
			if child.Number == i {
				child.Print(
					indent+indentIncr+node.NSpaces(maxLen-1),
					indentIncr,
				)
				childPrinted = true
				break
			}
		}
		if !childPrinted {
			fmt.Print("\n")
		}
	}
}

// NSpaces returns a string with the given number of spaces.
func (node Node) NSpaces(n int) string {
	return strings.Repeat(" ", n)
}

// FormatN formats the number.
func (node Node) FormatN(n int) string {
	switch node.Type {
	case TypeRoman:
		fallthrough
	case TypeRomanBig:
		return strings.ToUpper(roman.Roman(n))
	case TypeRomanSmall:
		return strings.ToLower(roman.Roman(n))
	case TypeAlpha:
		fallthrough
	case TypeAlphaSmall:
		return string('a' - 1 + n)
	case TypeAlphaBig:
		return string('A' - 1 + n)
	case TypeNumber:
		fallthrough
	default:
		return strconv.Itoa(n)
	}
}

// LongestN returns the longest representation
func (node Node) LongestN(n int) int {
	switch node.Type {
	case TypeRomanBig:
		fallthrough
	case TypeRomanSmall:
		fallthrough
	case TypeRoman:
		longest := 0
		for i := 1; i <= n; i++ {
			len := len(roman.Roman(i))
			if len > longest {
				longest = len
			}
		}
		return longest
	case TypeAlphaSmall:
		fallthrough
	case TypeAlphaBig:
		fallthrough
	case TypeAlpha:
		return 1
	case TypeNumber:
		fallthrough
	default:
		return len(strconv.Itoa(n))
	}
}
