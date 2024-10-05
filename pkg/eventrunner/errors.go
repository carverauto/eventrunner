package eventrunner

import "errors"

var (
	errBufferForcedWriteError = errors.New("forced write error")
	errPublishError           = errors.New("publish error")
	errConsumeEventError      = errors.New("consume event error")
	errConsumerFailed         = errors.New("consumer failed")
	errNilContext             = errors.New("nil context provided")
	errNilEvent               = errors.New("nil event provided")
	errNilCassandra           = errors.New("cassandra client is nil in CassandraEventSink")
)
