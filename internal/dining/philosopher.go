package dining

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/EGS-software/philosophical-deadlock/internal/logger"
)

// Mode defines the concurrency strategy used in the simulation
type Mode string

const (
	ModeDeadlock Mode = "deadlock"
	ModeSafe     Mode = "safe"
)

// Philosopher represents an independent actor (goroutine) competing for resources
type Philosopher struct {
	ID        int
	LeftFork  *Fork
	RightFork *Fork
	Mode      Mode
	MealsLeft int
	Logger    *logger.Logger
}

// Dine starts the philosopher's lifecycle (Think -> Hungry -> Eat)
func (p *Philosopher) Dine(wg *sync.WaitGroup) {
	defer wg.Done()

	for p.MealsLeft > 0 {
		p.think()
		p.eat()
		p.MealsLeft--
	}

	p.Logger.Log(p.ID, "✅ Finished all meals and left the table.")
}

func (p *Philosopher) think() {
	p.Logger.Log(p.ID, "💭 Is thinking...")
	if p.Mode == ModeDeadlock {
		// Short deterministic delay to align goroutines and force the deadlock condition quickly
		time.Sleep(50 * time.Millisecond)
	} else {
		// Randomized thinking time for realistic concurrency simulation
		time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)
	}
}

func (p *Philosopher) eat() {
	p.Logger.Log(p.ID, "🍽️  Is hungry and trying to pick up forks.")

	p.pickUpForks()

	// Critical Section: The philosopher holds both locks here
	p.Logger.Log(p.ID, fmt.Sprintf("🍝 Is EATING (holding forks #%d and #%d). Meals left: %d", p.LeftFork.ID, p.RightFork.ID, p.MealsLeft-1))
	time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)

	p.putDownForks()
}

// pickUpForks handles resource acquisition based on the selected concurrency strategy
func (p *Philosopher) pickUpForks() {
	if p.Mode == ModeDeadlock {
		// =========================================================================
		// DEADLOCK SCENARIO: Circular Wait + Hold and Wait (Coffman Conditions)
		// =========================================================================
		// Every philosopher acquires their LeftFork first.
		p.LeftFork.Lock()
		p.Logger.Log(p.ID, fmt.Sprintf("🔒 Picked up LEFT fork #%d. Waiting for RIGHT fork #%d...", p.LeftFork.ID, p.RightFork.ID))

		// Artificial delay while holding the first lock.
		// This guarantees all 5 philosophers grab their left fork simultaneously,
		// creating an unbreakable circular wait for the right fork.
		time.Sleep(200 * time.Millisecond)

		p.RightFork.Lock()
		p.Logger.Log(p.ID, fmt.Sprintf("🔒 Picked up RIGHT fork #%d.", p.RightFork.ID))
		return
	}

	// =========================================================================
	// SAFE SCENARIO: Resource Hierarchy Solution (Dijkstra's Asymmetric Approach)
	// =========================================================================
	// To break the Circular Wait condition, philosophers must acquire forks in a
	// strict global order (always lock the lowest ID fork first, then the highest).
	// For Philosophers 1 to 4: LeftFork has a lower ID than RightFork.
	// For Philosopher 5: RightFork (#1) has a lower ID than LeftFork (#5).
	firstFork, secondFork := p.LeftFork, p.RightFork
	if p.RightFork.ID < p.LeftFork.ID {
		firstFork, secondFork = p.RightFork, p.LeftFork
	}

	firstFork.Lock()
	p.Logger.Log(p.ID, fmt.Sprintf("🔒 Picked up first fork #%d (lowest ID).", firstFork.ID))

	secondFork.Lock()
	p.Logger.Log(p.ID, fmt.Sprintf("🔒 Picked up second fork #%d (highest ID).", secondFork.ID))
}

func (p *Philosopher) putDownForks() {
	p.LeftFork.Unlock()
	p.RightFork.Unlock()
	p.Logger.Log(p.ID, fmt.Sprintf("🔓 Put down forks #%d and #%d.", p.LeftFork.ID, p.RightFork.ID))
}
