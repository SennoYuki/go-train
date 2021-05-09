package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	o := NewOption()
	o.WithTimeout(time.Second * 10)
	o.AddServer(new(EmptyServer), new(EmptyServer))
	o.ListenSignal(syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	if err := o.Serve(); err != nil {
		fmt.Println(err)
	}
}

type EmptyServer struct {
}

func (s EmptyServer) Run() error {
	return nil
}

func (s EmptyServer) Stop() error {
	return nil
}

type Server interface {
	Run() error
	Stop() error
}

type Option struct {
	ss      []Server
	signals []os.Signal
	timeout time.Duration
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) WithTimeout(duration time.Duration) {
	o.timeout = duration
}

func (o *Option) ListenSignal(signals ...os.Signal) {
	o.signals = signals
}

func (o *Option) AddServer(servers ...Server) {
	o.ss = append(o.ss, servers...)
}

func (o *Option) Serve() error {
	g, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	for _, s := range o.ss {
		g.Go(func() error {
			<-ctx.Done()
			return s.Stop()
		})
		g.Go(s.Run)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, o.signals...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ch:
				if o.timeout != 0 {
					g.Go(func() error {
						time.After(o.timeout)
						return fmt.Errorf("reach max quit wait time %v", o.timeout)
					})
				}
				cancel()
				return nil
			}
		}
	})
	return g.Wait()
}
