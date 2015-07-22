package tetra

import (
	"log"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

func startWorkers(num int) {
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			uid := uuid.New()

			defer func() {
				if r := recover(); r != nil {
					wg.Done()

					log.Printf("Recovered in %s", uid)
					startWorkers(1)
				}
			}()

			debugf("Worker %s started", uid)

			for line := range tasks {
				ProcessLine(line)
				time.Sleep(5 * time.Millisecond)
			}
			wg.Done()
		}()
	}
}
