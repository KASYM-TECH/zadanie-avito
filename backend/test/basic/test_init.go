package basic

import (
	"avito/assembly"
	"avito/config"
	"avito/db"
	"avito/log"
	"context"
	"github.com/stretchr/testify/require"
	"github.com/txix-open/isp-kit/http/httpcli"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"testing"
)

type Test struct {
	Assertions *require.Assertions
	Server     *httptest.Server
	Cli        *httpcli.Client
	DbCli      db.DB
	TestId     uint32
	URL        string
}

func InitTest(t *testing.T) *Test {
	var (
		testId = rand.Uint32()
	)
	assert := require.New(t)
	ctx := context.Background()

	cfg := ConfigDefault().WithSchema("test_" + strconv.Itoa(int(testId)))
	_, err := cfg.Validate(ctx)
	assert.NoError(err)

	logger := log.NewLogger(cfg.AppMode)
	assembler := assembly.NewAssembler(logger)

	router, err := assembler.Assemble(ctx, cfg)
	assert.NoError(err)

	srv := httptest.NewServer(router)
	dbCli, err := db.Open(ctx, cfg.Dsn())
	assert.NoError(err)

	t.Cleanup(func() {
		_ = dbCli.DropSchema(cfg.DbSchema)
		srv.Close()
		err := dbCli.Close()
		assert.NoError(err)
	})

	return &Test{
		Assertions: assert,
		Server:     srv,
		DbCli:      dbCli,
		TestId:     rand.Uint32(),
		URL:        srv.URL,
		Cli:        httpcli.New(),
	}
}

func ConfigDefault() *config.Config {
	return config.LoadFromEnv("test.env")
}
