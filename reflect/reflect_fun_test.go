package reflect

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestReflectFun(t *testing.T) {
	err := HandleFun(handleImpl, context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestReflectFunMulti(t *testing.T) {
	err := HandleFun(handleMultiImpl, context.Background(), "abc", 12)
	if err != nil {
		t.Error(err)
	}
}

func handleImpl(ctx context.Context)  (int, error) {
	return 3, nil
}

func handleWithErrImpl(ctx context.Context)  (int, error) {
	return 3, errors.New("qq")
}

func handleMultiImpl(ctx context.Context, str string, num int) error {
	fmt.Println("handle data:", str, num)
	return nil
}
