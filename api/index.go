package api

type IndexParams struct {
}

type IndexResult struct {
	Data string
}

func Index(params interface{}) (interface{}, error) {
	result := &IndexResult{Data: "index"}
	return result, nil
}
