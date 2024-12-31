package runners

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func (r *Runner) RestartListeningAllAddr() error {
	wg := &sync.WaitGroup{}
	addr, curr, err := r.Ucase.Repository.GetAllAddrFromDB()
	if err != nil {
		return fmt.Errorf("Error getting adresses list from DB: %v", err)
	}

	for i, _ := range addr {
		wg.Add(1)
		log.Printf("Starting listening address: %s", addr[i])
		go r.BLockListener(addr[i], curr[i])
		time.Sleep(2 * time.Second)
	}
	return nil
}
