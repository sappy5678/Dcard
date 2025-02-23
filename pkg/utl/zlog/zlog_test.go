package zlog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sappy5678/dcard/pkg/utl/zlog"
)

func TestNew(t *testing.T) {
	log := zlog.New()
	assert.NotNil(t, log)
}

func TestLog(t *testing.T) {
	log := zlog.New()

	assert.NotPanics(t, func() {
		log.Log(context.Background(), "test", "test", nil, nil)
	})
	assert.NotPanics(t, func() {
		log.Log(context.Background(), "test", "test", errors.New("test"), nil)
	})
	assert.NotPanics(t, func() {
		log.Log(context.Background(), "test", "test", errors.New("test"), map[string]interface{}{"test": "test"})
	})
}
