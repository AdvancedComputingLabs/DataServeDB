package main

import (
	"fmt"
)

// Concept:
//TODO: write docs state 'Zen of DataserveDB'; inspired from https://en.wikipedia.org/wiki/Zen_of_Python
//1. KISS principle:
//1.a. Few types to keep data design simple. Some will take more memory (e.g. int32 than int8) but there is more RAM and disk space these days.

// NOTES:
// 1)
//weak conversion will convert to close result between different types or values e.g. 1 = true, 0 = false, "TRUE" = true.
//Strong conversion will require strong typing, "true" will not convert to true.

//Table Requirements:
// - Fields should have internal ids, it will help changing name of the field.
// - Fields should have positions so they can be rearranged.
// - Table fields meta data should not update when read or writes are happening to the table data. Reader lock on meta data?

//Next Tasks:
// 1)
//TODO: currently if field name is wrong it fails, but if certain fields are not in input then it will go through unless it is not nullable.
// Check if this correct behavior with regards to previous db server?

// 2)
//TODO: all dbtype conversions are weak conversions, change it later to choice between weak and strong conversions.

// 3) file names and paths are in packages; make central location and maybe package to handle file paths.
// 4)

func main() {
	fmt.Println("Hello World!")
}
