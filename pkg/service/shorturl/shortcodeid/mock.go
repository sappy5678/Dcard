package shortcodeid

type MockShortCodeIDRepository struct {
	NextIDFn func() string
}

func (m *MockShortCodeIDRepository) NextID() string {
	return m.NextIDFn()
}
