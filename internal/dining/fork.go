package dining

import "sync"

// Fork represents a shared resource protected by a mutual exclusion lock (Mutex).
// Embedding sync.Mutex allows calling Lock() and Unlock() directly on the Fork struct.
type Fork struct {
	ID int
	sync.Mutex
}
