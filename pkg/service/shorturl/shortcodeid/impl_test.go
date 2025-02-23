package shortcodeid

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	impl *impl
}

func (ts *TestSuite) SetupSuite() {
	ts.impl = New(1).(*impl)
}

func (ts *TestSuite) TestEncodeID() {
	tests := []struct {
		name    string
		input   uint64
		want    string
		wantErr bool
	}{
		{
			name:  "normal conversion",
			input: 12345,
			want:  "5ob", // 假設使用 base62 編碼
		},
		{
			name:    "zero value",
			input:   0,
			want:    "2",
			wantErr: true,
		},
		{
			name:  "boundary value 1",
			input: 1,
			want:  "3",
		},
	}

	for _, tc := range tests {
		ts.Run(tc.name, func() {
			got := ts.impl.encodeID(tc.input)
			ts.Require().Equal(tc.want, got)

			num, err := ts.impl.decodeID(got)
			ts.Require().NoError(err)
			ts.Require().Equal(tc.input, num)
		})
	}
}

func (ts *TestSuite) TestNextID() {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "normal case",
			want: "3-3",
		},
		{
			name: "next case",
			want: "3-4",
		},
		{
			name: "next case",
			want: "3-5",
		},
	}

	for _, tc := range tests {
		ts.Run(tc.name, func() {
			got := ts.impl.NextID()
			ts.Require().NotEmpty(got)
			ts.Require().Equal(tc.want, got)
		})
	}
}

func TestShortCodeIDSuite(t *testing.T) {
	ts := new(TestSuite)
	suite.Run(t, ts)
}
