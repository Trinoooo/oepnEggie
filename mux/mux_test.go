package mux

import (
	"errors"
	"fmt"
	"github.com/Trinoooo/oepnEggie"
	"github.com/Trinoooo/oepnEggie/mistake"
	"github.com/Trinoooo/oepnEggie/ptlpsr"
	"os"
	"testing"
)

var (
	handlerA = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return nil }
	handlerB = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return fmt.Errorf("handlerB") }
	handlerC = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return errors.New("handlerC") }
	handlerD = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return os.ErrClosed }
	handlerE = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return os.ErrDeadlineExceeded }
	handlerF = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return os.ErrExist }
	handlerG = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return os.ErrInvalid }

	handlerH = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return nil }
	handlerI = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return fmt.Errorf("handlerB") }
	handlerJ = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return errors.New("handlerC") }
	handlerK = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return os.ErrClosed }
	handlerL = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return os.ErrDeadlineExceeded }
	handlerM = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return os.ErrExist }
	handlerN = func(req ptlpsr.HttpRequester, resp ptlpsr.HttpResponser) error { return os.ErrInvalid }
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

var server *oepnEggie.Server

func TestMain(m *testing.M) {
	var err error
	server, err = oepnEggie.NewServerBuilder().Build()
	if err != nil {
		panic(err)
	}

	m.Run()
}

func TestDefaultMux_RegisterDuplicate(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	t.Run("alpha", func(t *testing.T) {
		for path, handler := range registerTableAlpha {
			if err := server.Get(path, handler); err != nil && err != mistake.NewInvalid() {

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
