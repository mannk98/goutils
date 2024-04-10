## How to use go sched?
```go
package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/sonnt85/gosutils/sched"
)

func main() {
	job := func() {
		t := time.Now()
		fmt.Println("Time's up! @", t.UTC())
	}
      // Run every 2 seconds but not now.
	sched.Every(2).ESeconds().NotImmediately().Run(job)
      
      // Run now and every X.
	sched.Every(5).EMinutes().Run(job)
	sched.Every().EDay().Run(job)
	sched.Every(10).WThursday().At("08:30").Run(PrintHello)
      
      // Keep the program from not exiting.
	runtime.Goexit()
}
```

## How it works?
By specifying the chain of calls, a `Job` struct is instantiated and a goroutine is starts observing the `Job`.

The goroutine will be on pause until:
* The next run scheduled is due. This will cause to execute the job.
* The `skipWait` channel is activated. This will cause to execute the job.
* The `Quit` channel is activated. This will cause to finish the job.

## Not immediate recurrent jobs
By default the behaviour of the recurrent jobs (Every(N) seconds, minutes, hours) is to start executing the job right away and then wait the required amount of time. By calling specifically `.NotImmediately()` you can override that behaviour and not execute it directly when the function `Run()` is called.

```go
sched.Every(5).EMinutes().NotImmediately().Run(job)
```

## License
Distributed under MIT license. See `LICENSE` for more information.