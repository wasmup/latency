package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	go loadAll(ctx, &wg)

	var min, max, ave time.Duration
	t0 := time.Now()
	n := 0

	a := make([]time.Duration, 0, 10_000_000)
	for ctx.Err() == nil {
		d := time.Since(t0)
		a = append(a, d)
		ave += d
		n++
		if n == 1 || d < min {
			min = d
		}
		if d > max {
			max = d
		}
		t0 = time.Now()
	}

	wg.Wait()
	ave /= time.Duration(n)

	fmt.Println("n", n)
	fmt.Println("min", min)
	fmt.Println("ave", ave)
	fmt.Println("max", max)

	fmt.Println()
	hist := Hist(100, a)
	maxWidth := 5
	err := Fprint(os.Stdout, hist, Linear(maxWidth))
	if err != nil {
		log.Fatal(err)
	}

}

func loadAll(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			var j uint64
			for {
				select {
				case <-ctx.Done():
					fmt.Println(j)
					wg.Done()
					return
				default:
					j++
				}
			}
		}()
	}
}

// https://github.com/aybabtme/uniplot

func Fprint(w io.Writer, h Histogram, s ScaleFunc) error {
	return fprintf(w, h, s, func(v time.Duration) string {
		return fmt.Sprint(v)
	})
}

type FormatFunc func(v time.Duration) string

var blocks = []string{
	"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█",
}

var barstring = func(v float64) string {
	decimalf := (v - math.Floor(v)) * 10.0
	decimali := math.Floor(decimalf)
	charIdx := int(decimali / 10.0 * 8.0)
	return strings.Repeat("█", int(v)) + blocks[charIdx]
}

func fprintf(w io.Writer, h Histogram, s ScaleFunc, f FormatFunc) error {
	tabw := tabwriter.NewWriter(w, 2, 2, 2, byte(' '), 0)

	yfmt := func(y int) string {
		if y > 0 {
			return strconv.Itoa(y)
		}
		return ""
	}

	for i, bkt := range h.Buckets {
		if bkt.Count == 0 {
			continue
		}
		sz := h.Scale(s, i)
		fmt.Fprintf(tabw, "%s-%s\t%.3g%%\t%s\n",
			f(bkt.Min), f(bkt.Max),
			float64(bkt.Count)*100.0/float64(h.Count),
			barstring(sz)+"\t"+yfmt(bkt.Count),
		)
	}

	return tabw.Flush()
}

// Histogram holds a count of values partionned over buckets.
type Histogram struct {
	// Min is the size of the smallest bucket.
	Min int
	// Max is the size of the biggest bucket.
	Max int
	// Count is the total size of all buckets.
	Count int
	// Buckets over which values are partionned.
	Buckets []Bucket
}

// Bucket counts a partion of values.
type Bucket struct {
	// Count is the number of values represented in the bucket.
	Count int
	// Min is the low, inclusive bound of the bucket.
	Min time.Duration
	// Max is the high, exclusive bound of the bucket. If
	// this bucket is the last bucket, the bound is inclusive
	// and contains the max value of the histogram.
	Max time.Duration
}

// Hist creates an histogram partionning input over `bins` buckets.
func Hist(bins int, input []time.Duration) Histogram {
	if len(input) == 0 || bins == 0 {
		return Histogram{}
	}

	min, max := input[0], input[0]
	for _, v := range input {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	if min == max {
		return Histogram{
			Min:     len(input),
			Max:     len(input),
			Count:   len(input),
			Buckets: []Bucket{{Count: len(input), Min: min, Max: max}},
		}
	}

	scale := (max - min) / time.Duration(bins)
	buckets := make([]Bucket, bins)
	for i := range buckets {
		bmin, bmax := time.Duration(i)*scale+min, time.Duration(i+1)*scale+min
		buckets[i] = Bucket{Min: bmin, Max: bmax}
	}

	minC, maxC := 0, 0
	for _, val := range input {
		minx := time.Duration(min)
		xdiff := val - minx
		bi := imin(int(xdiff/scale), len(buckets)-1)
		if bi < 0 || bi >= len(buckets) {
			log.Panicf("bi=%d\tval=%v\txdiff=%v\tscale=%v\tlen(buckets)=%d", bi, val, xdiff, scale, len(buckets))
		}
		buckets[bi].Count++
		minC = imin(minC, buckets[bi].Count)
		maxC = imax(maxC, buckets[bi].Count)
	}

	return Histogram{
		Min:     minC,
		Max:     maxC,
		Count:   len(input),
		Buckets: buckets,
	}
}

// PowerHist creates an histogram partionning input over buckets of power
// `pow`.
// func PowerHist(power float64, input []float64) Histogram {
// 	if len(input) == 0 || power <= 0 {
// 		return Histogram{}
// 	}

// 	minx, maxx := input[0], input[0]
// 	for _, val := range input {
// 		minx = math.Min(minx, val)
// 		maxx = math.Max(maxx, val)
// 	}

// 	fromPower := math.Floor(logbase(minx, power))
// 	toPower := math.Floor(logbase(maxx, power))

// 	buckets := make([]Bucket, int(toPower-fromPower)+1)
// 	for i, bkt := range buckets {
// 		bkt.Min = math.Pow(power, float64(i)+fromPower)
// 		bkt.Max = math.Pow(power, float64(i+1)+fromPower)
// 		buckets[i] = bkt
// 	}

// 	minC := 0
// 	maxC := 0
// 	for _, val := range input {
// 		powAway := logbase(val, power) - fromPower
// 		bi := int(math.Floor(powAway))
// 		buckets[bi].Count++
// 		minC = imin(buckets[bi].Count, minC)
// 		maxC = imax(buckets[bi].Count, maxC)
// 	}

// 	return Histogram{
// 		Min:     minC,
// 		Max:     maxC,
// 		Count:   len(input),
// 		Buckets: buckets,
// 	}
// }

// Scale gives the scaled count of the bucket at idx, using the
// provided scale func.
func (h Histogram) Scale(s ScaleFunc, idx int) float64 {
	bkt := h.Buckets[idx]
	scale := s(h.Min, h.Max, bkt.Count)
	return scale
}

// ScaleFunc is the type to implement to scale an histogram.
type ScaleFunc func(min, max, value int) float64

// Linear builds a ScaleFunc that will linearly scale the values of
// an histogram so that they do not exceed width.
func Linear(width int) ScaleFunc {
	return func(min, max, value int) float64 {
		if min == max {
			return 1
		}
		return float64(value-min) / float64(max-min) * float64(width)
	}
}

// func logbase(a, base float64) float64 {
// 	return math.Log2(a) / math.Log2(base)
// }

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
