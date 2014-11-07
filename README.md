# cachey

Expire-able cache

## Usage

```go
c := cachey.NewCache()
c.Set("foo", "bar", 3)
if v, ok := c.Get("foo"); ok {
    fmt.Println(v)
} else {
	fmt.Println("Not Found")
}

time.Sleep(3 * time.Second)

if v, ok := c.Get("foo"); ok {
    fmt.Println(v)
} else {
	fmt.Println("Not Found")
}
```

```
bar
Not Found
```

## Installation

```
$ go get github.com/mattn/cachey
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a mattn)
