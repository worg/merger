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
		pr   int
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
	orders = map[string]interface{}{
		`coffee`:    1,
		`milkBrand`: `Happy Cow`,
	}
	dSimple     = simpleT{20, `Springfield`}
	dComplex    = complexT{dSimple, `Homer`, 1}
	eComplex    = complexT{Name: `Homer`}
	dPointer    = pointerT{&simpleT{5, `Springfield`}}
	ePointer    = pointerT{nil}
	dComplexPtr = complexPtr{&complexT{
		simpleT{1, `Springfield`},
		`Homer`,
		1,
	},
		666,
	}
	eComplexPtr = complexPtr{nil, 999}
	dSlice      = sliceT{[]int{1, 2, 3}}
	eSlice      = sliceT{}
	dMap        = mapT{orders}
	eMap        = mapT{nil}
	tnow        = time.Now()
	tzero       time.Time
	dTime       = timeT{tnow}
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
	a := simpleT{ID: 666}

	if err := Merge(a, dSimple); err != ErrNoPtr {
		t.Error(`Failed to check for struct pointer`)
	}
}

func TestEqualType(t *testing.T) {
	a := simpleT{ID: 000}

	if err := Merge(&a, dPointer); err != ErrDistinctType {
		t.Error(`Failed to check for type equality`)
	}
}

func TestSimple(t *testing.T) {
	a := simpleT{ID: 2}

	if err := Merge(&a, dSimple); err != nil {
		t.Error(err)
	}

	if a.Address != dSimple.Address {
		t.Error(`Failed to merge fields`)
	}

	if a.ID != 2 {
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

	if a.ID != 0 {
		t.Error(`Failed to merge fields`)
	}
}

func TestComplex(t *testing.T) {
	a := complexT{simpleT{ID: 100}, `John`, 0}
	if err := Merge(&a, dComplex); err != nil {
		t.Error(err)
	}

	if a.Name != `John` {
		t.Error(`Falied to merge fields`)
	}

	if a.Address != `Springfield` {
		t.Error(`Failed to merge fields, expected Springfield got `, a.Address)
	}

	if a.ID != 100 {
		t.Error(`Failed to merge fields`)
	}

	if a.pr != 0 {
		t.Error(0, ` expected, got`, a.pr)
	}
}

func TestComplexEmpty(t *testing.T) {
	a := complexT{simpleT{ID: 200}, `Jane`, 2}
	if err := Merge(&a, eComplex); err != nil {
		t.Error(err)
	}

	if a.Name != `Jane` {
		t.Error(`Falied to merge fields`)
	}

	if a.Address != `` {
		t.Error(`Failed to merge fields`)
	}

	if a.ID != 200 {
		t.Error(`Failed to merge fields`)
	}
}

func TestComplexPtr(t *testing.T) {
	a := complexPtr{
		&complexT{simpleT{ID: 600}, `Doe`, 2},
		0,
	}

	if err := Merge(&a, dComplexPtr); err != nil {
		t.Error(err)
	}

	if a.complexT.ID != 600 {
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

	if a.complexT.ID != 1 {
		t.Error(`Failed to merge fields`)
	}

	if a.Name != `Homer` {
		t.Error(`Failed to merge fields`)
	}

	if a.Location != 90 {
		t.Error(`Failed to merge fields`)
	}
}

func TestTime(t *testing.T) {
	ntime := time.Now().Add(time.Hour)
	nt := timeT{}
	if err := Merge(&nt, dTime); err != nil {
		t.Error(`Failed to merge time field `, err)
	}

	if nt.Now != tnow {
		t.Error(tnow, `expected, got `, nt.Now)
	}

	// t.Log(tnow, `expected, got `, nt.Now)

	nt.Now = tzero
	if err := Merge(&nt, timeT{ntime}); err != nil {
		t.Error(`Failed to merge time field `, err)
	}

	if nt.Now != ntime {
		t.Error(ntime, ` expected, got `, nt.Now)
	}

	// t.Log(ntime, `expected, got `, nt.Now)

	if err := Merge(&nt, dTime); err != nil {
		t.Error(`Failed to merge time field `, err)
	}

	if nt.Now != ntime {
		t.Error(tnow, `expected, got `, nt.Now)
	}

	// t.Log(ntime, `expected, got `, nt.Now)
}

func TestMap(t *testing.T) {
	nmap := mapT{map[string]interface{}{
		`coffee`:  3,
		`veggies`: `lots`,
	},
	}

	if err := Merge(&nmap, dMap); err != nil {
		t.Error(`Failed to merge map `, err)
	}

	if nmap.Orders[`coffee`] != 3 {
		t.Errorf(`Failed to merge map expected %d, got %d`, 3, nmap.Orders[`coffee`])
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
