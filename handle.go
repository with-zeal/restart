package restart

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var TimeOut = time.Minute

func reload(listener net.Listener) {
	tl, ok := listener.(*net.TCPListener)
	if !ok {
		logger.Println(errors.New("listener is not tcp listener"))
		return
	}

	f, err := tl.File()
	if err != nil {
		logger.Println(err)
		return
	}

	os.Setenv("graceful", "true")
	for _, env := range envs {
		os.Setenv(env.key, env.value)
	}
	cmd := exec.Command(os.Args[0])
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// put socket FD at the first entry
	cmd.ExtraFiles = []*os.File{f}
	err = cmd.Start()
	if err != nil {
		logger.Println(err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		logger.Println(err)
	}
}

func serverQuit(server *http.Server) {
	server.SetKeepAlivesEnabled(false)
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Println(err)
	}
}

func handlerSignal(listener net.Listener, server *http.Server, afterTask []func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

wait:
	sig := <-quit

	switch sig {
	case syscall.SIGHUP:
		goto wait
	case syscall.SIGINT:
		for _, task := range afterTask {
			task()
		}
		go reload(listener)
		goto wait
	}
	serverQuit(server)
	return
}
