# API RATE LIMITING

This is a API rate limiting library written in Go.

- some reference [cloudflare-blog](https://blog.cloudflare.com/counting-things-a-lot-of-different-things/)

## TOKEN BUCKET ALGORITHM

Token Bucket is a container with a pre-defined capacity (2). Tokens are added to the bucket at a pre-defined rate (1 token every 10 sec). Once the bucket reached maximum capacity no more tokens can be added.

Each API request will consume 1 token from the bucket. If there are no tokens left in the bucket the request is dropped.

This algorithm persits following user data in memory

- **IP Address** unique identifier for every user.
- **Tokens left** in the bucket for user.
- **Timestamp** for when request are sent.

```go
type Bucket struct {
    mu         sync.Mutex
    capacity   TOKENS        // max tokens allowed
    tokens     TOKENS        // tokens present in the bucket
    lastRefill time.Time
    rate       time.Duration // rate at which tokens are refilled
}
```

In memory storage is used because if we use a Database reading from the disk can be slow and we can use time based expiration for the old user data.

Although token bucket is pretty straight forward algorithm but may cause race condition in a distributed system due to concurent request from same user.

<br />

```go
import (
    "net/http"

    limit "github.com/Kulvir-parmar/ratelimit/tokenbucket"
)

func main() {
    Users := limit.NewTokenBuckets()
    router := http.NewServeMux()

    http.ListenAndServe(":42069", limit.RateLimiter(router, Users))
}
```

### IMPORTANT

- `Users` is in memory representation of Token Bucket Store.
- Creating a Bucket Store and pass it to Middleware function is important step.
- Create a `Router` and attach all the routes on which you want Rate Limiting.
  - By default each user is allowed 2 request per 10 seconds (something like LEETCODE).
  - To uniquely identify each IP Address of each request is tracked.

**Memory Cleanup**

As the users grow (jox, coz you don't have any users) the Bucket Store grows which require more memory.
To manage the memory, we need to cleanup the Buckets of inactive users. So you can spawn a go routine in our application that will take care of cleaning up the memory.

**Use this code snippet in your main function**

**NOTE**: I'm not the best guy to write this kind of code. But I hope this works. If you have better ways use it or you can make a pull request for that, much appreciated.

```go
go func () {
    for {
        Users.ClearOldBuckets()

        time.Sleep(10 * time.Minute)
    }
}()
```

---

Wanna contact me ?? [twitter.com/Kulvirdotgg](https://twitter.com/kulvirdotgg)
