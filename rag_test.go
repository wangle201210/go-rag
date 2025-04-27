package rag

import "testing"

func TestIndex(t *testing.T) {
	ids, err := Index("./test_file/readme.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range ids {
		t.Log(id)
	}
	Index("./test_file/readme2.md")
	Index("./test_file/readme.html")
	Index("./test_file/test.pdf")
	Index("https://deepchat.thinkinai.xyz/docs/guide/advanced-features/shortcuts.html")
}

func TestRetriever(t *testing.T) {
	msg, err := Retrieve("这里有很多内容", 1.5)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range msg {
		t.Logf("content: %v, score: %v", m.Content, m.Score())
	}

	msg, err = Retrieve("代码解析", 1.5)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range msg {
		t.Logf(" content: %v, score: %v", m.Content, m.Score())
	}
}
