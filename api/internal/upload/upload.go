package upload

import (
	"cloud.google.com/go/storage"
)

type Uploader interface{}

type Client struct {
	ProfileAvatarBucket     *storage.BucketHandle
	ProfileAvatarBucketName string
	FeedPostBucket          *storage.BucketHandle
	FeedPostBucketName      string
	ThemeBucket             *storage.BucketHandle
	ThemeBucketName         string
	ChatMediaBucket         *storage.BucketHandle
	ChatMediaBucketName     string
}

func New(
	profileAvatarBucket *storage.BucketHandle,
	feedPostBucket *storage.BucketHandle,
	themeBucket *storage.BucketHandle,
	chatMediaBucket *storage.BucketHandle,
	profileAvatarBucketName string,
	feedPostBucketName string,
	themBucketName string,
	chatMediaBucketName string,
) *Client {
	return &Client{
		ProfileAvatarBucket:     profileAvatarBucket,
		ProfileAvatarBucketName: profileAvatarBucketName,
		FeedPostBucket:          feedPostBucket,
		FeedPostBucketName:      feedPostBucketName,
		ThemeBucket:             themeBucket,
		ThemeBucketName:         themBucketName,
		ChatMediaBucket:         chatMediaBucket,
		ChatMediaBucketName:     chatMediaBucketName,
	}
}
