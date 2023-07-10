package main

import (
	"image/color"

	"golang.org/x/image/colornames"
)

// ANCHOR Element Map
var ElementMap = map[int]Element{
	0: {
		Color:   colornames.Black,
		Name:    "Void",
		Density: 0,
		isFluid: true,
	},
	1: {
		Color:   colornames.Lightsteelblue,
		Name:    "Hydrogen",
		Density: 1,
		isFluid: true,
	},
	6: {
		Color:   colornames.Gray,
		Name:    "Carbon",
		Density: 22,
		isFluid: false,
	},
	8: {
		Color:   colornames.Aqua,
		Name:    "Oxygen",
		Density: 1,
		isFluid: true,
	},
	14: {
		Color:   colornames.Red,
		Name:    "Silicon",
		Density: 24,
		isFluid: false,
	},
	22: {
		Color:   colornames.Cornflowerblue,
		Name:    "Titanium",
		Density: 45,
		isFluid: false,
	},
	80: {
		Color:   colornames.White,
		Name:    "Mercury",
		Density: 135,
		isFluid: true,
	},
}

// ANCHOR Element Struct
type Element struct {
	Name    string
	Color   color.RGBA
	Density int
	isFluid bool
}
