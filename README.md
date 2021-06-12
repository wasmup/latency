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

6044097
5813540
6140386
5834088
5832977
5845063
5665101
5763578
n 5464038
min 33ns
ave 56ns
max 270.645µs
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