package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHashPassword(t *testing.T) {

	data := []struct {
		input string
	}{
		{"123456xz"},
		{"OcSOglmzaLKUCWn9mvkSQwJPdC"},
		{"SPf59uD5SpQjnOU8mng0EKbpO9"},
	}
	for _, item := range data {
		t.Run("Test Generate Hash Password", func(t *testing.T) {
			result, err := GenerateHashedPassword(item.input)
			assert.NoError(t, err)

			err = CompareHashAndPassword(item.input, result)
			assert.NoError(t, err)
		})
	}
}
