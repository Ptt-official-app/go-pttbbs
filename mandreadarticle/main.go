package main

import (
	"context"

	"github.com/Ptt-official-app/go-pttbbs/mand"
	"github.com/sirupsen/logrus"
)

func main() {
	brdname, path, err := initMain()
	if err != nil {
		logrus.Fatalf("unable to initMain: e: %v", err)
		return
	}

	ctx := context.Background()

	req := &mand.ArticleRequest{
		BoardName: brdname,
		Path:      path,
		MaxLength: -1,
	}

	resp, err := mand.Cli.Article(ctx, req)
	if err != nil {
		logrus.Errorf("unable to get content: e: %v", err)
		return
	}

	logrus.Infof("brdname: %v path: %v content: %v", brdname, path, string(resp.Content))
}
