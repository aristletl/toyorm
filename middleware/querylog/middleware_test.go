package querylog

import (
	"context"
	"fmt"
	"github.com/aristletl/toyorm"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{}
	builder.LogFunc(func(sql string, args ...any) {
		fmt.Println(sql)
	})
	db, err := toyorm.Open("sqlite3", "file:test.db?cache=shared&mode=memory", toyorm.DBWithMiddlewares(builder.Build()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = toyorm.NewSelector[TestModel](db).Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

type TestModel struct {
}
