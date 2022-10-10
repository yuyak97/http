# http サーバー

## /server/main.go 使い方

- GET /

```
curl -i http://localhost:7777/
```

- GET /sample.html

```
curl -i http://localhost:7777/sample.html
```

- GET /hello

```
curl -i http://localhost:7777/hello
```

- POST /hello

```
$ curl -v http://localhost:7777/hello -X POST -d "{"message": "hello"}"
```

## 参考

- [Go でシンプルな HTTP サーバを自作する](https://qiita.com/tutuz/items/ab1fd3c0ee6fa01e08b6#%E3%82%AF%E3%83%A9%E3%82%A4%E3%82%A2%E3%83%B3%E3%83%88%E3%81%8B%E3%82%89%E3%81%AE%E6%8E%A5%E7%B6%9A%E3%82%92%E5%BE%85%E3%81%A1%E5%8F%97%E3%81%91%E3%82%8B)
