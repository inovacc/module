package git

import "testing"

func TestNewGitBareStorage(t *testing.T) {
	bare, err := NewGitBareStorage()
	if err != nil {
		t.Fatalf("go: %v", err)
	}

	repo := bare.Repo()
	if repo == nil {
		t.Fatalf("go: %v", err)
	}

	if err = repo.SetRemote("golang.org/x/tools/gopls"); err != nil {
		t.Fatalf("go: %v", err)
	}

	if err = repo.Fetch(); err != nil {
		t.Fatalf("go: %v", err)
	}

	list, err := repo.GetTags()
	if err != nil {
		t.Fatalf("go: %v", err)
	}

	t.Log("tags:", list)
}
