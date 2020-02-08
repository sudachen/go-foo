package tests

import (
	"fmt"
	"github.com/sudachen/go-fp/lazy"
	"gotest.tools/assert"
	"gotest.tools/assert/cmp"
	"testing"
)

type Color struct {
	Color string
	Index int
}

var colors = []Color{
	{"White", 0},
	{"Yellow", 1},
	{"Blue", 2},
	{"Red", 3},
	{"Green", 4},
	{"Black", 5},
	{"Brown", 6},
	{"Azure", 7},
	{"Ivory", 8},
	{"Teal", 9},
	{"Silver", 10},
	{"Purple", 11},
	{"Navy blue", 12},
	{"Pea green", 13},
	{"Gray", 14},
	{"Orange", 15},
	{"Maroon", 16},
	{"Charcoal", 17},
	{"Aquamarine", 18},
	{"Coral", 19},
	{"Fuchsia", 20},
	{"Wheat", 21},
	{"Lime", 22},
	{"Crimson", 23},
	{"Khaki", 24},
	{"Hot pink", 25},
	{"Magenta", 26},
	{"Olden", 27},
	{"Plum", 28},
	{"Olive", 29},
	{"Cyan", 30},
}

func Test_NewPanic(t *testing.T) {
	assert.Assert(t, cmp.Panics(func() {
		lazy.New("")
	}))
	assert.Assert(t, cmp.Panics(func() {
		lazy.New(struct{ int }{0})
	}))
}

func Test_NewFromChan(t *testing.T) {
	c := make(chan Color)
	go func() {
		for _, x := range colors {
			c <- x
		}
		close(c)
	}()
	z := lazy.New(c)
	rs := z.Collect().([]Color)
	assert.DeepEqual(t, rs, colors)
}

func Test_Collect(t *testing.T) {
	z := lazy.New(colors)
	rs := z.Collect().([]Color)
	assert.DeepEqual(t, rs, colors)
}

func Test_ConqCollect(t *testing.T) {
	z := lazy.New(colors)
	rs := z.ConqCollect(8).([]Color)
	assert.DeepEqual(t, rs, colors)
	rs = z.ConqCollect(4).([]Color)
	assert.DeepEqual(t, rs, colors)
	rs = z.ConqCollect(2).([]Color)
	assert.DeepEqual(t, rs, colors)
	rs = z.ConqCollect(1).([]Color)
	assert.DeepEqual(t, rs, colors)
}

func Test_Filter(t *testing.T) {
	z := lazy.New(colors)
	rs := z.Filter(func(c Color) bool { return c.Index%2 == 0 }).ConqCollect(6).([]Color)
	for _, c := range rs {
		assert.Assert(t, c.Index%2 == 0)
	}
	for _, c := range colors {
		if c.Index%2 == 0 {
			assert.Assert(t, rs[c.Index/2].Index == c.Index)
		}
	}
}

func Test_Map1(t *testing.T) {
	z := lazy.New([]int{0, 1, 2, 3, 4})
	rs := z.Map(func(r int) string { return fmt.Sprint(r) }).ConqCollect(6).([]string)
	assert.Assert(t, len(rs) == 5)
	for i, r := range rs {
		assert.Assert(t, r == fmt.Sprint(i))
	}
}

func Test_Map2(t *testing.T) {
	z := lazy.New(colors)
	rs := z.Map(func(r Color) string { return r.Color }).ConqCollect(6).([]string)
	assert.Assert(t, len(rs) == len(colors))
	for i, r := range rs {
		assert.Assert(t, r == colors[i].Color)
	}
}

func Test_Map3(t *testing.T) {
	type R struct{ c string }
	z := lazy.New(colors)
	rs := z.Map(func(r Color) R { return R{r.Color} }).ConqCollect(6).([]R)
	assert.Assert(t, len(rs) == len(colors))
	for i, r := range rs {
		assert.Assert(t, r.c == colors[i].Color)
	}
}

func Test_Transf(t *testing.T) {
	z := lazy.New([]int{})
	assert.Assert(t, cmp.Panics(func() {
		z.Map(func(r int) {}).ConqCollect(6)
	}))
	assert.Assert(t, cmp.Panics(func() {
		z.Filter(func(r int) {}).ConqCollect(6)
	}))
}
