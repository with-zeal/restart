package restart

import (
	"net"
	"net/http"
	"os"
	"syscall"
)

type restartServer struct {
	server   *http.Server
	listener net.Listener
}

func NewServer(router http.Handler, addr string) *restartServer {
	return &restartServer{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (ser *restartServer) Run() (err error) {
	if os.Getenv("graceful") == "true" {
		f := os.NewFile(3, "")
		ser.listener, err = net.FileListener(f)
		syscall.Kill(os.Getppid(), syscall.SIGQUIT)
	} else {
		ser.listener, err = net.Listen("tcp", ser.server.Addr)
	}
	if err != nil {
		logger.Println(err)
		return
	}

	go func() {
		if len(beforeTask) != 0 {
			for _, task := range beforeTask {
				if task.isGo {
					go task.task()
				} else {
					task.task()
				}
			}
		}

		err = ser.server.Serve(ser.listener)
		if err != nil && err != http.ErrServerClosed {
			logger.Println(err)
			panic(err)
		}
	}()

	handlerSignal(ser.listener, ser.server, afterTask)
	return
}

func (ser *restartServer) RunTLS(certPath, keyPath string) (err error) {
	if os.Getenv("graceful") == "true" {
		f := os.NewFile(3, "")
		ser.listener, err = net.FileListener(f)
		syscall.Kill(os.Getppid(), syscall.SIGQUIT)
	} else {
		ser.listener, err = net.Listen("tcp", ser.server.Addr)
	}
	if err != nil {
		logger.Println(err)
		return
	}

	go func() {
		if len(beforeTask) != 0 {
			for _, task := range beforeTask {
				if task.isGo {
					go task.task()
				} else {
					task.task()
				}
			}
		}

		err = ser.server.ServeTLS(ser.listener, certPath, keyPath)
		if err != nil && err != http.ErrServerClosed {
			logger.Println(err)
			panic(err)
			return
		}
	}()

	handlerSignal(ser.listener, ser.server, afterTask)
	return
}
