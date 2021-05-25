package inmemory

const (
	emptyValue string = ""
)

// Store implements inmemory transaction storage.
//
// It MUST be used from single go-routine (it's not go-routine safe).
type Store struct {
	// transactions is a stack of currently open transactions.
	// Zero element transactions[0] is special, it's not issued by the user
	// (while all the transactions building on top of this one are)
	// and it represents the root state.
	// transactions[openTxCnt] represents transaction that gets currently targeted
	// by the commands issued on Store.
	transactions []transaction
	// openTxCnt represents the amount of currently open (nested) transactions.
	openTxCnt int
}

// NewStore returns new instance.
func NewStore() *Store {
	return &Store{
		transactions: []transaction{
			{
				state: make(map[string]string),
			},
		},
	}
}

func (s *Store) Set(key, value string) {
	s.transactions[s.openTxCnt].state[key] = value
}

func (s *Store) Get(key string) string {
	for txIdx := s.openTxCnt; txIdx >= 0; txIdx-- {
		value, ok := s.transactions[txIdx].state[key]
		if !ok {
			continue
		}
		return value
	}
	return emptyValue
}

func (s *Store) Delete(key string) {
	// Using empty value "" to denote deleted element.
	s.transactions[s.openTxCnt].state[key] = emptyValue
}

func (s *Store) Count(value string) int {
	matches := make(map[string]struct{})

	// Start searching for matches with the root transaction, and then apply other layers on top of one another.
	for _, tx := range s.transactions {
		for k, v := range tx.state {
			// Previously matched element got overwritten by this transaction, its value has changed,
			// so it is no longer a match.
			delete(matches, k)

			if v == value {
				matches[k] = struct{}{}
			}
		}
	}

	return len(matches)
}

func (s *Store) Begin() {
	tx := transaction{
		state: make(map[string]string),
	}
	s.transactions = append(s.transactions, tx)

	s.openTxCnt++
}

func (s *Store) Commit() {
	if s.openTxCnt == 0 {
		return
	}

	// Commit the state of the current transaction to its parent.
	for k, v := range s.transactions[s.openTxCnt].state {
		s.transactions[s.openTxCnt-1].state[k] = v
	}

	s.openTxCnt--

	// This implementation relies on empty value "" denoting deleted element.
	// Over time items with empty values will accumulate, so we need to prune these periodically.
	// The easiest way to prune it is to wait until there are no outstanding transactions issued on Store.
	//
	// This assumes such a chance will arrive eventually from time to time.
	if s.openTxCnt == 0 {
		for k, v := range s.transactions[s.openTxCnt].state {
			if v == emptyValue {
				delete(s.transactions[s.openTxCnt].state, k)
			}
		}
	}
}

func (s *Store) Rollback() {
	if s.openTxCnt == 0 {
		return
	}

	s.transactions = s.transactions[:s.openTxCnt]

	s.openTxCnt--
}

type transaction struct {
	// state represents state this transaction has accumulated.
	state map[string]string
}
