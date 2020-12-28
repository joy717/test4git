# pool-async
pool-async is a tool for goroutines with a pool.

#Quick start
```
package main

import "github.com/joy717/poolasync"
import "fmt"

func main() {
  pa := poolasync.NewDefaultPoolAsync()
  pa.DoWitError(func() error {
    fmt.Println("goroutine 1")
    return nil
  }).DoWitError(func() error {
    fmt.Println("goroutine 2")
    return fmt.Errorf("goroutine 2 err")
  })
  
  // NOTICE: should call pa.Wait() always to complete the jobs.
  // because we use channel, Wait() is the reciver function.
  if err := pa.Wait(); err != nil {
    fmt.Println("the 1st non-nil err: ", err)
  }
}
```
If you like this tool, star it to make this tool to be known by more people.
