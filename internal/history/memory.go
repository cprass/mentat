package history

type InMemoryStore struct {
	events []*ReviewEvent
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		events: make([]*ReviewEvent, 0),
	}
}

func (m *InMemoryStore) LoadAll() ([]*ReviewEvent, error) {
	return m.events, nil
}

func (m *InMemoryStore) Append(event *ReviewEvent) error {
	m.events = append(m.events, event)
	return nil
}

func (m *InMemoryStore) Close() error {
	m.events = make([]*ReviewEvent, 0)
	return nil
}
