package sectors

import (
	"fmt"

	"github.com/CyrusRoshan/Golf/utils"
)

type Func struct {
	Width float64
	F     func(float64) float64
	Df    func(float64) float64
}

func RandLineFuncConstrained(minSlope float64, maxSlope float64, minHeight float64, maxHeight float64, intercept float64, width float64) Func {
	slope := utils.RandFloat(minSlope, maxSlope)

	fmt.Println("-----")
	fmt.Println(slope)
	if slope*width+intercept > maxHeight {
		slope = (maxHeight - intercept) / width
	} else if slope*width+intercept < minHeight {
		slope = (minHeight - intercept) / width
	}
	fmt.Println(slope)

	return NewLineFunc(slope, intercept)
}

func NewLineFunc(slope float64, intercept float64) Func {
	f := func(x float64) float64 {
		return slope*x + intercept
	}
	df := func(x float64) float64 {
		return slope
	}

	return Func{
		F:  f,
		Df: df,
	}
}
