# github.com/chaimleib/errors

This Go package is a drop-in replacement to the built-in [`errors`](https://golang.org/pkg/errors/) package.

## What can I do?

Here's some code:

* Wrap errors ([try me](https://goplay.space/#6s48CkaR89-))

```go
timeoutErr := fmt.Errorf("timeout")
serverErr := errors.Wrap(timeoutErr, "server A failed")
loadErr := errors.Wrap(serverErr, "could not load resource")
fmt.Println("error:", loadErr)
// error: could not load resource
```

* Print traces unwrapping the error chain ([try me](https://goplay.space/#6s48CkaR89-))

```go
// ...
trace := errors.StackString(loadErr)
fmt.Println("trace:")
fmt.Println(trace)
// trace:
// could not load resource
// server A failed
// timeout
```

* Show the function call where errors happened ([try me](https://goplay.space/#jGQ4wzxt-NS))

```go
func ping(path string) error {
	em := errors.NewErrorMaker("%q", path)
	if !strings.HasPrefix(path, "/") {
		return em.Errorf("not an absolute path")
	}
	if err := requestMust200("https://example.com" + path); err != nil {
		return em.Wrap(err, "request failed")
	}
	return nil
}
// prog.go:21 main.ping("/health"): request failed
```

* Group errors ([try me](https://goplay.space/#auXQKNwP0VV))

```go
errs := []error{
  fmt.Errorf("server A failed"),
  fmt.Errorf("server B failed"),
}
wrapped := errors.Wrap(errors.Group(errs), "all servers failed")
fmt.Println(errors.StackString(wrapped))
// all servers failed
// [
//     server A failed
//     ,
//     server B failed
// ]
```

## License

Copyright 2019 Chaim Leib Halbert

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
