# Internal Guide: Testing the In-File Storage Layer in Go

This README is a practical, internal learning guide to help you test a file-backed storage implementation (infile storage) for a Todo app. It builds on your existing tests and shows best practices to make them reliable, deterministic, and maintainable.

Use this as a reference while iterating on `dal/infiletodo` tests.

---

## Objectives

- Understand what to test in a file-backed store and why
- Write deterministic tests that don’t interfere with each other
- Validate on-disk data (CSV/JSON) safely and meaningfully
- Cover persistence across restarts and corruption handling
- Add concurrency tests and run with the race detector

---

## Where the Interface Should Live

- Put the interface (e.g., `TodoStore`) in the higher-level package that uses it (typically `domain`).
- The infile store (concrete implementation) lives in `dal/infiletodo` and satisfies `domain.TodoStore`.
- Your service layer depends on `domain.TodoStore`, not on DAL packages.

This avoids circular imports and keeps responsibilities clear.

---

## What to Test for a File-Backed Store

- CRUD correctness:
  - List: empty on a fresh file
  - Create: unique IDs, sets timestamps, persists to disk
  - Get: found vs not found
  - Update: mutates fields, preserves identity and CreatedAt
  - Delete: removes and persists
- Persistence:
  - Data survives reopening the store (simulate process restart)
- Concurrency:
  - Parallel Create calls produce unique IDs; no data races or corruption
- File-level errors:
  - Corrupted file content yields meaningful errors (no panics)
- Determinism:
  - Tests don’t rely on current time or hard-coded temp file names
  - No shared state across parallel tests

---

## Test Setup Best Practices

- Use `t.TempDir()` for isolation. Don’t write into `os.TempDir()` with fixed names; parallel tests will collide.
- Prefer black-box tests via the public API. Testing unexported helpers (like `saveToFile` and `loadFromFile`) is okay if needed, but the bulk should exercise the public behavior.
- Use `testify/require` (or `assert`) for concise checks and better diffs.
- Avoid direct `==` comparison on `time.Time`. Use `require.WithinDuration(...)`.
- For file content checks:
  - Prefer parsing the file back into structs and comparing structs
  - If you must check raw bytes, verify prefixes or structured patterns, not entire timestamps

---

## Common Issues in the Initial Snippets (and Fixes)

- Using `os.TempDir()` with fixed filenames in parallel tests causes interference. Fix: `dir := t.TempDir(); file := filepath.Join(dir, "todos.csv")`.
- Missing imports like `fmt` (used for `Sprintf`).
- Non-deterministic time comparisons: `time.Now()` used in expected values will never match. Fix: assert with `WithinDuration` or compare only non-time fields when needed.
- Test name typo `TestList§`: should be `TestList`.
- Hard-coding ID `0` on Create expectations: if your store auto-increments, assert positivity and uniqueness instead of a specific ID value.
- Using `bytes.Contains` for whole-line matching with timestamps might be flaky; parse lines instead or check a stable prefix.

---

## Minimal Store Test Skeleton (Black-Box Style)

A compact template you can adapt across CRUD tests. Note the use of `t.TempDir()` and `require`.

```go
    package infiletodo_test

    import (
      "context"
      "path/filepath"
      "testing"
      "time"

      "github.com/stretchr/testify/require"

      "github.com/macesz/todo-go/domain"
      "github.com/macesz/todo-go/dal/infiletodo"
    )

    // Adjust this constructor helper to your actual API.
    func newStore(t *testing.T) domain.TodoStore {
      t.Helper()
      dir := t.TempDir()
      path := filepath.Join(dir, "todos.csv")
      s, err := infiletodo.New(path)
      require.NoError(t, err)
      return s
    }

    func TestList_EmptyOnFreshFile(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      got, err := s.List(context.Background())
      require.NoError(t, err)
      require.Empty(t, got)
    }

    func TestCreate_PersistsAndReturnsTodo(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      start := time.Now()
      todo, err := s.Create(context.Background(), "Write tests")
      require.NoError(t, err)

      require.Greater(t, todo.ID, 0)
      require.Equal(t, "Write tests", todo.Title)
      require.False(t, todo.Done)
      require.WithinDuration(t, start, todo.CreatedAt, 2*time.Second)

      // Verify via List (reads from file)
      list, err := s.List(context.Background())
      require.NoError(t, err)
      require.Len(t, list, 1)
      require.Equal(t, todo.ID, list[0].ID)
      require.Equal(t, todo.Title, list[0].Title)
      require.Equal(t, todo.Done, list[0].Done)
    }

---
```
## Testing Private Helpers (`saveToFile` / `loadFromFile`)

