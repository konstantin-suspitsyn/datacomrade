package shared

import (
	"fmt"
	"sync"

	"github.com/konstantin-suspitsyn/datacomrade/internal/utils/jsonlog"
)

func BackroundJob(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				jsonlog.PrintError(fmt.Errorf("%s", err).Error(), nil, nil)
			}

		}()
		fn()
	}()
}
