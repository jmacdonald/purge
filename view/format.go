package view

import (
	"fmt"
	"strconv"
)

// Declare size interval constants.
const (
    _ = iota
    KB float64 = 1 << (10*iota)
    MB
    GB
    TB
)

// Given a size in bytes, generates a presentable string
// in the unit most appropriate for the number of bytes.
func Size(sizeInBytes int64) (formattedSize string) {

	// Constituent components of the formatted size.
	var quantity float64
	var formattedQuantity string
	var unit string

	// Convert the size argument so that we can
	// preserve decimals when converting to other units.
	size := float64(sizeInBytes)

	// Determine units, calculating quantity once decided.
	switch {
	case size < KB:
		quantity = size
		unit = "bytes"
	case size >= KB && size < MB:
		quantity = size / KB
		unit = "KB"
	case size >= MB && size < GB:
		quantity = size / MB
		unit = "MB"
	case size >= GB && size < TB:
		quantity = size / GB
		unit = "GB"
	case size >= TB:
		quantity = size / TB
		unit = "TB"
	}

	// Use no decimal places for bytes, and one for anything else.
	if unit == "bytes" {
		formattedQuantity = strconv.FormatInt(sizeInBytes, 10)
	} else {
		formattedQuantity = strconv.FormatFloat(quantity, 'f', 1, 64)
	}

	// Create the formatted size using the final quantity and unit.
	formattedSize = fmt.Sprintf("%v %v", formattedQuantity, unit)

	return
}
