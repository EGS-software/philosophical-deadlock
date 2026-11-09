package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/EGS-software/philosophical-deadlock/internal/dining"
	"github.com/EGS-software/philosophical-deadlock/internal/logger"
)

const numPhilosophers = 5

func main() {
	// Parse CLI flags to switch modes seamlessly during the presentation
	modeFlag := flag.String("mode", "safe", "Simulation mode: 'safe' (prevents deadlock) or 'deadlock' (forces circular wait)")
	mealsFlag := flag.Int("meals", 3, "Number of meals each philosopher must eat before finishing")
	flag.Parse()

	mode := dining.Mode(*modeFlag)
	if mode != dining.ModeSafe && mode != dining.ModeDeadlock {
		fmt.Println("Invalid mode. Use --mode=safe or --mode=deadlock")
		return
	}

	fmt.Printf("=== Starting Dining Philosophers Simulation ===\n")
	fmt.Printf("Mode: %s | Philosophers: %d | Meals per Philosopher: %d\n\n", mode, numPhilosophers, *mealsFlag)

	log := logger.New()

	// Initialize shared resources (5 Forks)
	forks := make([]*dining.Fork, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = &dining.Fork{ID: i + 1}
	}

	// Initialize active entities (5 Philosophers)
	philosophers := make([]*dining.Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &dining.Philosopher{
			ID:        i + 1,
			LeftFork:  forks[i],
			RightFork: forks[(i+1)%numPhilosophers], // Modulo operator connects the last philosopher to the first fork (circular table)
			Mode:      mode,
			MealsLeft: *mealsFlag,
			Logger:    log,
		}
	}

	// WaitGroup waits for all philosopher goroutines to finish
	var wg sync.WaitGroup
	wg.Add(numPhilosophers)

	// Spawn a goroutine for each philosopher
	for _, p := range philosophers {
		go p.Dine(&wg)
	}

	// Block main thread until all philosophers finish eating (or deadlock occurs)
	wg.Wait()

	// Small delay to ensure the async logger flushes the final messages to stdout
	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n=== Simulation Completed Successfully! No Deadlocks Occurred. ===")
}
