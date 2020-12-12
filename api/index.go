package api

const INDEX_R = "/"

type IndexParams struct {
}

type IndexResult struct {
	Data string
}

func Index(remoteAddr string, userID string, params interface{}) (interface{}, error) {
	result := &IndexResult{Data: "index"}
	return result, nil
}
