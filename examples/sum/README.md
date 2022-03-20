## iOS integration

Set same value of `LDFLAGS` from cgo to Xcode
`OTHER_LDFLAGS` (Build Settings > Linking > Other Linker Flags).

## Setup gomobile

```shell
$ go install golang.org/x/mobile/cmd/gomobile
$ go install golang.org/x/mobile/cmd/gobind
$ asdf reshim golang
$ gomobile init
```
