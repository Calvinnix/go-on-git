package git

import (
	"strings"
	"testing"
)

func TestGetDiff_NoDiff(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if !diff.IsEmpty() {
		t.Error("expected diff to be empty in clean repo")
	}

	if diff.TotalHunks() != 0 {
		t.Errorf("expected 0 hunks, got %d", diff.TotalHunks())
	}
}

func TestGetDiff_ModifiedFile(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "line1\nline2\nline3\n", "initial")
	repo.WriteFile("test.txt", "line1\nmodified\nline3\n")

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if diff.IsEmpty() {
		t.Fatal("expected non-empty diff")
	}

	if len(diff.Files) != 1 {
		t.Fatalf("expected 1 file in diff, got %d", len(diff.Files))
	}

	fileDiff := diff.Files[0]
	if fileDiff.Path != "test.txt" {
		t.Errorf("expected path 'test.txt', got %q", fileDiff.Path)
	}

	if len(fileDiff.Hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(fileDiff.Hunks))
	}

	hunk := fileDiff.Hunks[0]
	// Verify hunk has both removed and added lines
	hasRemoved := false
	hasAdded := false
	for _, line := range hunk.Lines {
		if line.Type == LineRemoved {
			hasRemoved = true
			if !strings.Contains(line.Content, "line2") {
				t.Errorf("expected removed line to contain 'line2', got %q", line.Content)
			}
		}
		if line.Type == LineAdded {
			hasAdded = true
			if !strings.Contains(line.Content, "modified") {
				t.Errorf("expected added line to contain 'modified', got %q", line.Content)
			}
		}
	}

	if !hasRemoved {
		t.Error("expected hunk to have removed lines")
	}
	if !hasAdded {
		t.Error("expected hunk to have added lines")
	}
}

func TestGetDiff_MultipleHunks(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	// Create a file with content that will produce multiple hunks when modified
	original := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\n"
	repo.CommitFile("test.txt", original, "initial")

	// Modify lines far apart to create multiple hunks
	modified := "changed1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nchanged10\n"
	repo.WriteFile("test.txt", modified)

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if len(diff.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(diff.Files))
	}

	// With changes far apart, git should produce multiple hunks
	if len(diff.Files[0].Hunks) < 1 {
		t.Errorf("expected at least 1 hunk, got %d", len(diff.Files[0].Hunks))
	}
}

func TestGetDiff_MultipleFiles(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("file1.txt", "original1", "initial1")
	repo.CommitFile("file2.txt", "original2", "initial2")

	repo.WriteFile("file1.txt", "modified1")
	repo.WriteFile("file2.txt", "modified2")

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if len(diff.Files) != 2 {
		t.Errorf("expected 2 files in diff, got %d", len(diff.Files))
	}
}

func TestGetStagedDiff(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original\n", "initial")
	repo.WriteFile("test.txt", "modified\n")
	repo.Git("add", "test.txt")

	diff, err := GetStagedDiff()
	if err != nil {
		t.Fatalf("GetStagedDiff failed: %v", err)
	}

	if diff.IsEmpty() {
		t.Error("expected staged diff to not be empty")
	}

	if len(diff.Files) != 1 {
		t.Fatalf("expected 1 file in staged diff, got %d", len(diff.Files))
	}

	// Verify unstaged diff is empty
	unstagedDiff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if !unstagedDiff.IsEmpty() {
		t.Error("expected unstaged diff to be empty")
	}
}

func TestGetCombinedDiff(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original\n", "initial")

	// Stage a change
	repo.WriteFile("test.txt", "staged change\n")
	repo.Git("add", "test.txt")

	// Make additional unstaged change
	repo.WriteFile("test.txt", "staged change\nunstaged addition\n")

	combined, err := GetCombinedDiff()
	if err != nil {
		t.Fatalf("GetCombinedDiff failed: %v", err)
	}

	if combined.IsEmpty() {
		t.Error("expected combined diff to not be empty")
	}

	if combined.StagedDiff.IsEmpty() {
		t.Error("expected staged portion to not be empty")
	}

	if combined.UnstagedDiff.IsEmpty() {
		t.Error("expected unstaged portion to not be empty")
	}
}

func TestGetCombinedDiff_AllHunksCombined(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("file1.txt", "original1\n", "initial1")
	repo.CommitFile("file2.txt", "original2\n", "initial2")

	// Stage changes to file1
	repo.WriteFile("file1.txt", "modified1\n")
	repo.Git("add", "file1.txt")

	// Make unstaged changes to file2
	repo.WriteFile("file2.txt", "modified2\n")

	combined, err := GetCombinedDiff()
	if err != nil {
		t.Fatalf("GetCombinedDiff failed: %v", err)
	}

	hunks := combined.GetAllHunksCombined()
	if len(hunks) < 2 {
		t.Errorf("expected at least 2 hunks combined, got %d", len(hunks))
	}

	// Verify hunks have correct Staged flag
	stagedCount := 0
	unstagedCount := 0
	for _, h := range hunks {
		if h.Staged {
			stagedCount++
		} else {
			unstagedCount++
		}
	}

	if stagedCount == 0 {
		t.Error("expected at least one staged hunk")
	}
	if unstagedCount == 0 {
		t.Error("expected at least one unstaged hunk")
	}
}

