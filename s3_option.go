package storage

type S3OptionFunc func(*S3ObjectStorage) error

func S3Endpoint(url string) S3OptionFunc {
	return func(s3 *S3ObjectStorage) error {
		s3.Endpoint = url
		return nil
	}
}

func S3Region(region string) S3OptionFunc {
	return func(s3 *S3ObjectStorage) error {
		s3.Region = region
		return nil
	}
}
