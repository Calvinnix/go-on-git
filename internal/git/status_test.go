package git

import (
	"testing"
)

func TestGetStatus_EmptyRepo(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if !status.IsEmpty() {
		t.Error("expected status to be empty in clean repo")
	}

	if status.TotalFiles() != 0 {
		t.Errorf("expected 0 files, got %d", status.TotalFiles())
	}
}

func TestGetStatus_UntrackedFiles(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.WriteFile("untracked1.txt", "content1")
	repo.WriteFile("untracked2.txt", "content2")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(status.Untracked) != 2 {
		t.Errorf("expected 2 untracked files, got %d", len(status.Untracked))
	}

	if len(status.Staged) != 0 {
		t.Errorf("expected 0 staged files, got %d", len(status.Staged))
	}

	if len(status.Unstaged) != 0 {
		t.Errorf("expected 0 unstaged files, got %d", len(status.Unstaged))
	}

	// Verify files are marked as untracked
	for _, f := range status.Untracked {
		if !f.IsUntracked() {
			t.Errorf("expected file %s to be untracked", f.Path)
		}
		if f.IsStaged() {
			t.Errorf("expected file %s to not be staged", f.Path)
		}
		if f.IsUnstaged() {
			t.Errorf("expected file %s to not be marked as unstaged (it's untracked)", f.Path)
		}
	}
}

func TestGetStatus_StagedFiles(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.WriteFile("staged.txt", "content")
	repo.Git("add", "staged.txt")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(status.Staged) != 1 {
		t.Errorf("expected 1 staged file, got %d", len(status.Staged))
	}

	if len(status.Untracked) != 0 {
		t.Errorf("expected 0 untracked files, got %d", len(status.Untracked))
	}

	if status.Staged[0].Path != "staged.txt" {
		t.Errorf("expected path 'staged.txt', got %q", status.Staged[0].Path)
	}

	if !status.Staged[0].IsStaged() {
		t.Error("expected file to be marked as staged")
	}

	if status.Staged[0].IndexStatus != 'A' {
		t.Errorf("expected IndexStatus 'A', got %q", status.Staged[0].IndexStatus)
	}
}

func TestGetStatus_ModifiedUnstaged(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "modified")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(status.Unstaged) != 1 {
		t.Errorf("expected 1 unstaged file, got %d", len(status.Unstaged))
	}

	if len(status.Staged) != 0 {
		t.Errorf("expected 0 staged files, got %d", len(status.Staged))
	}

	if status.Unstaged[0].WorkStatus != 'M' {
		t.Errorf("expected WorkStatus 'M', got %q", status.Unstaged[0].WorkStatus)
	}

	if !status.Unstaged[0].IsUnstaged() {
		t.Error("expected file to be marked as unstaged")
	}
}

func TestGetStatus_ModifiedStaged(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")
	repo.WriteFile("test.txt", "modified")
	repo.Git("add", "test.txt")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(status.Staged) != 1 {
		t.Errorf("expected 1 staged file, got %d", len(status.Staged))
	}

	if status.Staged[0].IndexStatus != 'M' {
		t.Errorf("expected IndexStatus 'M', got %q", status.Staged[0].IndexStatus)
	}
}

func TestGetStatus_MixedStagedAndUnstaged(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original", "initial")

	// Stage a change
	repo.WriteFile("test.txt", "staged change")
	repo.Git("add", "test.txt")

	// Make another change (unstaged)
	repo.WriteFile("test.txt", "staged change\nwith unstaged addition")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	// File should appear in both staged and unstaged
	if len(status.Staged) != 1 {
		t.Errorf("expected 1 staged file, got %d", len(status.Staged))
	}

	if len(status.Unstaged) != 1 {
		t.Errorf("expected 1 unstaged file, got %d", len(status.Unstaged))
	}

	// Both should reference the same file
	if status.Staged[0].Path != status.Unstaged[0].Path {
		t.Errorf("expected same file in staged and unstaged")
	}

	// TotalFiles should count it only once
	if status.TotalFiles() != 1 {
		t.Errorf("expected TotalFiles to be 1, got %d", status.TotalFiles())
	}
}

func TestGetStatus_DeletedFile(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "content", "initial")
	repo.DeleteFile("test.txt")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(status.Unstaged) != 1 {
		t.Errorf("expected 1 unstaged file, got %d", len(status.Unstaged))
	}

	if status.Unstaged[0].WorkStatus != 'D' {
		t.Errorf("expected WorkStatus 'D', got %q", status.Unstaged[0].WorkStatus)
	}
}

func TestGetStatus_StagedDeleted(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "content", "initial")
	repo.Git("rm", "test.txt")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(status.Staged) != 1 {
		t.Errorf("expected 1 staged file, got %d", len(status.Staged))
	}

	if status.Staged[0].IndexStatus != 'D' {
		t.Errorf("expected IndexStatus 'D', got %q", status.Staged[0].IndexStatus)
	}
}