It’s okay to test them when they encode critical file logic. Keep them in the same package (`package infiletodo`) to access unexported methods.

Key improvements:
- Use `t.TempDir()` per test
- Parse the file content to validate structure instead of `bytes.Contains`
- If you must check raw bytes, verify only the stable parts (ID, Title, Done) and then separately assert CreatedAt parses as RFC3339

Example: validating `saveToFile` by parsing CSV
```go

    package infiletodo

    import (
      "bufio"
      "context"
      "encoding/csv"
      "fmt"
      "os"
      "path/filepath"
      "testing"
      "time"

      "github.com/stretchr/testify/require"
      "github.com/macesz/todo-go/domain"
    )

    func Test_saveToFile_WritesExpectedCSV(t *testing.T) {
      t.Parallel()
      dir := t.TempDir()
      file := filepath.Join(dir, "todos.csv")

      ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
      s := &InFileStore{
        filePath: file,
        data: map[int]domain.Todo{
          1: {ID: 1, Title: "A", Done: false, CreatedAt: ts},
          2: {ID: 2, Title: "B", Done: true, CreatedAt: ts},
        },
      }

      err := s.saveToFile()
      require.NoError(t, err)

      f, err := os.Open(file)
      require.NoError(t, err)
      defer f.Close()

      r := csv.NewReader(bufio.NewReader(f))
      rows, err := r.ReadAll()
      require.NoError(t, err)
      require.Len(t, rows, 2)

      // Example row format: ID,Title,Done,CreatedAt(RFC3339)
      for _, row := range rows {
        require.Len(t, row, 4)
        // Validate CreatedAt parses
        _, perr := time.Parse(time.RFC3339, row[3])
        require.NoError(t, perr, fmt.Sprintf("invalid timestamp: %q", row[3]))
      }
    }

Example: validating `loadFromFile` from a prepared file

    func Test_loadFromFile_ReadsCSVIntoMap(t *testing.T) {
      t.Parallel()
      dir := t.TempDir()
      file := filepath.Join(dir, "todos.csv")

      // Prepare file content
      ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC).Format(time.RFC3339)
      content := "1,Todo 1,false," + ts + "\n" + "2,Todo 2,true," + ts + "\n"
      require.NoError(t, os.WriteFile(file, []byte(content), 0o600))

      s := &InFileStore{filePath: file, data: make(map[int]domain.Todo)}
      err := s.loadFromFile()
      require.NoError(t, err)

      require.Len(t, s.data, 2)
      require.Equal(t, "Todo 1", s.data[1].Title)
      require.Equal(t, false, s.data[1].Done)
      require.Equal(t, "Todo 2", s.data[2].Title)
      require.Equal(t, true, s.data[2].Done)
    }

```
---

## CRUD Through Public API (Preferred)

These are higher-value tests. They ensure end-to-end behavior including file I/O.
```go

Create

    func TestCreate(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      got, err := s.Create(context.Background(), "Test Todo")
      require.NoError(t, err)
      require.Greater(t, got.ID, 0)
      require.Equal(t, "Test Todo", got.Title)
      require.False(t, got.Done)
    }

List

    func TestList(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      _, _ = s.Create(context.Background(), "Todo 1")
      _, _ = s.Create(context.Background(), "Todo 2")

      got, err := s.List(context.Background())
      require.NoError(t, err)
      require.Len(t, got, 2)
    }

Get (found and not found)

    func TestGet(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      created, _ := s.Create(context.Background(), "X")

      got, err := s.Get(context.Background(), created.ID)
      require.NoError(t, err)
      require.Equal(t, created.ID, got.ID)

      _, err = s.Get(context.Background(), 999_999)
      require.Error(t, err) // Prefer require.ErrorIs if you export ErrNotFound
    }

Update

    func TestUpdate(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      created, _ := s.Create(context.Background(), "Old")
      updated, err := s.Update(context.Background(), created.ID, "New", true)
      require.NoError(t, err)

      require.Equal(t, created.ID, updated.ID)
      require.Equal(t, "New", updated.Title)
      require.True(t, updated.Done)
      // CreatedAt should be preserved
      require.True(t, updated.CreatedAt.Equal(created.CreatedAt))
    }

Delete

    func TestDelete(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      created, _ := s.Create(context.Background(), "A")
      require.NoError(t, s.Delete(context.Background(), created.ID))

      _, err := s.Get(context.Background(), created.ID)
      require.Error(t, err)

      // Deleting again should return not found
      err = s.Delete(context.Background(), created.ID)
      require.Error(t, err)
    }

```
---


