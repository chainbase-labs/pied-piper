package geth

import (
	"testing"

	"github.com/ddl-hust/pied-piper/common/config"
	"github.com/stretchr/testify/require"
)

// test kafka admin init once.
func TestSendTx(t *testing.T) {
	t.Run("test on send tx", func(t *testing.T) {

		cfg := config.GetConf()

		cli, err := NewGethClient(cfg)
		require.Nil(t, err)

		err = cli.SendTx(cfg, "0xa43fe16908251ee70ef74718545e4fe6c5ccec9f")
		require.Nil(t, err)

	})
}