func TestGetStatus_RenamedFile(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("old-name.txt", "content", "initial")
	repo.Git("mv", "old-name.txt", "new-name.txt")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(status.Staged) != 1 {
		t.Errorf("expected 1 staged file, got %d", len(status.Staged))
	}

	if status.Staged[0].IndexStatus != 'R' {
		t.Errorf("expected IndexStatus 'R', got %q", status.Staged[0].IndexStatus)
	}

	if status.Staged[0].Path != "new-name.txt" {
		t.Errorf("expected path 'new-name.txt', got %q", status.Staged[0].Path)
	}

	if status.Staged[0].OriginalPath != "old-name.txt" {
		t.Errorf("expected OriginalPath 'old-name.txt', got %q", status.Staged[0].OriginalPath)
	}
}

func TestGetStatus_MultipleFileTypes(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("existing.txt", "original", "initial")

	// Create various file states
	repo.WriteFile("untracked.txt", "new file")
	repo.WriteFile("staged-new.txt", "staged content")
	repo.Git("add", "staged-new.txt")
	repo.WriteFile("existing.txt", "modified content")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(status.Untracked) != 1 {
		t.Errorf("expected 1 untracked file, got %d", len(status.Untracked))
	}

	if len(status.Staged) != 1 {
		t.Errorf("expected 1 staged file, got %d", len(status.Staged))
	}

	if len(status.Unstaged) != 1 {
		t.Errorf("expected 1 unstaged file, got %d", len(status.Unstaged))
	}

	if status.TotalFiles() != 3 {
		t.Errorf("expected 3 total files, got %d", status.TotalFiles())
	}
}

func TestFileStatus_StatusDescription(t *testing.T) {
	tests := []struct {
		name        string
		status      FileStatus
		expected    string
	}{
		{
			name: "untracked",
			status: FileStatus{
				IndexStatus: '?',
				WorkStatus:  '?',
			},
			expected: "untracked",
		},
		{
			name: "staged modified",
			status: FileStatus{
				IndexStatus: 'M',
				WorkStatus:  ' ',
			},
			expected: "staged: modified",
		},
		{
			name: "staged added",
			status: FileStatus{
				IndexStatus: 'A',
				WorkStatus:  ' ',
			},
			expected: "staged: added",
		},
		{
			name: "staged deleted",
			status: FileStatus{
				IndexStatus: 'D',
				WorkStatus:  ' ',
			},
			expected: "staged: deleted",
		},
		{
			name: "staged renamed",
			status: FileStatus{
				IndexStatus: 'R',
				WorkStatus:  ' ',
			},
			expected: "staged: renamed",
		},
		{
			name: "staged copied",
			status: FileStatus{
				IndexStatus: 'C',
				WorkStatus:  ' ',
			},
			expected: "staged: copied",
		},
		{
			name: "unstaged modified",
			status: FileStatus{
				IndexStatus: ' ',
				WorkStatus:  'M',
			},
			expected: "modified",
		},
		{
			name: "unstaged deleted",
			status: FileStatus{
				IndexStatus: ' ',
				WorkStatus:  'D',
			},
			expected: "deleted",
		},
		{
			name: "staged and modified",
			status: FileStatus{
				IndexStatus: 'M',
				WorkStatus:  'M',
			},
			expected: "staged: modified, modified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.StatusDescription()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestStatusResult_IsEmpty(t *testing.T) {
	t.Run("empty status", func(t *testing.T) {
		status := &StatusResult{}
		if !status.IsEmpty() {
			t.Error("expected empty status to return true for IsEmpty")
		}
	})

	t.Run("with staged file", func(t *testing.T) {
		status := &StatusResult{
			Staged: []FileStatus{{Path: "test.txt"}},
		}
		if status.IsEmpty() {
			t.Error("expected non-empty status for staged files")
		}
	})

	t.Run("with unstaged file", func(t *testing.T) {
		status := &StatusResult{
			Unstaged: []FileStatus{{Path: "test.txt"}},
		}
		if status.IsEmpty() {
			t.Error("expected non-empty status for unstaged files")
		}
	})

	t.Run("with untracked file", func(t *testing.T) {
		status := &StatusResult{
			Untracked: []FileStatus{{Path: "test.txt"}},
		}
		if status.IsEmpty() {
			t.Error("expected non-empty status for untracked files")
		}
	})
}

func TestGetStatus_Subdirectories(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.WriteFile("dir1/file1.txt", "content1")
	repo.WriteFile("dir1/dir2/file2.txt", "content2")
	repo.Git("add", "-A")

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if len(status.Staged) != 2 {
		t.Errorf("expected 2 staged files, got %d", len(status.Staged))
	}

	// Check paths include subdirectories
	foundFile1 := false
	foundFile2 := false
	for _, f := range status.Staged {
		if f.Path == "dir1/file1.txt" {
			foundFile1 = true
		}
		if f.Path == "dir1/dir2/file2.txt" {
			foundFile2 = true
		}
	}

	if !foundFile1 {
		t.Error("expected to find dir1/file1.txt")
	}
	if !foundFile2 {
		t.Error("expected to find dir1/dir2/file2.txt")
	}
}
