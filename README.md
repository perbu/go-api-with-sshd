# A Golang API with a built-in ssh server

This is made for a talk about Go and ssh and not really meant to be used. 

## Compilation

Make sure to put a file with your authorized keys in the "backdoor" directory. It'll be compiled into the application.


## API Usage:

Access the users:
```shell
http http://localhost:8080/user/Bob
http http://localhost:8080/user/Peter
http http://localhost:8080/user/Alice
```

Let's give Bob another pet:
```shell
http POST http://localhost:8080/user/Bob/addpet name=Fluffy
```
