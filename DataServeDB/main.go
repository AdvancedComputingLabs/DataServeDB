package main

import (
	"fmt"
)

//TODO:
// 1) use db's own conversion functions [done; sort of, see note], note: for UInt conversion in Length parsing db's convert package is not used as UInt conversion function is not there.
// 2) add weakconversion to field properties [done]
// 3) DbNull conversion test

//Table Requirements:
// - Fields should have internal ids, it will help changing name of the field.
// - Fields should have positions so they can be rearranged.
// - Table fields meta data should not update when read or writes are happening to the table data. Reader lock on meta data?

//Next Tasks:
// - Create table and send data; see if it is processed correctly.

// Story:
// - Data comes in json format, which is table record.
// - It is converted to correct Go object(s) using table metadata which includes fields metadata.
// - Table name is in the url.

// - Table fields are not stored by name, but by internal id. This allows to easily change the name of the field.
// - Field mapper takes the name of the field and returns its meta data.

//TODO: currently if field name is wrong it fails, but if certain fields are not in input then it will go through unless it is not nullable.
// Check if this correct behavior with regards to previous db server?

//TODO: all dbtype conversions are weak conversions, change it later to choice between weak and strong conversions.
//weak conversion will convert to close result between different types or values e.g. 1 = true, 0 = false, "TRUE" = true.
//Strong conversion will require strong typing, "true" will not convert to true.

//TODO: write docs state 'Zen of DataserveDB'; inspired from https://en.wikipedia.org/wiki/Zen_of_Python
//1. KISS principle:
//1.a. Few types to keep data design simple. Some will take more memory (e.g. int32 than int8) but there is more RAM and disk space these days.

func main() {
	fmt.Println("Hello World!")
}
