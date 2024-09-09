package tests

import (
	"avito/assembly"
	"avito/config"
	"avito/db"
	"avito/log"
	"context"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"testing"
)

type Test struct {
	Assertions *require.Assertions
	Server     *httptest.Server
	DbCli      db.DB
	TestID     uint32
	URL        string
}

func InitTest(t *testing.T) *Test {
	var (
		testID = rand.Uint32()
	)
	assert := require.New(t)
	ctx := context.Background()

	cfg := ConfigDefault().WithSchema("test_" + strconv.Itoa(int(testID)))
	_, err := cfg.Validate(ctx)
	assert.NoError(err)

	logger := log.NewLogger(cfg.AppMode)
	assembler := assembly.NewAssembler(logger)

	router, err := assembler.Assemble(ctx, cfg)
	assert.NoError(err)

	srv := httptest.NewServer(router)
	dbCli, err := db.Open(ctx, cfg.Dsn())
	assert.NoError(err)

	err = dbCli.CreateSchema(cfg.DbSchema)
	assert.NoError(err)

	t.Cleanup(func() {
		srv.Close()
		err := dbCli.Close()
		assert.NoError(err)
	})

	return &Test{
		Assertions: assert,
		Server:     srv,
		DbCli:      dbCli,
		TestID:     rand.Uint32(),
	}
}

func ConfigDefault() *config.Config {
	return &config.Config{
		DbUsername: "postgres",
		DbPassword: "postgres",
		DbHost:     "localhost",
		DbPort:     "5432",
		DbName:     "avitotest",
		DbSSL:      "disable",
		AppMode:    "dev",
	}
}
