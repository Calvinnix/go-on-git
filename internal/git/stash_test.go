package git

import (
	"strings"
	"testing"
)

func TestGetStashes_NoStashes(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	stashes, err := GetStashes()
	if err != nil {
		t.Fatalf("GetStashes failed: %v", err)
	}

	if len(stashes) != 0 {
		t.Errorf("expected 0 stashes, got %d", len(stashes))
	}
}

func TestGetStashes_SingleStash(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "modified")
	repo.Git("stash", "push", "-m", "Test stash message")

	stashes, err := GetStashes()
	if err != nil {
		t.Fatalf("GetStashes failed: %v", err)
	}

	if len(stashes) != 1 {
		t.Fatalf("expected 1 stash, got %d", len(stashes))
	}

	stash := stashes[0]
	if stash.Index != 0 {
		t.Errorf("expected Index 0, got %d", stash.Index)
	}

	if stash.Message != "Test stash message" {
		t.Errorf("expected message 'Test stash message', got %q", stash.Message)
	}
}

func TestGetStashes_MultipleStashes(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")

	// Create multiple stashes
	repo.WriteFile("test.txt", "change1")
	repo.Git("stash", "push", "-m", "First stash")

	repo.WriteFile("test.txt", "change2")
	repo.Git("stash", "push", "-m", "Second stash")

	repo.WriteFile("test.txt", "change3")
	repo.Git("stash", "push", "-m", "Third stash")

	stashes, err := GetStashes()
	if err != nil {
		t.Fatalf("GetStashes failed: %v", err)
	}

	if len(stashes) != 3 {
		t.Fatalf("expected 3 stashes, got %d", len(stashes))
	}

	// Stashes are in LIFO order, so newest is at index 0
	if stashes[0].Index != 0 {
		t.Errorf("expected first stash to have index 0, got %d", stashes[0].Index)
	}
	if stashes[0].Message != "Third stash" {
		t.Errorf("expected first stash message 'Third stash', got %q", stashes[0].Message)
	}

	if stashes[2].Index != 2 {
		t.Errorf("expected last stash to have index 2, got %d", stashes[2].Index)
	}
	if stashes[2].Message != "First stash" {
		t.Errorf("expected last stash message 'First stash', got %q", stashes[2].Message)
	}
}

func TestGetStashes_Branch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "modified")
	repo.Git("stash", "push", "-m", "Stash on master")

	stashes, err := GetStashes()
	if err != nil {
		t.Fatalf("GetStashes failed: %v", err)
	}

	if len(stashes) == 0 {
		t.Fatal("expected at least one stash")
	}

	// Branch should be captured
	if stashes[0].Branch != "master" && stashes[0].Branch != "main" {
		t.Errorf("expected branch 'master' or 'main', got %q", stashes[0].Branch)
	}
}

func TestGetStashes_WIPStash(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "modified")

	// Stash without message creates a WIP stash
	repo.Git("stash", "push")

	stashes, err := GetStashes()
	if err != nil {
		t.Fatalf("GetStashes failed: %v", err)
	}

	if len(stashes) == 0 {
		t.Fatal("expected at least one stash")
	}

	// WIP stash should still have branch info
	if stashes[0].Branch == "" {
		t.Error("expected branch to be set for WIP stash")
	}
}

func TestGetStashDiff(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original\n", "initial")
	repo.WriteFile("test.txt", "modified\n")
	repo.Git("stash", "push", "-m", "Test stash")

	diff, err := GetStashDiff(0)
	if err != nil {
		t.Fatalf("GetStashDiff failed: %v", err)
	}

	if diff.IsEmpty() {
		t.Error("expected non-empty stash diff")
	}

	// Stash diff shows as unstaged diff (what would be applied)
	if diff.UnstagedDiff == nil || diff.UnstagedDiff.IsEmpty() {
		t.Error("expected unstaged diff portion to have content")
	}
}

func TestGetStashDiff_MultipleFiles(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("file1.txt", "original1\n", "initial1")
	repo.CommitFile("file2.txt", "original2\n", "initial2")

	repo.WriteFile("file1.txt", "modified1\n")
	repo.WriteFile("file2.txt", "modified2\n")
	repo.Git("stash", "push", "-m", "Multi-file stash")

	diff, err := GetStashDiff(0)
	if err != nil {
		t.Fatalf("GetStashDiff failed: %v", err)
	}

	if len(diff.UnstagedDiff.Files) != 2 {
		t.Errorf("expected 2 files in stash diff, got %d", len(diff.UnstagedDiff.Files))
	}
}

