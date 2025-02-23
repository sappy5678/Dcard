package shortcode

type MockShortCodeIDRepository struct {
	NextIDFunc func() string
}

func (m *MockShortCodeIDRepository) NextID() string {
	return m.NextIDFunc()
}
