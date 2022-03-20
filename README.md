# go-wasm3

Golang wrapper for WASM3 (partial)

## Setup (macOS)

Install `docker` & `asdf`

```shell
$ asdf plugin add golang https://github.com/kennyp/asdf-golang.git
$ asdf plugin add cmake  https://github.com/asdf-community/asdf-cmake.git
$ asdf plugin add tinygo https://github.com/schmir/asdf-tinygo.git  
```

## Related projects

- Strongly inspired by [go-wasm3][1], even though I wrote this binding from
  scratch
- If you are looking for a webassembly runtime that can be used in go, check out
  [wazero][2] instead, which does not depend on any CGO library, making it much
  more relevant for use with gomobile

<!-- links -->

[1]: https://github.com/matiasinsaurralde/go-wasm3
[2]: https://github.com/tetratelabs/wazero
