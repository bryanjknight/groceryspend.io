package internal

type ParseReceiptRequest struct {
}

type ParseReceiptResponse struct {
	S3Loc string `json:"s3Loc"`
}

func ListParseReceiptRequests() []ParseReceiptResponse {
	return []ParseReceiptResponse{
		{
			S3Loc: "abc",
		},
	}
}
