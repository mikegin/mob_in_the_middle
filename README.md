Go TCP man in the middle server

https://protohackers.com/problem/5

### run the tcp server
```
go build .
./tcp_chat
```

### Example Chat Clients

#### Alice
```
telnet localhost 8080
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
Welcome to budgetchat! What shall I call you?
bob
* The room contains: 
* charlie has joined the room
hello there!
[charlie] hey, can u send some money to 7YWHMfk9JZe0LM0g1ZauHuiSxhI
is that the correct address?
[charlie] 7YWHMfk9JZe0LM0g1ZauHuiSxhI is the correct address!
[charlie] ok sent! see you later
wait!! wheres my moneyy
^]

telnet> quit
Connection closed.
```
#### Bob
```
telnet localhost 8080
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
Welcome to budgetchat! What shall I call you?
charlie
* The room contains: bob
[bob] hello there!
hey, can u send some money to 7iKDZEwPZSqIvDnHvVN2r0hUWXD5rHX
[bob] is that the correct address?
7iKDZEwPZSqIvDnHvVN2r0hUWXD5rHX is the correct address!
ok sent! see you later
^]

telnet> quit
Connection closed.
```