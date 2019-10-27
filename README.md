![Travis CI build status](https://travis-ci.com/chaimleib/errors.svg?branch=master)

# github.com/chaimleib/errors

This Go package is a drop-in replacement to the built-in [`errors`](https://golang.org/pkg/errors/) package. It is designed for Go 1.13, while under 1.12, the `Is`, `As` and `Unwrap` functions have been backported.

## Why?

**Error messages exist to help you fix the problem.** I kept finding that regular error messages required too much searching to understand the issue. I like this better:

```
main.main() main.go:34 error getting user profile
github.com/chaimleib/client.UserProfile(ctx, "alice") userprofile.go:124 error authenticating
github.com/chaimleib/client.Authenticate(ctx, "bob", pw) client.go:63 password expired
```

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

## What else can I do?

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

## License

Copyright 2019 Chaim Leib Halbert

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this source code except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
