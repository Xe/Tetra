package tetra

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

func startWorkers(num int) {
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			uid := uuid.New()
			debugf("Worker %s started", uid)

			for line := range tasks {
				ProcessLine(line)
				time.Sleep(5 * time.Millisecond)
			}
			wg.Done()
		}()
	}
}
