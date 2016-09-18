# Blogk

"Blockchain backed microblogging"

(will try to add more later, running out of time in this hackathon)

## Usage

0. Install go...

1. Install dependancies
```
$ go get -u ./...
```

2. Run the `blogk` server
```
$ go run cmd/blogk/main.go
```

3. In a new shell, try creating a user or a post
Make a user (uses a generated public/private key that are not saved)
```
go run cmd/examples/add_user/main.go
```

Make a post
```
go run cmd/examples/make_post/main.go
```

check back later!