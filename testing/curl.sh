curl -X POST -d '{"age":18,"name":"scguo"}' http://127.0.0.1:8090/demo

echo ""

curl -X POST -d '{"name":"scguo","age":1}{"name":"scguo","age":18}' -H "Transfer-Encoding: chunked" 127.0.0.1:8090/stream