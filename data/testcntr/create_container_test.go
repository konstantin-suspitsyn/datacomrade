package testcntr

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx)
	assert.NoError(t, err)
}
