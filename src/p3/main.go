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
	errgroup.WithContext(context.Background())
}

type Server interface {
	Run() error
	Stop() error
}

type Option struct {
	timeout time.Duration
	ss      []*Server
	extra   []func() error
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) WithTimeout(duration time.Duration) {
	o.timeout = duration
}

func (o *Option) AddDefaultSignalHandler() {
	o.extra = append(o.extra, defaultSignalHandle)
}

func (o *Option) Serve() error {
	g, ctx := errgroup.WithContext(context.Background())
	for _, s := range o.ss {
		g.Go(func() error {
			<-ctx.Done()
			return (*s).Stop()
		})
		g.Go((*s).Run)
	}
	for _, x := range o.extra {
		g.Go(x)
	}
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
	return g.Wait()
}

func defaultSignalHandle() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	return fmt.Errorf("receive signal %s", <-ch)
}
