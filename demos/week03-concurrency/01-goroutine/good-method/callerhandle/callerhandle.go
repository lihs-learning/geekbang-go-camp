package callerhandle

import (
"fmt"
"time"
)

type Some struct {}

// ListDirectory leave concurrency to the caller
func (s Some) ListDirectory(dir string, handle func(string, error)(stop bool)) {
	fmt.Println("listing", dir)
	time.Sleep(5 * time.Minute)
	for handle("some", nil) {}
}

//func Demo() {
//	cnt := 0
//	foo := Some{}
//	go foo.ListDirectory("/path/to", func(s string, err error) (stop bool) {
//		if err != nil {
//			return true
//		}
//		fmt.Println(s)
//		cnt += 1
//		if cnt > 5 {
//			return true
//		}
//		return false
//	})
//  // do other work
//}