func TestApplyStash(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "stashed content")
	repo.Git("stash", "push", "-m", "Test stash")

	// Verify file is restored to original
	content := repo.ReadFile("test.txt")
	if content != "original" {
		t.Fatalf("expected file to be restored to 'original', got %q", content)
	}

	err := ApplyStash(0)
	if err != nil {
		t.Fatalf("ApplyStash failed: %v", err)
	}

	// Verify stashed changes are applied
	content = repo.ReadFile("test.txt")
	if content != "stashed content" {
		t.Errorf("expected 'stashed content', got %q", content)
	}

	// Verify stash is still present
	stashes, _ := GetStashes()
	if len(stashes) != 1 {
		t.Errorf("expected stash to still exist after apply, got %d stashes", len(stashes))
	}
}

func TestApplyStash_SpecificIndex(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")

	// Create multiple stashes
	repo.WriteFile("test.txt", "first change")
	repo.Git("stash", "push", "-m", "First")

	repo.WriteFile("test.txt", "second change")
	repo.Git("stash", "push", "-m", "Second")

	// Apply stash@{1} (the first/older stash)
	err := ApplyStash(1)
	if err != nil {
		t.Fatalf("ApplyStash(1) failed: %v", err)
	}

	content := repo.ReadFile("test.txt")
	if content != "first change" {
		t.Errorf("expected 'first change', got %q", content)
	}
}

func TestPopStash(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "stashed content")
	repo.Git("stash", "push", "-m", "Test stash")

	// Verify we have a stash
	stashes, _ := GetStashes()
	if len(stashes) != 1 {
		t.Fatalf("expected 1 stash before pop, got %d", len(stashes))
	}

	err := PopStash(0)
	if err != nil {
		t.Fatalf("PopStash failed: %v", err)
	}

	// Verify stashed changes are applied
	content := repo.ReadFile("test.txt")
	if content != "stashed content" {
		t.Errorf("expected 'stashed content', got %q", content)
	}

	// Verify stash is removed
	stashes, _ = GetStashes()
	if len(stashes) != 0 {
		t.Errorf("expected stash to be removed after pop, got %d stashes", len(stashes))
	}
}

func TestPopStash_MiddleStash(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")

	// Create multiple stashes
	repo.WriteFile("test.txt", "first")
	repo.Git("stash", "push", "-m", "First")

	repo.WriteFile("test.txt", "second")
	repo.Git("stash", "push", "-m", "Second")

	repo.WriteFile("test.txt", "third")
	repo.Git("stash", "push", "-m", "Third")

	// Pop the middle stash (index 1 = "Second")
	err := PopStash(1)
	if err != nil {
		t.Fatalf("PopStash(1) failed: %v", err)
	}

	// Verify correct changes were applied
	content := repo.ReadFile("test.txt")
	if content != "second" {
		t.Errorf("expected 'second', got %q", content)
	}

	// Verify we have 2 stashes remaining
	stashes, _ := GetStashes()
	if len(stashes) != 2 {
		t.Errorf("expected 2 stashes remaining, got %d", len(stashes))
	}
}

func TestDropStash(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "stashed content")
	repo.Git("stash", "push", "-m", "Test stash")

	// Verify we have a stash
	stashes, _ := GetStashes()
	if len(stashes) != 1 {
		t.Fatalf("expected 1 stash before drop, got %d", len(stashes))
	}

	err := DropStash(0)
	if err != nil {
		t.Fatalf("DropStash failed: %v", err)
	}

	// Verify stash is removed
	stashes, _ = GetStashes()
	if len(stashes) != 0 {
		t.Errorf("expected stash to be removed after drop, got %d stashes", len(stashes))
	}

	// Verify file is still at original (not modified by drop)
	content := repo.ReadFile("test.txt")
	if content != "original" {
		t.Errorf("expected 'original' (drop shouldn't apply changes), got %q", content)
	}
}

