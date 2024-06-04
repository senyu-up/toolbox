package addr

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestExternalIP(t *testing.T) {
	t.Run("测试获取ip功能", func(t *testing.T) {
		ip, err := ExternalIP()
		assert.Equal(t, err, nil)
		split := strings.Split(ip, ".")
		assert.Equal(t, len(split), 4)
	})

}
