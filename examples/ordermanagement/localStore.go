package main

type itemDetails struct {
	quantity     int
	pricePerItem float64
}

// Local map storing item inventory details. This can be retrieved from external store
var itemInventory = map[string]itemDetails{
	"pen":      {125, 5.5},
	"pencil":   {150, 2.75},
	"clips":    {500, 1.35},
	"stapler":  {300, 7.25},
	"notepad":  {100, 8.15},
	"marker":   {300, 2.25},
	"eraser":   {450, 1.75},
	"sharpner": {350, 2.25},
}

// Local map storing customer level discounts. This again can be maintained from external store
var levelDiscount = map[string]float64{
	"diamond": 20,
	"gold":    15,
	"silver":  10,
	"bronze":  5,
}
