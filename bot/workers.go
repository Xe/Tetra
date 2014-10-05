package tetra

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

func (tetra *Tetra) startWorkers(num int) {
	for i := 0; i < num; i++ {
		tetra.wg.Add(1)
		go func() {
			uid := uuid.New()
			debugf("Worker %s started", uid)

			for line := range tetra.tasks {
				tetra.ProcessLine(line)
				time.Sleep(5 * time.Millisecond)
			}
			tetra.wg.Done()
		}()
	}
}
