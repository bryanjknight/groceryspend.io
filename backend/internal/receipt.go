package internal

type ParseReceiptRequest struct {
	ID    string `json:"id"`
	S3Loc string `json:"s3Loc"`
}

type ParseReceiptResult struct {
	S3Loc string `json:"s3Loc"`
}
