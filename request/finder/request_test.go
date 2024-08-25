package finder

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRequester(t *testing.T) {
	r := regexp.MustCompile(`[?\s]`)
	c1 := r.Split("apps i ./asd", 2)
	c2 := r.Split(`apps.install?["./asd"]`, 2)
	assert.Equal(t, 2, len(c1))
	assert.Equal(t, c1[0], c2[0])
}
