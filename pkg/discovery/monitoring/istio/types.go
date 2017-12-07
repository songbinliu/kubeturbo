package istio

import (
	"bytes"
	"fmt"
)

type ObjectMetric struct {
	UID              string  `json:"uid,omitempty"`
	Latency          float64 `json:"response_time,omitempty"`
	RequestPerSecond float64 `json:"req_per_second,omitempty"`
}

func (o *ObjectMetric) String() string {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(fmt.Sprintf("latency=%.5f, rps=%.5f", o.Latency, o.RequestPerSecond))
	if len(o.UID) > 0 {
		buffer.WriteString(fmt.Sprintf(", uid=%v", o.UID))
	}

	return buffer.String()
}

func NewObjectMetric(uid string, resp, rps float64) *ObjectMetric {
	return &ObjectMetric{
		UID:              uid,
		Latency:          resp,
		RequestPerSecond: rps,
	}
}

type MetricSet map[string]*ObjectMetric

func NewMetricSet() MetricSet {
	ret := make(MetricSet)
	return ret
}

/*
func (mset MetricSet) AddMetric(uid string, resp, rps float64) {
	obj := NewObjectMetric(uid, resp, rps)
	obj.UID = ""
	mset[uid] = obj
}

func (mset MetricSet) AddorSetResponeTime(uid string, resp float64) {
	obj, exist := mset[uid]
	if !exist {
		mset.AddMetric(uid, resp, 0.0)
	} else {
		if obj.Latency < resp {
			obj.Latency = resp
		}
	}
}

func (mset MetricSet) AddorSetRPS(uid string, rps float64) {
	obj, exist := mset[uid]
	if !exist {
		mset.AddMetric(uid, 0.0, rps)
	} else {
		if obj.RequestPerSecond < rps {
			obj.RequestPerSecond = rps
		}
	}
}*/

func (mset MetricSet) String() string {
	buffer := bytes.NewBufferString("")

	buffer.WriteString(fmt.Sprintf("size = %d\n", len(mset)))
	for k, v := range mset {
		buffer.WriteString(fmt.Sprintf("%s, uid=%v\n", v.String(), k))
	}

	return buffer.String()
}
