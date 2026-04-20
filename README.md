# Blog

## Ready

*`go 1.26+`*

```bash
go install github.com/google/wire/cmd/wire@latest
```

## Install

```bash
git clone https://github.com/wantnotshould/blog
cd blog/internal
wire ./wire
cd ..
go mod tidy
go build
```

## 🙏 Acknowledgements

- Inspired by the [alist](https://github.com/AlistGo/alist) project.

## ⚖️ License

MIT License. See [LICENSE](./LICENSE).
