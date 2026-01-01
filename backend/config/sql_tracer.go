package config

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type sqlTraceStartKey struct{}
type sqlTraceSQLKey struct{}

type SQLTracer struct{}

func (t *SQLTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx = context.WithValue(ctx, sqlTraceStartKey{}, time.Now())
	ctx = context.WithValue(ctx, sqlTraceSQLKey{}, data.SQL)
	return ctx
}

func (t *SQLTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	start, _ := ctx.Value(sqlTraceStartKey{}).(time.Time)
	latency := time.Since(start)

	sql, _ := ctx.Value(sqlTraceSQLKey{}).(string)
	sql = strings.TrimSpace(sql)
	sql = strings.Join(strings.Fields(sql), " ") // single-line for terminal readability

	if data.Err != nil {
		log.Printf("SQL | ERROR | %s | %s | %v", latency, sql, data.Err)
		return
	}

	log.Printf("SQL | OK    | %s | %s", latency, sql)
}
