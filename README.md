# API RATE LIMITING

This is a API rate limiting library written in Go. The following algorithms are implemented curently to throttle user request (you can request any specific algorithms you need.) - Token Bucket

## TOKEN BUCKET ALGORITHM

Token Bucket is a container with a pre-defined capacity (1). Tokens are added to the bucket at a pre-defined rate (1 token every 10 sec). Once the bucket reached maximum capacity no more tokens can be added.

Each API request will consume 1 token from the bucket. If there are no tokens left in the bucket the request is dropped.

This algorithm persits following user data in memory

- **User ID** unique identifier for every user.
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

    limit "github.com/Kulvir-parmar/ratelimit/xxx"
)

func main() {
    Users := limit.NewTokenBuckets()
    router := http.NewServeMux()

    http.ListenAndServe(":42069", limit.RateLimiter(Users, router))
}
```

`xxx` can be limiting algorithm of your choice.
But currently we have only `tokenbucket` implemented. WIP other algorithms.

### IMPORTANT

- `Users` is in memory representation of Token Bucket Store.
- Creating a Bucket Store and pass it to Middleware function is important step.
- Create a `Router` and attach all the routes on which you want Rate Limiting.
  - By default each user is allowed 1 request per 10 seconds (copied this from Leetcode).
  - To uniquely identify each user `userId` should be passed (!IMPORTANT)
  - If request does not have `userId` 400 BAD REQUEST is thrown.

**Memory Cleanup**

As the users grow (jox, coz your don't have any user) the Bucket Store grows and It require much more memory.
To manage the memory, we need to cleanup the Buckets of inactive users. So you can spawn a go routine in our application that will take care of cleaning up the memory.

**Use this code snippet in your main function**

**NOTE**: I'm not the best guy to write this kind of code. But I hope this works. If you have better ways use it or you can make a pull request for that, much appreciated.

```go
go func () {
    for {
        Users.ClearOldBuckets()

        time.Sleep(4 * time.Hour)
    }
}()
```

4 Hours is just a random number you can choose anything that works for you.

---

Wanna contact me ?? [twitter.com/Kulvirdotgg](https://twitter.com/kulvirdotgg)
