package idempotency

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type RedisStore struct {
	addr string
}

func NewRedisStore(addr string) *RedisStore {
	return &RedisStore{addr: addr}
}

func (s *RedisStore) Reserve(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", s.addr)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	} else {
		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
	}

	command := respArray("SET", key, "1", "NX", "PX", fmt.Sprintf("%d", ttl.Milliseconds()))
	if _, err := conn.Write([]byte(command)); err != nil {
		return false, err
	}

	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return false, err
	}

	line = strings.TrimSpace(line)
	switch {
	case line == "+OK":
		return true, nil
	case line == "$-1":
		return false, nil
	case strings.HasPrefix(line, "-"):
		return false, errors.New(strings.TrimPrefix(line, "-"))
	default:
		return false, fmt.Errorf("unexpected redis response: %s", line)
	}
}

func respArray(values ...string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("*%d\r\n", len(values)))
	for _, value := range values {
		builder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
	}
	return builder.String()
}
