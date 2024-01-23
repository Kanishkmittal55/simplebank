package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	db "github.com/kanishkmittal55/simplebank/db/sqlc"
	"github.com/kanishkmittal55/simplebank/db/util"
	"github.com/kanishkmittal55/simplebank/worker"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskPayload := worker.NewRedisTaskDistributor(redisOpt)

	server, err := NewServer(config, store, taskPayload)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
