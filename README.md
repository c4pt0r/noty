noty
====

go notify util

```
package main

import (
	"github.com/c4pt0r/noty"
	log "github.com/ngaut/logging"
)

func init() {
	noty.Register("test", func(p interface{}) (interface{}, error) {
		log.Info(p)
		return p, nil
	})
}

func main() {
	r, err := noty.SignalSync("test", "param")
	if err != nil {
		log.Error(err)
	}
	log.Info(r)
}
```
