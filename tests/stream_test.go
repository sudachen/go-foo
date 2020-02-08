package tests

import (
	"github.com/sudachen/go-fp/lazy"
	"gotest.tools/assert"
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
	rs := z.Filter(func(c Color) bool { return c.Index%2 == 0 }).Collect().([]Color)
	for _, c := range rs {
		assert.Assert(t, c.Index%2 == 0)
	}
	for _, c := range colors {
		if c.Index%2 == 0 {
			assert.Assert(t, rs[c.Index/2].Index == c.Index)
		}
	}
}

func Test_ConqFilter(t *testing.T) {
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
