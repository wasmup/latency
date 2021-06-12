```sh
go run .

5203228
5199442
5202763
5223122
5187032
5247427
5228973
5192881
n 4977017
min 58ns
ave 1.89µs
max 3.559549ms
```

```go
package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
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

	for ctx.Err() == nil {
		d := time.Since(t0)
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
```

---

```sh
go run .

5650519
5571047
5739384
5686872
5668563
5694935
5719850
5718231
n 5161255
min 23ns
ave 58ns
max 306.238µs
```


```go
package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	go loadAll(ctx, &wg)

	var min, max, ave time.Duration
	n := 0

	for ctx.Err() == nil {
		t0 := time.Now()
		d := time.Since(t0)
		ave += d
		n++
		if n == 1 || d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}

	wg.Wait()
	ave /= time.Duration(n)

	fmt.Println("n", n)
	fmt.Println("min", min)
	fmt.Println("ave", ave)
	fmt.Println("max", max)
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
```