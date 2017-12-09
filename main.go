/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"html"

	"github.com/gopherjs/jquery"
	"github.com/asciian/iris/runtime"
	"github.com/asciian/iris/runtime/ilos"
	"github.com/asciian/iris/runtime/ilos/class"
	"github.com/asciian/iris/runtime/ilos/instance"
)

const version = "526f28d"

var stream = make(chan string)

type Dom struct{}

func (Dom) Write(p []byte) (n int, err error) {
	jQuery("#output").Append(string(p))
	return len(p), nil
}

func (dom Dom) Read(p []byte) (n int, err error) {
	s := <-stream
	fmt.Fprint(dom, s, "\n")
	copy(p, []byte(s))
	return len(p), nil
}

var jQuery = jquery.NewJQuery

func main() {
	prompt := true
	dom := new(Dom)
	runtime.TopLevel.StandardInput = instance.NewStream(dom, nil)
	runtime.TopLevel.StandardOutput = instance.NewStream(nil, dom)
	runtime.TopLevel.ErrorOutput = instance.NewStream(nil, dom)
	fmt.Fprintf(dom, `Welcome to Iris (%v). Iris is an ISLisp implementation on Go.
This REPL works on JavaScript with gopherjs.

Copyright &copy; 2017 asciian All Rights Reserved.`, version)
	jQuery("#version").SetHtml(version)
	jQuery("#input").On(jquery.KEYDOWN, func(e jquery.Event) {
		if !e.ShiftKey && e.KeyCode == 13 {
			if prompt {
				fmt.Fprint(dom, "\n> ")
			}
			go func() {
				input := jQuery("#input").Html()
				jQuery("#input").SetHtml("")
				stream <- jQuery(`<span>` + input + `</span>`).Text()
			}()
		}
	})
	for {
		prompt = true
		jQuery("#prompt").Show()
		runtime.TopLevel.StandardInput = instance.NewStream(dom, nil)
		exp, err := runtime.Read(runtime.TopLevel)
		if err != nil {
			if !ilos.InstanceOf(class.EndOfStream, err) {
				fmt.Fprint(dom, html.EscapeString(err.String()))
			}
			continue
		}
		prompt = false
		jQuery("#prompt").Hide()
		ret, err := runtime.Eval(runtime.TopLevel, exp)
		if err != nil {
			fmt.Fprint(dom, html.EscapeString(err.String()))
			continue
		}
		fmt.Fprint(dom, html.EscapeString(ret.String()))
	}
}
