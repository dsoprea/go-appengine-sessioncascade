## Overview

This package equips Go with session support under Google AppEngine. It provides a storage implementation that's built on top of [gorilla/sessions](http://www.gorillatoolkit.org/pkg/sessions). Thanks goes to the [boj/redistore](https://github.com/boj/redistore) project for providing an insightful foundation to start from.


## Description

The purpose of this storage implementation is to allow you to implement one or more AppEngine-specific backends to store your session data. For example, the recommended mode is to use both Memcache and Datastore. When you do updates, they will be pushed to both. When you do a retrieval, it'll try retrieving the session from Memcache (faster and cheaper) before trying Datastore. 

There is also a "request" backend, which will store into the current request so that further retrievals within the same request will not incur the cost of hitting the services. However, use of this backend is not recommended as updates from concurrent requests will not be seen from the current request. If the session has already been stored to the current request context then it'll never hit one of the distributed services and you risk stale data.


## Usage

To choose a backend, pass a constant, or a bitwise-OR of more than one, to `NewCascadeStore()`:

- RequestBackend
- MemcacheBackend
- DatastoreBackend

For convenience, constants for commonly OR'd modes are provided:

- DistributedBackends (MemcacheBackend | DatastoreBackend)
- AllBackends (RequestBackend | MemcacheBackend | DatastoreBackend)


## Example

```go
package handlers

import (
    "net/http"

    "google.golang.org/appengine"
    "google.golang.org/appengine/log"

    "github.com/dsoprea/goappenginesessioncascade"
)

const (
    sessionName = "MainSession"
)

var (
    sessionSecret = []byte("SessionSecret")
    sessionStore = cascadestore.NewCascadeStore(cascadestore.DistributedBackends, sessionSecret)
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := appengine.NewContext(r)

    if session, err := sessionStore.Get(r, sessionName); err != nil {
        panic(err)
    } else {
        if vRaw, found := session.Values["ExistingKey"]; found == false {
            log.Debugf(ctx, "Existing value not found.")
        } else {
            v := vRaw.(string)
            log.Debugf(ctx, "Existing value: [%s]", v)
        }

        session.Values["NewKey"] = "NewValue"
        if err := session.Save(r, w); err != nil {
            panic(err)
        }
    }
}
```


## Other Notes

- There is always the chance that you are going to store things bigger than Memcache can handle. By default, AppEngine has a one-megabyte size-limit for values in Memcache. *go-appengine-sessioncascade* will, by default, not try to push to Memcache if the serialized value is larger than this. Since you are probably ignoring Memcache errors, we could probably, technically, just ignore this as well. However, you would see error messages in your logging that might catch you off-guard. So, we prefer preventing them. If you want to change the size limit, define the `MaxMemcacheSessionSizeBytes` environment variable.

- For general session store usage information, see the [gorilla/sessions homepage](http://www.gorillatoolkit.org/pkg/sessions).
- This project uses [go-appengine-logging](https://github.com/dsoprea/go-appengine-logging) for logging and, as a result, provides stacktraces with its errors.
