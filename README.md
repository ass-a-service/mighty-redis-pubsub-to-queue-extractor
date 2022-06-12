# Mighty Redis Pub/Sub to Queue extractor

MRPSTQE doesn't really sound like a name our marketing department
could use, so please if you can think of a cool name let us know.

## Purpose

A simple yet mighty Redis Pub/sub to Queue extractor written in Go.
It's meant to [Web Scale](https://www.youtube.com/watch?v=b2F-DItXtZs).

## Usage

Define the following environment variables and GO:

- REDIS_HOST: I know this is hard to believe, but this variable
actually is the Redis host.
- REDIS_PORT: Redis port but with an underscore between the
words so our OCD doesn't spike.
- REDIS_CHANNEL: The name of the redis channel to subscribe to.
- DEBUG: Set to 1 to get MORE verbosity (is it ever enough?).
Spoiler: setting it to 2 won't give you more verbosity just yet.
- REDIS_QUEUE: The name of the redis queue to queue all
the messages received from REDIS_CHANNEL.
- QUEUE_TYPE: Allowed values are `FIFO` and `LIFO`, which stands
for First In F*** off and Last In Fast Oatmeals, respectively.