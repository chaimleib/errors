![Travis CI build status](https://travis-ci.com/chaimleib/errors.svg?branch=master)

# github.com/chaimleib/errors

https://godoc.org/github.com/chaimleib/errors

![Doctor Gopher treating patient](doctorPatientGophers.png)

This Go package is a drop-in replacement to the built-in [`errors`](https://golang.org/pkg/errors/) package. It is designed for Go 1.13, while under 1.12, the `Is`, `As` and `Unwrap` functions have been backported.

## Why?

**Error messages exist to help you fix the problem.** I kept finding that regular error messages required too much searching to understand the issue. I like this better:

```
main.main() main.go:34 error getting user profile
~/client.UserProfile(ctx, "alice") userprofile.go:124 error authenticating
~/client.Authenticate(ctx, "bob", pw) client.go:63 password expired
```

> Your `go.mod` main module name is abbreviated as `~`.

That's much more helpful than this:

```
password expired
```

What operation required a password? Where did this get generated? How did we get there?

With enhanced errors, you can begin fixing the problem right away:

1. You know exactly where the error came from.
2. You know what code path led to the error.
3. You know what inputs each function was seeing.

All this, without having to use [`rg`](https://github.com/BurntSushi/ripgrep) or open up any files!

## How?

1. At the beginning of your methods, have this:

```go
func (client Client) Authenticate(ctx context.Context, user, pw string) (UserProfile, error) {
  b := errors.NewBuilder("ctx, %q, pw", user)
  // Note: You don't have to show values for all the arguments.
  // An idea: for brevity, try "[%d]aSlice", len(aSlice) instead of showing the
  // whole slice or map, etc.
```

2. Whenever you return an error:

```go
if err != nil {
  return nil, b.Wrap(err, "error parsing %q", value)
}
```

3. At the top of the program, print the full stack trace:

```go
if err != nil {
  fmt.Println(errors.StackString(err))
}
```

## Another example ([try me!](https://goplay.space/#HE4BuAJaZYA)):

```go
func main() {
	b := errors.NewBuilder("")
	if err := FileHasHello("greet.txt"); err != nil {
		err = b.Wrap(err, "program failed")
		fmt.Println(errors.StackString(err))
		return
	}
}

func FileHasHello(fpath string) error {
	b := errors.NewBuilder("%q", fpath)
	buf, err := ioutil.ReadFile(fpath)
	if err != nil {
		return b.Wrap(err, "could not open file")
	}
	if !bytes.Contains(buf, []byte("hello")) {
		return b.Errorf("could not find `hello` in file")
	}
	return nil
}

/* output:
main.main() prog.go:14 program failed
main.FileHasHello("greet.txt") prog.go:24 could not open file
open greet.txt: No such file or directory
No such file or directory
*/
```

## What else can I do?

* Print just one level of stack trace using `StackStringAt(err)`. Index into `Stack(err)` to select an error by its `Unwrap()` depth.

* Use `Is`, `As` and `Unwrap` in Go 1.12 (added officially in Go 1.13)

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

* Get call site info with `NewFuncInfo(calldepth)`. (This is for debugging output only. It is bad design to write application logic around these values.)

```go
fi := errors.NewFuncInfo(1)
fmt.Println(fi.File(), fi.Line(), fi.FuncName())
// /absolute/path/to/main.go 12 main.main
```

* Customize your stack trace formatting by writing your our `StackString()` function. All the necessary plumbing (like `FuncInfo()`, `ArgStringer()` and `Wrapper()`) is exposed as public methods.

## License

Copyright 2019-2020 Chaim Leib Halbert

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this source code except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
