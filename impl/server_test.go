package impl

import (
	"github.com/Trinoooo/oepnEggie/logs"
	"github.com/Trinoooo/oepnEggie/types"
	"github.com/luci/go-render/render"
	"os"
	"os/signal"
	"sync"
	"testing"
)

func TestServerBuilder_Build(t *testing.T) {
	if testing.Short() {
		t.Log("TestServerBuilder_Build skip in short mode")
		t.SkipNow()
	}
	builder := NewServerBuilder()
	server, err := builder.WithProtocolParser(nil).WithMultiplexer(nil).Build()
	if err != nil {
		t.Error(err)
	}

	// 同步注册handler
	server.Get("/abc", func(req types.HttpRequester, resp types.HttpResponser) error {
		return nil
	})
	server.Post("/def", func(req types.HttpRequester, resp types.HttpResponser) error {
		return nil
	})

	t.Log(render.Render(server))
}

func TestServer_Serve(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Serve()
		if err != nil {
			logs.V1().Println(err)
		}
	}()

	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, os.Kill, os.Interrupt)
	select {
	case <-signalCh:
		err := server.Close()
		if err != nil {
			logs.V1().Println(err)
		}
	}

	wg.Wait()
	logs.V1().Println("gracefully exit")
}
