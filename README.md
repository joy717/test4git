# pool-async
pool-async is a tool for goroutines with a pool.

#Quick start
```
package main

import "github.com/joy717/pool-async/utils"

func main() {
  pa := utils.NewDefaultPoolAsync()
  pa.DoWitError(func() error {
    fmt.Println("goroutine 1")
    return nil
  }).DoWitError(func() error {
    fmt.Println("goroutine 2")
    return fmt.Errorf("goroutine 2 err")
  })
  
  if err := pa.Wait(); err != nil {
    t.Println("the 1st non-nil err: ", err)
  }
}
```
