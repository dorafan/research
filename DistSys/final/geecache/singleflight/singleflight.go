package singleflight

import "sync"

// call is an in-flight or completed Do call
//Represents an ongoing or closed request. Use sync.WaitGroup lock to avoid reentry
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression.
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}

// Do executes and returns the results of the given function, making
// sure that only one execution is in-flight for a given key at a
// time. If a duplicate comes in, the duplicate caller waits for the
// original to complete and receives the same results.
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		//If the request is in progress, wait
		c.wg.Wait()
		//End of request, return result
		return c.val, c.err
	}
	c := new(call)
	//Lock before request
	c.wg.Add(1)
	//Add to g.m, indicating that the key has a corresponding request being processed
	g.m[key] = c
	g.mu.Unlock()

	//Call fn to initiate a request
	c.val, c.err = fn()
	//End of request
	c.wg.Done()

	g.mu.Lock()
	//update g.m
	delete(g.m, key)
	g.mu.Unlock()

	//return result
	return c.val, c.err
}
