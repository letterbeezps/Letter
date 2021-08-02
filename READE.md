# Letter

A golang http request library implemented using [fasthttp](https://github.com/valyala/fasthttp)
The project mainly refers to [grequest](https://github.com/levigross/grequests).

## Examples

### Get

```golang
ro := &letter.RequestOptions{
    Params: map[string]string{"name": "zp"},
}
myJsonStruct := verifyOkResponse(letter.Get("http://127.0.0.1:6001/get", ro), t)
{
  "args": {
    "name": "zp"
  }, 
  "headers": {
    "Host": "127.0.0.1:6001", 
    "User-Agent": "letterClient"
  }, 
  "origin": "172.17.0.1", 
  "url": "http://127.0.0.1:6001/get?name=zp"
}
```

### POST

```golang
pData := letter.PostData{
    Name: "zp",
    RelationShip: letter.Person{
    Name: "dd",
    },
}
resp := letter.Post("http://127.0.0.1:6001/post", &letter.RequestOptions{JSON: pData, IsAjax: true})

{
  "args": {}, 
  "data": "{\"name\":\"zp\",\"relationship\":{\"name\":\"dd\"}}\n", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Content-Length": "43", 
    "Content-Type": "application/json", 
    "Host": "127.0.0.1:6001", 
    "User-Agent": "letterClient", 
    "X-Requested-With": "XMLHttpRequest"
  }, 
  "json": {
    "name": "zp", 
    "relationship": {
      "name": "dd"
    }
  }, 
  "origin": "172.17.0.1", 
  "url": "http://127.0.0.1:6001/post"
}
```

## TodoList

* PUT, PATCH, HEAD, DELETE, OPTIONS
* gzip
* work pool

## Reference

[Grequest](https://github.com/levigross/grequests)
