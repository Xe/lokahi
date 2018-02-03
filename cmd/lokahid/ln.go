package main

import (
	"context"
	"time"

	"github.com/Xe/ln"
	"github.com/twitchtv/twirp"
)

func makeLnHooks() *twirp.ServerHooks {
	hooks := &twirp.ServerHooks{}

	hooks.RequestRouted = func(ctx context.Context) (context.Context, error) {
		ctx = withStartTime(ctx)

		method, ok := twirp.MethodName(ctx)
		if !ok {
			return ctx, nil
		}

		pkg, ok := twirp.PackageName(ctx)
		if !ok {
			return ctx, nil
		}

		svc, ok := twirp.ServiceName(ctx)
		if !ok {
			return ctx, nil
		}

		ctx = ln.WithF(ctx, ln.F{
			"twirp_method":  method,
			"twirp_package": pkg,
			"twirp_service": svc,
		})

		return ctx, nil
	}

	hooks.ResponseSent = func(ctx context.Context) {
		f := ln.F{}
		now := time.Now()
		t, ok := getStartTime(ctx)
		if ok {
			f["response_time"] = now.Sub(t)
		}

		ln.Log(ctx, f, ln.Action("response sent"))
	}

	hooks.Error = func(ctx context.Context, e twirp.Error) context.Context {
		f := ln.F{}

		for k, v := range e.MetaMap() {
			f["twirp_meta_"+k] = v
		}

		ln.Error(ctx, e, f, ln.Action("twirp error"), ln.F{
			"twirp_error_code": e.Code(),
			"twirp_error_msg":  e.Msg(),
		})

		return ctx
	}

	return hooks
}

type ctxKey int

const (
	startTimeKey ctxKey = iota
)

func withStartTime(ctx context.Context) context.Context {
	return context.WithValue(ctx, startTimeKey, time.Now())
}

func getStartTime(ctx context.Context) (time.Time, bool) {
	t, ok := ctx.Value(startTimeKey).(time.Time)
	if !ok {
		return time.Time{}, false
	}

	return t, true
}
