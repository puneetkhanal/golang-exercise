curl --data "password=angryMonkey" http://localhost:9090/hash
curl --data "password=hello" http://localhost:9090/hash
curl --data "password=hi" http://localhost:9090/hash
curl --data "password=test" http://localhost:9090/hash

# shutdown should wait for all tasks to finish
curl -i http://localhost:9090/shutdown