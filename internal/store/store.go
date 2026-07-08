package store

import (
	"log"
	"strings"
	"sync"

	"tailscope/internal/model"
)

type Store struct {
	mu sync.RWMutex

	capacity int

	next uint64

	logs []model.LogEntry

	subscribers map[int]chan model.LogEntry

	nextSubscriberID int
}

func New(capacity int) *Store {

	return &Store{

		capacity: capacity,

		logs: make(
			[]model.LogEntry,
			0,
			capacity,
		),

		subscribers: make(
			map[int]chan model.LogEntry,
		),
	}
}

func (s *Store) Add(
	entry model.LogEntry,
) {
	log.Printf(
		"Add(): subscribers=%d id=%d",
		len(s.subscribers),
		entry.ID,
	)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.next++

	entry.ID = s.next

	if len(s.logs) >= s.capacity {

		s.logs = s.logs[1:]
	}

	s.logs = append(
		s.logs,
		entry,
	)

	// Notify listeners

	for _, subscriber := range s.subscribers {

		select {

		case subscriber <- entry:

		default:

			// Slow clients are skipped

		}
	}
}

func (s *Store) Size() int {

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.logs)
}

func (s *Store) Latest(
	limit int,
) []model.LogEntry {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 ||
		limit > len(s.logs) {

		limit = len(s.logs)
	}

	start :=
		len(s.logs) - limit

	result :=
		make(
			[]model.LogEntry,
			limit,
		)

	copy(
		result,
		s.logs[start:],
	)

	return result
}

func (s *Store) Search(
	query Query,
) []model.LogEntry {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result :=
		make(
			[]model.LogEntry,
			0,
		)

	for i := len(s.logs) - 1; i >= 0; i-- {

		entry :=
			s.logs[i]

		if !matches(
			entry,
			query,
		) {

			continue
		}

		result =
			append(
				result,
				entry,
			)

		if query.Limit > 0 &&
			len(result) >= query.Limit {

			break
		}
	}

	return result
}

func (s *Store) Subscribe() (
	<-chan model.LogEntry,
	func(),
) {

	s.mu.Lock()
	defer s.mu.Unlock()

	id :=
		s.nextSubscriberID

	s.nextSubscriberID++

	channel :=
		make(
			chan model.LogEntry,
			100,
		)

	s.subscribers[id] =
		channel

	unsubscribe :=
		func() {

			s.mu.Lock()
			defer s.mu.Unlock()

			if ch, ok :=
				s.subscribers[id]; ok {

				delete(
					s.subscribers,
					id,
				)

				close(ch)
			}
		}

	return channel, unsubscribe
}

func matches(
	entry model.LogEntry,
	query Query,
) bool {

	if query.TraceID != "" &&
		entry.TraceID != query.TraceID {

		return false
	}

	if query.CorrelationID != "" &&
		entry.CorrelationID != query.CorrelationID {

		return false
	}

	if query.Level != "" &&
		!strings.EqualFold(
			entry.Level,
			query.Level,
		) {

		return false
	}

	if query.Pod != "" &&
		entry.Pod != query.Pod {

		return false
	}

	if query.Namespace != "" &&
		entry.Namespace != query.Namespace {

		return false
	}

	if query.Text != "" {

		search :=
			strings.ToLower(
				query.Text,
			)

		if !strings.Contains(
			strings.ToLower(
				entry.Raw,
			),
			search,
		) {

			return false
		}
	}

	return true
}
