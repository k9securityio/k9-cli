/*
Copyright © 2022 The K9CLI Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import "fmt"

type IllegalArgumentError struct {
	argument string
	error    string
}

func (e IllegalArgumentError) Error() string {
	return e.argument + `, ` + e.error
}
func (e IllegalArgumentError) Arg() string {
	return e.argument
}

type AggregateError struct {
	isPartial bool
	errors    []error
}

func (e AggregateError) IsPartial() bool {
	return e.isPartial
}
func (e AggregateError) Errors() []error {
	return e.errors
}
func (e *AggregateError) Error() string {
	var o string
	for i, v := range e.errors {
		if i > 0 {
			o += "\n"
		}
		o += fmt.Sprintf("%v", v)
	}
	return o
}
