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

## why?

I have a lot of sensors reporting information from a lot of sources over a variety of transports.  Often data across multiple sensors at the same time describe the same thing.  For example, I may have separate devices measuring temperature and air pressure.  Readings taken near the same time together both describe the weather.  By aggregating the readings in one place I can have a separate script which queries these aggregated documents and submit useful weather reports, say to Weather Undeground, or my own time-series database.