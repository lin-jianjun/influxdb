package query

import (
	"github.com/influxdata/flux"
	platform "github.com/influxdata/influxdb"
)

// BucketAwareOperationSpec specifies an operation that reads or writes buckets
type BucketAwareOperationSpec interface {
	BucketsAccessed(orgID *platform.ID) (readBuckets, writeBuckets []platform.BucketFilter)
}

// BucketsAccessed returns the set of buckets read and written by a query spec
func BucketsAccessed(q *flux.Spec, orgID *platform.ID) (readBuckets, writeBuckets []platform.BucketFilter, err error) {
	err = q.Walk(func(o *flux.Operation) error {
		bucketAwareOpSpec, ok := o.Spec.(BucketAwareOperationSpec)
		if ok {
			opBucketsRead, opBucketsWritten := bucketAwareOpSpec.BucketsAccessed(orgID)
			readBuckets = append(readBuckets, opBucketsRead...)
			writeBuckets = append(writeBuckets, opBucketsWritten...)
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return readBuckets, writeBuckets, nil
}
