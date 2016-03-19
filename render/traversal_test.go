package render

import (
	"github.com/longkai/xiaolongtongxue.com/env"
	"os"
	"testing"
)

func TestTraversal(t *testing.T) {
	Traversal(env.Config().ArticleRepo)
}

func TestCp(t *testing.T) {
	dest := "_____tmp.ttt"
	src := "traversal_test.go"
	err := copyFile(src, dest)
	if err != nil {
		t.Errorf("copy from %s to %s fail, %v\n", src, dest)
	}

	_, err = os.Stat(dest)
	if os.IsNotExist(err) {
		t.Errorf("dest %s not exist!\n", dest)
	}

	// delete it
	if err = os.Remove(dest); err != nil {
		t.Errorf("remove tmp dest file fail, %v\n", err)
	}
}
