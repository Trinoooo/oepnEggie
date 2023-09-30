package impl

import (
	"errors"
	"fmt"
	"github.com/Trinoooo/oepnEggie/types"
	"os"
	"testing"
)

var (
	handlerA = func(req types.HttpRequester, resp types.HttpResponser) error { return nil }
	handlerB = func(req types.HttpRequester, resp types.HttpResponser) error { return fmt.Errorf("handlerB") }
	handlerC = func(req types.HttpRequester, resp types.HttpResponser) error { return errors.New("handlerC") }
	handlerD = func(req types.HttpRequester, resp types.HttpResponser) error { return os.ErrClosed }
	handlerE = func(req types.HttpRequester, resp types.HttpResponser) error { return os.ErrDeadlineExceeded }
	handlerF = func(req types.HttpRequester, resp types.HttpResponser) error { return os.ErrExist }
	handlerG = func(req types.HttpRequester, resp types.HttpResponser) error { return os.ErrInvalid }

	handlerH = func(req types.HttpRequester, resp types.HttpResponser) error { return nil }
	handlerI = func(req types.HttpRequester, resp types.HttpResponser) error { return fmt.Errorf("handlerB") }
	handlerJ = func(req types.HttpRequester, resp types.HttpResponser) error { return errors.New("handlerC") }
	handlerK = func(req types.HttpRequester, resp types.HttpResponser) error { return os.ErrClosed }
	handlerL = func(req types.HttpRequester, resp types.HttpResponser) error { return os.ErrDeadlineExceeded }
	handlerM = func(req types.HttpRequester, resp types.HttpResponser) error { return os.ErrExist }
	handlerN = func(req types.HttpRequester, resp types.HttpResponser) error { return os.ErrInvalid }
)
var registerTableAlpha = map[string]Handler{
	"/a": handlerA,
	"/b": handlerB,
	"/c": handlerC,
	"/d": handlerD,
	"/e": handlerE,
	"/f": handlerF,
	"/g": handlerG,
}

var registerTableBeta = map[string]Handler{
	"/h": handlerH,
	"/i": handlerI,
	"/j": handlerJ,
	"/k": handlerK,
	"/l": handlerL,
	"/m": handlerM,
	"/n": handlerN,
}

func TestDefaultMux_RegisterDuplicate(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	t.Run("alpha", func(t *testing.T) {
		for path, handler := range registerTableAlpha {
			if err := server.Get(path, handler); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("beta", func(t *testing.T) {
		for path, handler := range registerTableAlpha {
			if err := server.Get(path, handler); err != nil {
				t.Error(err)
			}
		}
	})
}

func TestDefaultMux_RegisterRace(t *testing.T) {
	t.Run("alpha", func(t *testing.T) {
		for path, handler := range registerTableAlpha {
			if err := server.Get(path, handler); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("beta", func(t *testing.T) {
		for path, handler := range registerTableBeta {
			if err := server.Get(path, handler); err != nil {
				t.Error(err)
			}
		}
	})
}

func TestDefaultMux_Close(t *testing.T) {
	t.Log("do something...")
}
