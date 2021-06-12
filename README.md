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


---

## Histogram

```sh
go run . 
5375653
5512829
5435067
5461223
5514545
5382318
5582644
5486227
n 5229320
min 81ns
ave 1.785µs
max 4.258098ms

# Histogram

81ns-42.661µs          99.3%      █████▏  5190278
42.661µs-85.241µs      0.625%     ▏       32706
85.241µs-127.821µs     0.0908%    ▏       4747
127.821µs-170.401µs    0.021%     ▏       1099
170.401µs-212.981µs    0.00551%   ▏       288
212.981µs-255.561µs    0.00142%   ▏       74
255.561µs-298.141µs    0.000516%  ▏       27
298.141µs-340.721µs    0.000287%  ▏       15
340.721µs-383.301µs    0.000287%  ▏       15
383.301µs-425.881µs    7.65e-05%  ▏       4
425.881µs-468.461µs    0.000134%  ▏       7
468.461µs-511.041µs    7.65e-05%  ▏       4
511.041µs-553.621µs    3.82e-05%  ▏       2
553.621µs-596.201µs    5.74e-05%  ▏       3
596.201µs-638.781µs    9.56e-05%  ▏       5
638.781µs-681.361µs    7.65e-05%  ▏       4
681.361µs-723.941µs    3.82e-05%  ▏       2
723.941µs-766.521µs    3.82e-05%  ▏       2
766.521µs-809.101µs    3.82e-05%  ▏       2
809.101µs-851.681µs    3.82e-05%  ▏       2
851.681µs-894.261µs    5.74e-05%  ▏       3
894.261µs-936.841µs    5.74e-05%  ▏       3
979.421µs-1.022001ms   5.74e-05%  ▏       3
1.064581ms-1.107161ms  1.91e-05%  ▏       1
1.107161ms-1.149741ms  3.82e-05%  ▏       2
1.277481ms-1.320061ms  1.91e-05%  ▏       1
1.320061ms-1.362641ms  1.91e-05%  ▏       1
1.362641ms-1.405221ms  1.91e-05%  ▏       1
1.405221ms-1.447801ms  1.91e-05%  ▏       1
1.447801ms-1.490381ms  1.91e-05%  ▏       1
1.490381ms-1.532961ms  1.91e-05%  ▏       1
1.532961ms-1.575541ms  1.91e-05%  ▏       1
2.171661ms-2.214241ms  1.91e-05%  ▏       1
2.214241ms-2.256821ms  1.91e-05%  ▏       1
2.256821ms-2.299401ms  1.91e-05%  ▏       1
2.427141ms-2.469721ms  1.91e-05%  ▏       1
2.512301ms-2.554881ms  1.91e-05%  ▏       1
2.725201ms-2.767781ms  3.82e-05%  ▏       2
2.852941ms-2.895521ms  3.82e-05%  ▏       2
2.938101ms-2.980681ms  1.91e-05%  ▏       1
3.023261ms-3.065841ms  1.91e-05%  ▏       1
3.236161ms-3.278741ms  3.82e-05%  ▏       2
3.406481ms-3.449061ms  1.91e-05%  ▏       1
4.215501ms-4.258081ms  1.91e-05%  ▏       1
```