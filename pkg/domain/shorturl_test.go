package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortURL_IsValid(t *testing.T) {
	testCases := []struct {
		name     string
		shortURL ShortURL
		currTime uint64
		expected bool
	}{
		{
			name: "valid short URL",
			shortURL: ShortURL{
				ShortCode:   "abc123",
				OriginalURL: "https://example.com",
				ShortURL:    "https://short.url/abc123",
				ExpireTime:  2000,
				CreatedTime: 1000,
			},
			currTime: 1500,
			expected: true,
		},
		{
			name: "expired short URL",
			shortURL: ShortURL{
				ShortCode:   "abc123",
				OriginalURL: "https://example.com",
				ShortURL:    "https://short.url/abc123",
				ExpireTime:  1500,
				CreatedTime: 1000,
			},
			currTime: 2000,
			expected: false,
		},
		{
			name: "invalid original URL",
			shortURL: ShortURL{
				ShortCode:   "abc123",
				OriginalURL: "not-a-url",
				ShortURL:    "https://short.url/abc123",
				ExpireTime:  2000,
				CreatedTime: 1000,
			},
			currTime: 1500,
			expected: false,
		},
		{
			name: "empty original URL",
			shortURL: ShortURL{
				ShortCode:   "abc123",
				OriginalURL: "not-a-url",
				ShortURL:    "https://short.url/abc123",
				ExpireTime:  2000,
				CreatedTime: 1000,
			},
			currTime: 1500,
			expected: false,
		},
		{
			name: "empty short code",
			shortURL: ShortURL{
				ShortCode:   "",
				OriginalURL: "https://example.com",
				ShortURL:    "https://short.url/abc123",
				ExpireTime:  2000,
				CreatedTime: 1000,
			},
			currTime: 1500,
			expected: false,
		},
		{
			name: "invalid expired",
			shortURL: ShortURL{
				ShortCode:   "abc123",
				OriginalURL: "https://example.com",
				ShortURL:    "not-a-url",
				ExpireTime:  1000,
				CreatedTime: 2000,
			},
			currTime: 100,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.shortURL.IsValid(tc.currTime)
			assert.Equal(t, tc.expected, result)
		})
	}
}