func TestDiffResult_GetAllHunks(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("file1.txt", "original1\n", "initial1")
	repo.CommitFile("file2.txt", "original2\n", "initial2")

	repo.WriteFile("file1.txt", "modified1\n")
	repo.WriteFile("file2.txt", "modified2\n")

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	hunks := diff.GetAllHunks()
	if len(hunks) != 2 {
		t.Errorf("expected 2 hunks from GetAllHunks, got %d", len(hunks))
	}
}

func TestHunk_GeneratePatch(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "line1\nline2\nline3\n", "initial")
	repo.WriteFile("test.txt", "line1\nmodified\nline3\n")

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if len(diff.Files) == 0 || len(diff.Files[0].Hunks) == 0 {
		t.Fatal("expected diff with hunks")
	}

	patch := diff.Files[0].Hunks[0].GeneratePatch(&diff.Files[0])

	// Verify patch contains required parts
	if !strings.Contains(patch, "diff --git") {
		t.Error("expected patch to contain 'diff --git'")
	}
	if !strings.Contains(patch, "@@") {
		t.Error("expected patch to contain hunk header '@@'")
	}
	if !strings.Contains(patch, "-line2") {
		t.Error("expected patch to contain removed line")
	}
	if !strings.Contains(patch, "+modified") {
		t.Error("expected patch to contain added line")
	}
}

func TestGetUntrackedFileDiff(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.WriteFile("new-file.txt", "line1\nline2\nline3\n")

	diff := GetUntrackedFileDiff("new-file.txt")
	if diff == nil {
		t.Fatal("expected non-nil diff for untracked file")
	}

	if diff.Path != "new-file.txt" {
		t.Errorf("expected path 'new-file.txt', got %q", diff.Path)
	}

	if len(diff.Hunks) == 0 {
		t.Fatal("expected hunks in untracked file diff")
	}

	// All lines should be additions
	for _, hunk := range diff.Hunks {
		for _, line := range hunk.Lines {
			if line.Type == LineRemoved {
				t.Error("unexpected removed line in untracked file diff")
			}
		}
	}
}

func TestParseDiff_HunkLineNumbers(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "line1\nline2\nline3\n", "initial")
	repo.WriteFile("test.txt", "line1\nmodified\nline3\n")

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if len(diff.Files) == 0 || len(diff.Files[0].Hunks) == 0 {
		t.Fatal("expected diff with hunks")
	}

	hunk := diff.Files[0].Hunks[0]

	// Hunk should have valid line numbers
	if hunk.StartOld <= 0 {
		t.Errorf("expected positive StartOld, got %d", hunk.StartOld)
	}
	if hunk.StartNew <= 0 {
		t.Errorf("expected positive StartNew, got %d", hunk.StartNew)
	}
}

func TestParseDiff_FileHeaders(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original\n", "initial")
	repo.WriteFile("test.txt", "modified\n")

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if len(diff.Files) == 0 {
		t.Fatal("expected at least one file in diff")
	}

	headers := diff.Files[0].Header
	if len(headers) == 0 {
		t.Fatal("expected file headers")
	}

	// First header should be "diff --git"
	if !strings.HasPrefix(headers[0], "diff --git") {
		t.Errorf("expected first header to start with 'diff --git', got %q", headers[0])
	}

	// Should have --- and +++ lines
	hasMinus := false
	hasPlus := false
	for _, h := range headers {
		if strings.HasPrefix(h, "---") {
			hasMinus = true
		}
		if strings.HasPrefix(h, "+++") {
			hasPlus = true
		}
	}

	if !hasMinus {
		t.Error("expected --- header line")
	}
	if !hasPlus {
		t.Error("expected +++ header line")
	}
}

func TestDiffResult_TotalHunks(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("file1.txt", "original1\n", "initial1")
	repo.CommitFile("file2.txt", "original2\n", "initial2")

	repo.WriteFile("file1.txt", "modified1\n")
	repo.WriteFile("file2.txt", "modified2\n")

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	total := diff.TotalHunks()
	if total != 2 {
		t.Errorf("expected 2 total hunks, got %d", total)
	}
}

func TestGetDiff_DeletedFile(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "content\n", "initial")
	repo.DeleteFile("test.txt")

	// File deleted but not staged
	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	// Git shows deleted content as removed lines
	if diff.IsEmpty() {
		t.Error("expected diff for deleted file")
	}
}

