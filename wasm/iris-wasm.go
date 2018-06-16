/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"fmt"
	"html"
	"strings"
	"syscall/js"

	"github.com/asciian/iris/runtime"
	"github.com/asciian/iris/runtime/ilos/instance"
)

func eval(s string) string {
	r := strings.NewReader(s)
	w := new(bytes.Buffer)
	runtime.TopLevel.StandardInput = instance.NewStream(r, nil)
	runtime.TopLevel.StandardOutput = instance.NewStream(nil, w)
	runtime.TopLevel.ErrorOutput = instance.NewStream(nil, w)
	runtime.TopLevel.Function.Set(instance.NewSymbol("READ"), nil)
	runtime.TopLevel.Function.Set(instance.NewSymbol("READ-LINE"), nil)
	runtime.TopLevel.Function.Set(instance.NewSymbol("READ-CHAR"), nil)
	p, err := runtime.Read(runtime.TopLevel)
	if err != nil {
		fmt.Fprint(w, html.EscapeString(err.String()))
		return w.String()
	}
	e, err := runtime.Eval(runtime.TopLevel, p)
	if err != nil {
		fmt.Fprint(w, html.EscapeString(err.String()))
		return w.String()
	}
	fmt.Fprint(w, html.EscapeString(e.String()))
	return w.String()
}

func main() {
	input := js.Global.Get("document").Call("getElementById", "input")
	button := js.Global.Get("document").Call("getElementById", "button")
	output := js.Global.Get("document").Call("getElementById", "output")
	button.Call("addEventListener", "click", js.NewCallback(func(args []js.Value) {
		output.Set("value", eval(input.Get("value").String()))
	}))
	select {}
}
