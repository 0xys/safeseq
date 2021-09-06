package sequencer

import (
	"sync"

	"github.com/0xys/safeseq/models"
)

type Transactions []*models.Transaction

type SortedWaitlist struct {
	AccountId  string
	dict       map[uint64]*models.Transaction
	txs        Transactions
	beginIndex int
}

// add new transaction to sequence
func (l *SortedWaitlist) Add(transaction *models.Transaction) bool {
	if l.dict[transaction.Nonce] != nil {
		return false
	}

	l.txs = append(l.txs, transaction)
	l.txs.Sort(l.beginIndex)
	return true
}

func (l *SortedWaitlist) Pop() (bool, *models.Transaction) {
	if len(l.dict) == 0 {
		return false, nil
	}

	tx := models.CopyTransaction(l.txs[l.beginIndex])
	delete(l.dict, l.txs[l.beginIndex].Nonce)
	l.beginIndex++
	return true, tx
}

func (l *SortedWaitlist) Len() int {
	return len(l.dict)
}

func (t Transactions) Sort(begin int) {
	quicksort(&t, begin, len(t)-1)
}

func quicksort(array *Transactions, low int, high int) {
	if low < high {
		pi := partition(array, low, high)

		quicksort(array, low, pi-1)
		quicksort(array, pi+1, high)
	}
}

func partition(array *Transactions, low int, high int) int {
	pivot := (*array)[high]

	start := low - 1

	for i := 0; i < len(*array); i++ {
		if (*array)[i].Nonce < pivot.Nonce {
			start++
			(*array)[start], (*array)[i] = (*array)[i], (*array)[start]
		}
	}
	(*array)[start+1], (*array)[high] = (*array)[high], (*array)[start+1]
	return start + 1
}

type Sequencer struct {
	lock        sync.Mutex
	Waitlists   map[string]*SortedWaitlist // account id to waitlist mapping
	SubmitQueue chan *models.Transaction

	SuccessLists map[string]*models.Transaction
	FailureLists map[string]*models.Transaction
}

func NewSequencer() *Sequencer {
	return &Sequencer{
		Waitlists:   make(map[string]*SortedWaitlist),
		SubmitQueue: make(chan *models.Transaction),
	}
}

func (s *Sequencer) Add(account string, tx *models.Transaction) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	res := s.Waitlists[account].Add(tx)
	if !res {
		return false
	}

	if len(s.SubmitQueue) > 10 {
		return false
	}

	found, next := s.Waitlists[account].Pop()
	if !found || next == nil {
		return false
	}
	s.SubmitQueue <- next

	return true
}
