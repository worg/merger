// Copyright (c) 2014 Hiram Jerónimo Pérez worg{at}linuxmail[dot]org
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package merger is an utility to merge structs of the same type
package merger

import (
	"testing"
	"time"
)

type (
	simpleT struct {
		ID      int
		Address string
	}

	complexT struct {
		simpleT
		Name string
	}

	pointerT struct {
		*simpleT
	}

	complexPtr struct {
		*complexT
		Location int
	}

	sliceT struct {
		Answers []int
	}

	mapT struct {
		Orders map[string]interface{}
	}

	timeT struct {
		Now time.Time
	}
)

var (
	orders      = make(map[string]interface{})
	dSimple     = simpleT{20, `Springfield`}
	dComplex    = complexT{simpleT{10, `Springfield`}, `Homer`}
	eComplex    = complexT{Name: `Homer`}
	dPointer    = pointerT{&simpleT{5, `Springfield`}}
	ePointer    = pointerT{nil}
	dComplexPtr = complexPtr{&complexT{
		simpleT{1, `Springfield`},
		`Homer`,
	},
		666,
	}
	eComplexPtr = complexPtr{nil, 999}
	dSlice      = sliceT{[]int{1, 2, 3}}
	eSlice      = sliceT{}
	dMap        = mapT{orders}
	eMap        = mapT{nil}
	dTime       = timeT{time.Now()}
	eTime       = timeT{}
)

func TestNil(t *testing.T) {
	if err := Merge(nil, nil); err != ErrNilArguments {
		t.Error(`Failed to check for nil values`)
	}

	if err := Merge(dPointer, nil); err != ErrNilArguments {
		t.Error(`Failed to check for nil values`)
	}

	if err := Merge(nil, dSimple); err != ErrNilArguments {
		t.Error(`Failed to check for nil values`)
	}
}

func TestStructPtr(t *testing.T) {
	a := simpleT{Id: 666}

	if err := Merge(a, dSimple); err != ErrNoPtr {
		t.Error(`Failed to check for struct pointer`)
	}
}

func TestEqualType(t *testing.T) {
	a := simpleT{Id: 000}

	if err := Merge(&a, dPointer); err != ErrDistinctType {
		t.Error(`Failed to check for type equality`)
	}
}

func TestSimple(t *testing.T) {
	a := simpleT{Id: 2}

	if err := Merge(&a, dSimple); err != nil {
		t.Error(err)
	}

	if a.Address != dSimple.Address {
		t.Error(`Failed to merge fields`)
	}

	if a.Id != 2 {
		t.Error(`Failed to merge fields`)
	}
}

func TestSimpleEmpty(t *testing.T) {
	a := simpleT{Address: `Omaha`}

	if err := Merge(&a, simpleT{}); err != nil {
		t.Error(err)
	}

	if a.Address != `Omaha` {
		t.Error(`Failed to merge fields`)
	}

	if a.Id != 0 {
		t.Error(`Failed to merge fields`)
	}
}

func TestComplex(t *testing.T) {
	a := complexT{simpleT{Id: 100}, `John`}
	if err := Merge(&a, dComplex); err != nil {
		t.Error(err)
	}

	if a.Name != `John` {
		t.Error(`Falied to merge fields`)
	}

	if a.Address != `Springfield` {
		t.Error(`Failed to merge fields`)
	}

	if a.Id != 100 {
		t.Error(`Failed to merge fields`)
	}
}

func TestComplexEmpty(t *testing.T) {
	a := complexT{simpleT{Id: 200}, `Jane`}
	if err := Merge(&a, eComplex); err != nil {
		t.Error(err)
	}

	if a.Name != `Jane` {
		t.Error(`Falied to merge fields`)
	}

	if a.Address != `` {
		t.Error(`Failed to merge fields`)
	}

	if a.Id != 200 {
		t.Error(`Failed to merge fields`)
	}
}

func TestComplexPtr(t *testing.T) {
	a := complexPtr{
		&complexT{simpleT{Id: 600}, `Doe`},
		0,
	}

	if err := Merge(&a, dComplexPtr); err != nil {
		t.Error(err)
	}

	if a.complexT.Id != 600 {
		t.Error(`Failed to merge fields`)
	}

	if a.Name != `Doe` {
		t.Error(`Failed to merge fields`)
	}

	if a.Location != 666 {
		t.Error(`Failed to merge fields`)
	}

	a = complexPtr{
		nil,
		90,
	}

	if err := Merge(&a, dComplexPtr); err != nil {
		t.Error(err)
	}

	if a.complexT.Id != 1 {
		t.Error(`Failed to merge fields`)
	}

	if a.Name != `Homer` {
		t.Error(`Failed to merge fields`)
	}

	if a.Location != 90 {
		t.Error(`Failed to merge fields`)
	}
}

/*
func TestComplexPtrEmpty() {
	a := complexPtr{
		nil,
		0,
	}

}
*/
