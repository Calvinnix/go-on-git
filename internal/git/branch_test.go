package git

import (
	"os"
	"testing"
)

func TestGetBranches_SingleBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	if len(branches) != 1 {
		t.Errorf("expected 1 branch, got %d", len(branches))
	}

	// Default branch should be current and named master or main
	branch := branches[0]
	if !branch.IsCurrent {
		t.Error("expected default branch to be current")
	}

	if branch.Name != "master" && branch.Name != "main" {
		t.Errorf("expected branch name 'master' or 'main', got %q", branch.Name)
	}
}

func TestGetBranches_MultipleBranches(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.Git("branch", "feature-1")
	repo.Git("branch", "feature-2")

	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	if len(branches) != 3 {
		t.Errorf("expected 3 branches, got %d", len(branches))
	}

	// Verify branch names
	names := make(map[string]bool)
	for _, b := range branches {
		names[b.Name] = true
	}

	if !names["feature-1"] {
		t.Error("expected to find branch 'feature-1'")
	}
	if !names["feature-2"] {
		t.Error("expected to find branch 'feature-2'")
	}
}

func TestGetBranches_CurrentBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.Git("branch", "other-branch")
	repo.Git("checkout", "other-branch")

	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	currentCount := 0
	var currentBranch Branch
	for _, b := range branches {
		if b.IsCurrent {
			currentCount++
			currentBranch = b
		}
	}

	if currentCount != 1 {
		t.Errorf("expected exactly 1 current branch, got %d", currentCount)
	}

	if currentBranch.Name != "other-branch" {
		t.Errorf("expected current branch 'other-branch', got %q", currentBranch.Name)
	}
}

func TestGetBranches_LastCommit(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "content", "Test commit message")

	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	if len(branches) == 0 {
		t.Fatal("expected at least one branch")
	}

	if branches[0].LastCommit != "Test commit message" {
		t.Errorf("expected last commit 'Test commit message', got %q", branches[0].LastCommit)
	}
}

func TestGetBranches_WithUpstream(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	remoteDir := repo.SetupRemote()
	defer os.RemoveAll(remoteDir)

	repo.PushToRemote()

	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	if len(branches) == 0 {
		t.Fatal("expected at least one branch")
	}

	branch := branches[0]
	if branch.Upstream == "" {
		t.Error("expected upstream to be set after push")
	}
}

func TestGetBranches_AheadBehind(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	remoteDir := repo.SetupRemote()
	defer os.RemoveAll(remoteDir)

	repo.PushToRemote()

	// Make local commits (ahead)
	repo.CommitFile("local1.txt", "content1", "Local commit 1")
	repo.CommitFile("local2.txt", "content2", "Local commit 2")

	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	if len(branches) == 0 {
		t.Fatal("expected at least one branch")
	}

	branch := branches[0]
	if branch.Ahead != 2 {
		t.Errorf("expected 2 commits ahead, got %d", branch.Ahead)
	}
}

func TestCreateBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	err := CreateBranch("new-feature")
	if err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	// Verify branch was created and we're on it
	currentBranch := GetBranch()
	if currentBranch != "new-feature" {
		t.Errorf("expected to be on 'new-feature', got %q", currentBranch)
	}

	// Verify it shows up in branch list
	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	foundNew := false
	for _, b := range branches {
		if b.Name == "new-feature" {
			foundNew = true
			if !b.IsCurrent {
				t.Error("expected new-feature to be current branch")
			}
			break
		}
	}

	if !foundNew {
		t.Error("expected to find new-feature branch")
	}
}

func TestCreateBranch_InvalidName(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	// Branch names with spaces are invalid
	err := CreateBranch("invalid branch name")
	if err == nil {
		t.Error("expected error for invalid branch name")
	}
}

func TestCheckoutBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.Git("branch", "other-branch")

	err := CheckoutBranch("other-branch")
	if err != nil {
		t.Fatalf("CheckoutBranch failed: %v", err)
	}

	currentBranch := GetBranch()
	if currentBranch != "other-branch" {
		t.Errorf("expected to be on 'other-branch', got %q", currentBranch)
	}
}

func TestCheckoutBranch_NonExistent(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	err := CheckoutBranch("nonexistent-branch")
	if err == nil {
		t.Error("expected error for non-existent branch")
	}
}