func TestGetDiff_NewFileStaged(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.InitialCommit()
	repo.WriteFile("new.txt", "new content\n")
	repo.Git("add", "new.txt")

	diff, err := GetStagedDiff()
	if err != nil {
		t.Fatalf("GetStagedDiff failed: %v", err)
	}

	if diff.IsEmpty() {
		t.Error("expected staged diff for new file")
	}

	if len(diff.Files) == 0 {
		t.Fatal("expected file in diff")
	}

	// New file should show "new file mode" in headers
	hasNewFileHeader := false
	for _, h := range diff.Files[0].Header {
		if strings.Contains(h, "new file") {
			hasNewFileHeader = true
			break
		}
	}

	if !hasNewFileHeader {
		t.Error("expected 'new file' in headers for newly added file")
	}
}

func TestCombinedDiffResult_GetFileDiff(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "original\n", "initial")
	repo.WriteFile("test.txt", "staged\n")
	repo.Git("add", "test.txt")
	repo.WriteFile("test.txt", "staged\nunstaged\n")

	combined, err := GetCombinedDiff()
	if err != nil {
		t.Fatalf("GetCombinedDiff failed: %v", err)
	}

	hunks := combined.GetAllHunksCombined()
	if len(hunks) == 0 {
		t.Fatal("expected hunks")
	}

	for _, hunk := range hunks {
		fileDiff := combined.GetFileDiff(&hunk)
		if fileDiff == nil {
			t.Error("GetFileDiff returned nil for valid hunk")
		}
	}
}

func TestCombinedDiffResult_IsEmpty(t *testing.T) {
	t.Run("empty combined diff", func(t *testing.T) {
		combined := &CombinedDiffResult{
			StagedDiff:   &DiffResult{},
			UnstagedDiff: &DiffResult{},
		}
		if !combined.IsEmpty() {
			t.Error("expected empty combined diff to return true for IsEmpty")
		}
	})

	t.Run("with staged diff", func(t *testing.T) {
		combined := &CombinedDiffResult{
			StagedDiff: &DiffResult{
				Files: []FileDiff{{Path: "test.txt"}},
			},
			UnstagedDiff: &DiffResult{},
		}
		if combined.IsEmpty() {
			t.Error("expected non-empty combined diff")
		}
	})

	t.Run("with unstaged diff", func(t *testing.T) {
		combined := &CombinedDiffResult{
			StagedDiff: &DiffResult{},
			UnstagedDiff: &DiffResult{
				Files: []FileDiff{{Path: "test.txt"}},
			},
		}
		if combined.IsEmpty() {
			t.Error("expected non-empty combined diff")
		}
	})

	t.Run("nil diffs", func(t *testing.T) {
		combined := &CombinedDiffResult{
			StagedDiff:   nil,
			UnstagedDiff: nil,
		}
		if !combined.IsEmpty() {
			t.Error("expected empty for nil diffs")
		}
	})
}

func TestParseDiff_ContextLines(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	// Create a file with enough context to have context lines in diff
	content := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10\n"
	repo.CommitFile("test.txt", content, "initial")

	// Modify middle line
	modified := "line1\nline2\nline3\nline4\nMODIFIED\nline6\nline7\nline8\nline9\nline10\n"
	repo.WriteFile("test.txt", modified)

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if len(diff.Files) == 0 || len(diff.Files[0].Hunks) == 0 {
		t.Fatal("expected hunks")
	}

	hunk := diff.Files[0].Hunks[0]

	// Count line types
	contextCount := 0
	addedCount := 0
	removedCount := 0

	for _, line := range hunk.Lines {
		switch line.Type {
		case LineContext:
			contextCount++
		case LineAdded:
			addedCount++
		case LineRemoved:
			removedCount++
		}
	}

	if contextCount == 0 {
		t.Error("expected context lines in hunk")
	}
	if addedCount == 0 {
		t.Error("expected added lines in hunk")
	}
	if removedCount == 0 {
		t.Error("expected removed lines in hunk")
	}
}

func TestHunk_Metadata(t *testing.T) {
	repo := NewTestRepo(t)
	defer repo.Cleanup()

	repo.CommitFile("test.txt", "line1\nline2\nline3\n", "initial")
	repo.WriteFile("test.txt", "line1\nmodified\nline3\n")

	diff, err := GetDiff()
	if err != nil {
		t.Fatalf("GetDiff failed: %v", err)
	}

	if len(diff.Files) == 0 || len(diff.Files[0].Hunks) == 0 {
		t.Fatal("expected hunks")
	}

	hunk := diff.Files[0].Hunks[0]

	if hunk.FilePath != "test.txt" {
		t.Errorf("expected FilePath 'test.txt', got %q", hunk.FilePath)
	}

	if hunk.FileIndex != 0 {
		t.Errorf("expected FileIndex 0, got %d", hunk.FileIndex)
	}

	if hunk.HunkIndex != 0 {
		t.Errorf("expected HunkIndex 0, got %d", hunk.HunkIndex)
	}

	if !strings.HasPrefix(hunk.Header, "@@") {
		t.Errorf("expected hunk header to start with '@@', got %q", hunk.Header)
	}
}
