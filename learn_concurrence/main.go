package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ErrExit = errors.New("exit sign")
var ErrExitSign = errors.New("error sign")
var ErrOtherExitSign = errors.New("a one chan error sign")
var ErrHttpListen = errors.New("http listen error")

func main() {
	server := func(ctx context.Context, gc context.CancelFunc, srv *http.Server) error {

		g, ctx := errgroup.WithContext(ctx)

		// 启动服务
		g.Go(func() error {
			err := StartHttpServer(ctx, srv)
			if err != nil {
				gc()
			}
			return err
		})

		// 关闭服务
		g.Go(func() error {
			err := CloseHttpServer(ctx, srv)
			if err != nil {
				gc()
			}
			return err
		})

		// 监听信号
		g.Go(func() error {
			err := SignHandle(ctx)
			if err != nil {
				gc()
			}
			return err
		})

		if err := g.Wait(); err != nil {
			return err
		}
		return nil
	}

	srv := &http.Server{Addr: ":8080"}
	ctx, gc := context.WithCancel(context.Background())
	err := server(ctx, gc, srv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println("over")
}

func SignHandle(ctx context.Context) error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	for {

		select {
		case s := <-c:
			switch s {
			case syscall.SIGTERM:
				//TODO:服务器关闭逻辑处理
				return errors.Wrapf(ErrExit, "程序退出")
			default:
				return errors.Wrapf(ErrExitSign, "错误的退出信号")
			}
		case <-ctx.Done():
			//TODO:服务器关闭逻辑处理
			return errors.Wrapf(ErrOtherExitSign, "来自其他程序通知退出")
		}

	}
}

func StartHttpServer(_ context.Context, srv *http.Server) error {
	if err := srv.ListenAndServe(); err != nil {
		return errors.Wrapf(ErrHttpListen, fmt.Sprintf("http server 启动错误 error: %v", err))
	}
	return nil

	// 模拟错误
	//ticker := time.NewTicker(time.Second * 5)
	//select {
	//case <-ticker.C:
	//	return errors.Wrapf(ErrExit, "http server listen 错误")
	//}
	//return nil
}

func CloseHttpServer(ctx context.Context, srv *http.Server) error {
	ticker := time.NewTicker(time.Second * 10)
	select {
	case <-ticker.C:
		srv.Shutdown(ctx)
		return errors.Wrapf(ErrExit, "http server 关闭")
	case <-ctx.Done():
		srv.Shutdown(ctx)
		return errors.Wrapf(ErrOtherExitSign, "来自其他程序通知 http server 关闭")
	}
	return nil
}
