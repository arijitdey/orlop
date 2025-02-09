// Copyright (c) 2020 Ketch Kloud, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package orlop

import (
	"context"
	"go.ketch.com/lib/orlop/v2/config"
	"go.ketch.com/lib/orlop/v2/env"
	"go.ketch.com/lib/orlop/v2/log"
	"go.ketch.com/lib/orlop/v2/logging"
	"go.ketch.com/lib/orlop/v2/service"
	"go.uber.org/fx"
)

func FxOptions(o config.Config) fx.Option {
	return o.Options()
}

func FxContext(ctx context.Context) fx.Option {
	return fx.Provide(func() context.Context { return ctx })
}

func Populate(ctx context.Context, prefix string, e env.Environment, module fx.Option, cfg config.Config, targets ...interface{}) error {
	e.Load()

	if err := Unmarshal(prefix, cfg); err != nil {
		return err
	}

	app := fx.New(
		logging.WithLogger(log.New()),
		FxContext(ctx),
		FxOptions(cfg),
		fx.Supply(service.Name(prefix)),
		Module,
		module,
		fx.Populate(targets...),
	)

	if err := app.Err(); err != nil {
		return err
	}

	if err := app.Start(ctx); err != nil {
		return err
	}

	if err := app.Stop(ctx); err != nil {
		return err
	}

	return app.Err()
}