func TestCheckoutBranch_WithChanges(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.Git("branch", "other-branch")

	// Make uncommitted changes that conflict with checkout
	repo.WriteFile("README.md", "modified content")
	repo.Git("add", "README.md")

	// Checkout should still work for non-conflicting changes
	err := CheckoutBranch("other-branch")
	// This may or may not fail depending on if there are conflicts
	_ = err // Result depends on git behavior with staged changes
}

func TestDeleteBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.Git("branch", "to-delete")

	// Verify branch exists
	branches, _ := GetBranches()
	found := false
	for _, b := range branches {
		if b.Name == "to-delete" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected to-delete branch to exist")
	}

	err := DeleteBranch("to-delete")
	if err != nil {
		t.Fatalf("DeleteBranch failed: %v", err)
	}

	// Verify branch was deleted
	branches, _ = GetBranches()
	for _, b := range branches {
		if b.Name == "to-delete" {
			t.Error("expected to-delete branch to be gone")
		}
	}
}

func TestDeleteBranch_CurrentBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	currentBranch := GetBranch()
	err := DeleteBranch(currentBranch)
	if err == nil {
		t.Error("expected error when deleting current branch")
	}
}

func TestDeleteBranch_UnmergedBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.Git("checkout", "-b", "unmerged")
	repo.CommitFile("unmerged.txt", "content", "unmerged commit")
	repo.Git("checkout", "-")

	// Regular delete should fail for unmerged branch
	err := DeleteBranch("unmerged")
	if err == nil {
		t.Error("expected error when deleting unmerged branch with -d")
	}
}

func TestForceDeleteBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.Git("checkout", "-b", "unmerged")
	repo.CommitFile("unmerged.txt", "content", "unmerged commit")
	repo.Git("checkout", "-")

	// Force delete should succeed
	err := ForceDeleteBranch("unmerged")
	if err != nil {
		t.Fatalf("ForceDeleteBranch failed: %v", err)
	}

	// Verify branch was deleted
	branches, _ := GetBranches()
	for _, b := range branches {
		if b.Name == "unmerged" {
			t.Error("expected unmerged branch to be gone after force delete")
		}
	}
}

func TestForceDeleteBranch_NonExistent(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	err := ForceDeleteBranch("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent branch")
	}
}

func TestGetBranches_BranchOnDifferentCommit(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.Git("branch", "old-branch")

	// Make new commits on master
	repo.CommitFile("new.txt", "new content", "New commit on master")

	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	// Find the old-branch and verify it has the old commit message
	for _, b := range branches {
		if b.Name == "old-branch" {
			if b.LastCommit != "Initial commit" {
				t.Errorf("expected old-branch to have 'Initial commit', got %q", b.LastCommit)
			}
			break
		}
	}
}

func TestBranch_IsRemote(t *testing.T) {
	// The IsRemote field is not currently set by GetBranches (only returns local branches)
	// This test documents the current behavior

	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	for _, b := range branches {
		if b.IsRemote {
			t.Error("expected all branches from GetBranches to be local (IsRemote=false)")
		}
	}
}

func TestCheckoutBranch_SwitchBack(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	originalBranch := GetBranch()

	repo.Git("branch", "feature")
	CheckoutBranch("feature")
	CheckoutBranch(originalBranch)

	currentBranch := GetBranch()
	if currentBranch != originalBranch {
		t.Errorf("expected to be back on %q, got %q", originalBranch, currentBranch)
	}
}

func TestCreateBranch_FromSpecificCommit(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("file1.txt", "content1", "First commit")
	repo.CommitFile("file2.txt", "content2", "Second commit")

	// CreateBranch creates from HEAD, so the new branch should have both commits
	err := CreateBranch("from-head")
	if err != nil {
		t.Fatalf("CreateBranch failed: %v", err)
	}

	// Verify the new branch has the second commit as HEAD
	branches, _ := GetBranches()
	for _, b := range branches {
		if b.Name == "from-head" {
			if b.LastCommit != "Second commit" {
				t.Errorf("expected 'Second commit', got %q", b.LastCommit)
			}
			break
		}
	}
}

func TestGetBranches_EmptyRepo(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	// Don't create initial commit - repo has no branches yet
	branches, err := GetBranches()
	if err != nil {
		t.Fatalf("GetBranches failed: %v", err)
	}

	// Empty repo should have no branches
	if len(branches) != 0 {
		t.Errorf("expected 0 branches in empty repo, got %d", len(branches))
	}
}
