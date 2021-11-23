# DistSys Project start-up

this project contains : [literature](./literature.pdf), [design](./design.pdf), and [final report](./final report.pdf)

### Implementation Focused

Implement a basic distributed cache system using Golang from scratch.



The results of some time-consuming operations are temporarily stored for the first request, and the temporarily stored data is directly returned when the same request is encountered later. This is a simple  understanding of caching.

In computer systems, caching is ubiquitous. For example, when we visit a webpage, webpages and referenced JS/CSS and other static files, according to different strategies, they will be cached locally in the browser or on the CDN server. Sometimes, you feel that the web page loads a lot faster; for example, the number of likes on twitter, it is impossible for everyone to find all the likes in the database and count them every time they visit. The operation of the database is very time-consuming. It is difficult to support such a large amount of traffic, so generally this kind of data like praise is cached in the Redis service cluster.



What we want to do is do implement a basic distributed cache system using Golang from scratch. There are some developed commercial distributed cache system like redis, memcached. We decide to build a basic cache system with fundamental features using Golang. The features will at least include:

1. Stand-alone cache and HTTP-based distributed cache
2. Caching strategy
3. Achieve load balancing



### Challenge

To design a distributed cache system, it is necessary to consider various aspects of resource control, elimination strategy, concurrency, and distributed node communication. Moreover, for different application scenarios, different features need to be weighed, and different trade-offs correspond to different implementations.

##### We need to implement a reasonable elimination strategy.

When there is no enough memory, we need to delete some data. Which strategy is the best, randomly or follow the time sequence, or there is a better method?



##### Concurrent write conflict

Access to the cache is generally impossible to be serial. To deal with concurrent scenarios, modification operations (including adding, updating and deleting) need to be locked.



### Papers

> Consistent Hashing and Random Trees: Distributed Caching Protocols for Relieving Hot Spots on the World Wide Web

> Scaling Memcache at Facebook

> DistCache: Provable Load Balancing for LargeScale Storage Systems with Distributed Caching

> Sales ranks, Burgers-like equations, and least-recently-used caching

> Least-recently-used caching with dependent requests



