package runners

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//USING THIS RUNNER FOR LISTENING TX FROM ALL ADRESSES WE HAVE WHEN WE RESTART THE APP

func (r *Runner) RestartListeningAllAddr() error {
	wg := &sync.WaitGroup{}
	addr, curr, err := r.Ucase.Repository.GetAllAddr()
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
