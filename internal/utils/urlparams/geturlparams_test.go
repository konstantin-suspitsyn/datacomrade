package urlparams

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/konstantin-suspitsyn/datacomrade/configs"
	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?param=value", nil)
		result, err := GetString("param", req)
		assert.NoError(t, err)
		assert.Equal(t, "value", result)
	})

	t.Run("MissingParam", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?other=1", nil)
		result, err := GetString("param", req)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoParameterInUrl))
		assert.Empty(t, result)
	})
}

func TestGetInt(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?param=123", nil)
		result, err := GetInt("param", req)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), result)
	})

	t.Run("MissingParam", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?other=1", nil)
		result, err := GetInt("param", req)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoParameterInUrl))
		assert.Equal(t, int64(-1), result)
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?param=abc", nil)
		result, err := GetInt("param", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ERROR converting string to int64")
		assert.Equal(t, int64(-1), result)
	})

	t.Run("NegativeValue", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?param=-5", nil)
		result, err := GetInt("param", req)
		assert.NoError(t, err)
		assert.Equal(t, int64(-5), result)
	})
}

func TestGetPager(t *testing.T) {
	t.Run("DefaultsWhenMissing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		pager, err := GetPager(req)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), pager.CurrentPage)
		assert.Equal(t, 0, pager.PageSize)
	})

	t.Run("PageNegative", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?"+configs.PAGE_PARAM+"=-1", nil)
		_, err := GetPager(req)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNegativaParameterInUrl))
	})

	t.Run("ItemsPerPageNegative", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?" + configs.ITEMS_PER_PAGE_PARAM + "=-1", nil)
		_, err := GetPager(req)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNegativaParameterInUrl))
		assert.Contains(t, err.Error(), configs.ITEMS_PER_PAGE_PARAM)
	})

	t.Run("ValidValues", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?page=2&"+configs.ITEMS_PER_PAGE_PARAM+"=10", nil)
		pager, err := GetPager(req)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), pager.CurrentPage)
		assert.Equal(t, 10, pager.PageSize)
	})

	t.Run("PageProvidedItemsPerPageMissing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?page=3", nil)
		pager, err := GetPager(req)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), pager.CurrentPage)
		assert.Equal(t, 0, pager.PageSize)
	})

	t.Run("ItemsPerPageProvidedPageMissing", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?"+configs.ITEMS_PER_PAGE_PARAM+"=5", nil)
		pager, err := GetPager(req)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), pager.CurrentPage)
		assert.Equal(t, 5, pager.PageSize)
	})
}
