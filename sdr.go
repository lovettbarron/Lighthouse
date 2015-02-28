// Collection of useful SDR functions
package sdr

import (
	//	"bufio"
	_ "encoding/binary"
	_ "flag"
	_ "fmt"
	_ "math"
	_ "os"
	_ "os/signal"
	_ "reflect"
	_ "strconv"
	_ "strings"
	_ "sync"
	_ "time"

	rtl "github.com/jpoirier/gortlsdr"
)

func freqHz(freqStr string) (freq uint32, err error) {
	var f64 float64
	upper := strings.ToUpper(freqStr)

	switch {
	case strings.HasSuffix(upper, "K"):
		upper = strings.TrimSuffix(upper, "K")
		f64, err = strconv.ParseFloat(upper, 64)
		freq = uint32(f64 * 1e3)
	case strings.HasSuffix(upper, "M"):
		upper = strings.TrimSuffix(upper, "M")
		f64, err = strconv.ParseFloat(upper, 64)
		freq = uint32(f64 * 1e6)
	default:
		if last := len(upper) - 1; last >= 0 {
			upper = upper[:last]
		}
		f64, err = strconv.ParseFloat(upper, 64)
		freq = uint32(f64)
	}
	return
}
