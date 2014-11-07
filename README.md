# cachey

Expire-able cache

## Usage

```go
c := cachey.NewCache()
c.Set("foo", "bar")
if v, ok := c.Get("foo"); ok {
    fmt.Println(v)
}
c.Delete("foo")
```

## Installation

```
$ go get github.com/mattn/cachey
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a mattn)
