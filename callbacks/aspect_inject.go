/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package callbacks

import (
	"context"

	"github.com/cloudwego/eino/internal/callbacks"
	"github.com/cloudwego/eino/schema"
)

// OnStart Fast inject callback input / output aspect for component developer
// e.g.
//
//	func (t *testChatModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (resp *schema.Message, err error) {
//		defer func() {
//			if err != nil {
//				callbacks.OnEnd(ctx, err)
//			}
//		}()
//
//		ctx = callbacks.OnStart(ctx, &model.CallbackInput{
//			Messages: input,
//			Tools:    nil,
//			Extra:    nil,
//		})
//
//		// do smt
//
//		ctx = callbacks.OnEnd(ctx, &model.CallbackOutput{
//			Message: resp,
//			Extra:   nil,
//		})
//
//		return resp, nil
//	}
//

// OnStart invokes the OnStart logic for the particular context, ensuring that all registered
// handlers are executed in reverse order (compared to add order) when a process begins.
func OnStart[T any](ctx context.Context, input T) context.Context {
	ctx, _ = callbacks.On(ctx, input, callbacks.OnStartHandle[T], TimingOnStart)

	return ctx
}

// OnEnd invokes the OnEnd logic of the particular context, allowing for proper cleanup
// and finalization when a process ends.
// handlers are executed in normal order (compared to add order).
func OnEnd[T any](ctx context.Context, output T) context.Context {
	ctx, _ = callbacks.On(ctx, output, callbacks.OnEndHandle[T], TimingOnEnd)

	return ctx
}

// OnStartWithStreamInput invokes the OnStartWithStreamInput logic of the particular context, ensuring that
// every input stream should be closed properly in handler.
// handlers are executed in reverse order (compared to add order).
func OnStartWithStreamInput[T any](ctx context.Context, input *schema.StreamReader[T]) (
	nextCtx context.Context, newStreamReader *schema.StreamReader[T]) {

	return callbacks.On(ctx, input, callbacks.OnStartWithStreamInputHandle[T], TimingOnStartWithStreamInput)
}

// OnEndWithStreamOutput invokes the OnEndWithStreamOutput logic of the particular, ensuring that
// every input stream should be closed properly in handler.
// handlers are executed in normal order (compared to add order).
func OnEndWithStreamOutput[T any](ctx context.Context, output *schema.StreamReader[T]) (
	nextCtx context.Context, newStreamReader *schema.StreamReader[T]) {

	return callbacks.On(ctx, output, callbacks.OnEndWithStreamOutputHandle[T], TimingOnEndWithStreamOutput)
}

// OnError invokes the OnError logic of the particular, notice that error in stream will not represent here.
// handlers are executed in normal order (compared to add order).
func OnError(ctx context.Context, err error) context.Context {
	ctx, _ = callbacks.On(ctx, err, callbacks.OnErrorHandle, TimingOnError)

	return ctx
}

// ReuseHandlers initializes a new context with the provided RunInfo, while using the same handlers already exist.
func ReuseHandlers(ctx context.Context, info *RunInfo) context.Context {
	return callbacks.ReuseHandlers(ctx, info)
}

// InitCallbacks initializes a new context with the provided RunInfo and handlers.
// Any previously set RunInfo and Handlers for this ctx will be overwritten.
func InitCallbacks(ctx context.Context, info *RunInfo, handlers ...Handler) context.Context {
	return callbacks.InitCallbacks(ctx, info, handlers...)
}
