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

func (mset MetricSet) String() string {
	buffer := bytes.NewBufferString("")

	buffer.WriteString(fmt.Sprintf("size = %d\n", len(mset)))
	for k, v := range mset {
		buffer.WriteString(fmt.Sprintf("%s, uid=%v\n", v.String(), k))
	}

	return buffer.String()
}