## Persistence Across Reopen

Simulates a process restart by constructing a new store pointing to the same file.

```go
    func TestPersistsAcrossReopen(t *testing.T) {
      t.Parallel()
      dir := t.TempDir()
      path := filepath.Join(dir, "todos.csv")

      s1, err := infiletodo.New(path)
      require.NoError(t, err)
      created, _ := s1.Create(context.Background(), "Persist me")

      s2, err := infiletodo.New(path)
      require.NoError(t, err)
      got, err := s2.Get(context.Background(), created.ID)
      require.NoError(t, err)
      require.Equal(t, created.ID, got.ID)
    }

```
---

## Corrupted File Handling

Write garbage into the data file and ensure operations fail with a meaningful error.
```go
    func TestCorruptedFile_ReturnsError(t *testing.T) {
      t.Parallel()
      dir := t.TempDir()
      path := filepath.Join(dir, "todos.csv")

      require.NoError(t, os.WriteFile(path, []byte("{not-csv"), 0o600))

      s, err := infiletodo.New(path)
      require.NoError(t, err) // if constructor defers reading
      _, err = s.List(context.Background())
      require.Error(t, err, "expected error due to corrupted file")
    }

```
If your constructor reads immediately, change the expectation to `require.Error(t, err)` at `New`.

---

## Concurrency: Unique IDs and Race Safety

Run with the race detector: `go test -race ./...`.
```go
    func TestConcurrentCreateUniqueIDs(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      const n = 200
      ids := make(chan int, n)
      var wg sync.WaitGroup
      // wg.Add(n)

      for i := 0; i < n; i++ {
      	wg.Add(1)

        go func() {
          defer wg.Done()

          todo, err := s.Create(context.Background(), "T")
          require.NoError(t, err)
          ids <- todo.ID
        }()
      }

	      wg.Wait()
      close(ids)

      seen := map[int]struct{}{}
      for id := range ids {
        if _, ok := seen[id]; ok {
          t.Fatalf("duplicate id: %d", id)
        }
        seen[id] = struct{}{}
      }
    }
```
---

## Context Handling (Optional)

If the store checks `ctx` (for example during disk I/O), validate that a canceled context results in a prompt error.

```go
    func TestRespectsContextCancellation(t *testing.T) {
      t.Parallel()
      s := newStore(t)

      ctx, cancel := context.WithCancel(context.Background())
      cancel()

      _, err := s.List(ctx)
      // If your store propagates context errors:
      // require.ErrorIs(t, err, context.Canceled)
      // Otherwise, skip this test.
      require.Error(t, err)
    }

---

## Running Tests

- All tests: `go test ./...`
- Verbose: `go test -v ./...`
- Without cache (iteration): `go test -count=1 ./...`
- Race detector: `go test -race ./...`
- Shuffle (catches hidden coupling): `go test -shuffle=on ./...`
- Coverage report:
  - `go test -coverprofile=cover.out ./... && go tool cover -func=cover.out`
  - `go tool cover -html=cover.out`

---

## Checklist Before Submitting

- Uses `t.TempDir()`; no shared files between tests
- Doesn’t rely on exact timestamps; uses `WithinDuration` or parses RFC3339
- CRUD happy and error paths covered
- Persistence and corruption cases covered
- Concurrency test included; passes `-race`
- Tests are black-box where possible; helper tests limited and focused
- Clear assertions with `require`/`assert`

---

## Appendix: Upgrading Your Current Tests

- Replace `filepath.Join(os.TempDir(), "...")` with per-test paths via `t.TempDir()`.
- Import `fmt` if you use `fmt.Sprintf`.
- Fix `TestList§` to `TestList`.
- When asserting `Create`, don’t expect `ID == 0`; assert the properties you control: non-empty title, `Done == false`, `ID > 0`.
- For file content checks, prefer parsing CSV and validating each field instead of using `bytes.Contains` on a full line with a dynamic timestamp.