func TestDropStash_SpecificIndex(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")

	// Create multiple stashes
	repo.WriteFile("test.txt", "first")
	repo.Git("stash", "push", "-m", "First")

	repo.WriteFile("test.txt", "second")
	repo.Git("stash", "push", "-m", "Second")

	// Drop stash@{1} (the "First" stash)
	err := DropStash(1)
	if err != nil {
		t.Fatalf("DropStash(1) failed: %v", err)
	}

	stashes, _ := GetStashes()
	if len(stashes) != 1 {
		t.Errorf("expected 1 stash remaining, got %d", len(stashes))
	}

	// The remaining stash should be "Second"
	if stashes[0].Message != "Second" {
		t.Errorf("expected remaining stash to be 'Second', got %q", stashes[0].Message)
	}
}

func TestApplyStash_InvalidIndex(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "stashed")
	repo.Git("stash", "push", "-m", "Only stash")

	// Try to apply a non-existent stash
	err := ApplyStash(99)
	if err == nil {
		t.Error("expected error for invalid stash index")
	}
}

func TestPopStash_InvalidIndex(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "stashed")
	repo.Git("stash", "push", "-m", "Only stash")

	err := PopStash(99)
	if err == nil {
		t.Error("expected error for invalid stash index")
	}
}

func TestDropStash_InvalidIndex(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "stashed")
	repo.Git("stash", "push", "-m", "Only stash")

	err := DropStash(99)
	if err == nil {
		t.Error("expected error for invalid stash index")
	}
}

func TestApplyStash_WithConflict(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "stashed content")
	repo.Git("stash", "push", "-m", "Test stash")

	// Make a conflicting change
	repo.WriteFile("test.txt", "conflicting content")

	err := ApplyStash(0)
	// This may or may not produce an error depending on git's merge behavior
	// Git may auto-merge or report a conflict
	_ = err // We're just testing that it doesn't crash
}

func TestStash_OnDifferentBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.Git("checkout", "-b", "feature")
	repo.WriteFile("test.txt", "feature changes")
	repo.Git("stash", "push", "-m", "Feature stash")

	stashes, _ := GetStashes()
	if len(stashes) == 0 {
		t.Fatal("expected at least one stash")
	}

	if stashes[0].Branch != "feature" {
		t.Errorf("expected branch 'feature', got %q", stashes[0].Branch)
	}
}

func TestApplyStash_OnDifferentBranch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "stashed on master")
	repo.Git("stash", "push", "-m", "Master stash")

	// Create and switch to a new branch
	repo.Git("checkout", "-b", "other-branch")

	// Apply stash on different branch
	err := ApplyStash(0)
	if err != nil {
		t.Fatalf("ApplyStash on different branch failed: %v", err)
	}

	content := repo.ReadFile("test.txt")
	if content != "stashed on master" {
		t.Errorf("expected 'stashed on master', got %q", content)
	}
}

func TestStash_NewFile(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.WriteFile("new-file.txt", "new content")
	repo.Git("add", "new-file.txt")
	repo.Git("stash", "push", "-m", "Stash with new file")

	// Verify new file is gone
	if repo.FileExists("new-file.txt") {
		t.Error("expected new file to be stashed away")
	}

	// Apply stash
	ApplyStash(0)

	// Verify new file is back
	if !repo.FileExists("new-file.txt") {
		t.Error("expected new file to be restored")
	}

	content := repo.ReadFile("new-file.txt")
	if content != "new content" {
		t.Errorf("expected 'new content', got %q", content)
	}
}

func TestGetStashDiff_InvalidIndex(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	_, err := GetStashDiff(0)
	if err == nil {
		t.Error("expected error for stash diff with no stashes")
	}
}

func TestStash_MessageParsing(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "modified")

	// Test various message formats
	repo.Git("stash", "push", "-m", "Message with: colons: in: it")

	stashes, err := GetStashes()
	if err != nil {
		t.Fatalf("GetStashes failed: %v", err)
	}

	if len(stashes) == 0 {
		t.Fatal("expected at least one stash")
	}

	// Message should be parsed correctly despite colons
	if !strings.Contains(stashes[0].Message, "colons") {
		t.Errorf("expected message to contain 'colons', got %q", stashes[0].Message)
	}
}
