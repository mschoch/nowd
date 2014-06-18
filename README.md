# nowd

the now daemon, an http server for aggregating sensor data.

## features

* sensor messages are submitted in JSON
* messages expire after configurable timeout
* root url reports on all current (non-expired) sensor readings

## build

    go get github.com/mschoch/nowd

## run

    nowd

## submit sensor data

    curl -XPOST http://localhost:4793/outside -d '{"t":30.0}'
    curl -XPOST http://localhost:4793/inside -d '{"t":20.0}'

## read current state

    curl http://localhost:4793/
    {
    	"outside": {"t":30.0},
    	"inside": {"t":20.0}
    }