/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"html"
	"strings"

	"github.com/ta2gch/iris/runtime"
	"github.com/ta2gch/iris/runtime/ilos/instance"
	"github.com/ta2gch/jquery"
)

type Dom struct{}

func (Dom) Write(p []byte) (n int, err error) {
	jQuery("#output").Append(string(p))
	return len(p), nil
}

func (dom Dom) Read(p []byte) (n int, err error) {
	input := strings.Replace(jQuery("#input").Html(), "<br>", " ", -1)
	jQuery("#input").SetHtml("")
	copy(p, input)
	dom.Write([]byte(input + "\n"))
	return len(p), nil
}

var jQuery = jquery.NewJQuery

func main() {
	dom := new(Dom)
	runtime.TopLevel.StandardInput = instance.NewStream(dom, nil)
	runtime.TopLevel.StandardOutput = instance.NewStream(nil, dom)
	runtime.TopLevel.ErrorOutput = instance.NewStream(nil, dom)
	runtime.TopLevel.Function.Set(instance.NewSymbol("READ"), nil)
	runtime.TopLevel.Function.Set(instance.NewSymbol("READ-LINE"), nil)
	runtime.TopLevel.Function.Set(instance.NewSymbol("READ-CHAR"), nil)
	fmt.Fprintf(dom, `Welcome to Iris %v. Iris is an ISLisp implementation on Go.
This REPL works with gopherjs and has no methods to get input.

Copyright &copy; 2017 TANIGUCHI Masaya All Rights Reserved.`, runtime.Version)
	jQuery("body").On(jquery.KEYDOWN, func(e jquery.Event) {
		if !e.ShiftKey && e.KeyCode == 13 {
			fmt.Fprint(dom, "\n> ")
			exp, err := runtime.Read(runtime.TopLevel)
			if err != nil {
				fmt.Fprint(dom, html.EscapeString(err.String()))
				return
			}
			ret, err := runtime.Eval(runtime.TopLevel, exp)
			if err != nil {
				fmt.Fprint(dom, html.EscapeString(err.String()))
				return
			}
			fmt.Fprint(dom, html.EscapeString(ret.String()))
			return
		}
	})
}
