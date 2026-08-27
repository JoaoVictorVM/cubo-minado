package game

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func useTempConfigDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	original := userConfigDir
	userConfigDir = func() (string, error) { return dir, nil }
	t.Cleanup(func() { userConfigDir = original })
	return dir
}

func recordPath(t *testing.T, dir string) string {
	t.Helper()
	return filepath.Join(dir, recordDirName, recordFileName)
}

func writeRecordFile(t *testing.T, dir, contents string) {
	t.Helper()
	path := recordPath(t, dir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func assertAllNil(t *testing.T, times *BestTimes) {
	t.Helper()
	if times == nil {
		t.Fatal("BestTimes is nil")
	}
	if times.Easy != nil || times.Medium != nil || times.Hard != nil {
		t.Errorf("BestTimes = %+v, want all nil", times)
	}
}

func assertRecord(t *testing.T, got *int, want int, label string) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s is nil, want %d", label, want)
	}
	if *got != want {
		t.Errorf("%s = %d, want %d", label, *got, want)
	}
}

func TestGetBestTimesNoFile(t *testing.T) {
	useTempConfigDir(t)

	times, err := GetBestTimes()
	if err != nil {
		t.Fatalf("GetBestTimes returned error: %v", err)
	}
	assertAllNil(t, times)
}

func TestGetBestTimesCorruptFile(t *testing.T) {
	dir := useTempConfigDir(t)
	writeRecordFile(t, dir, "{ isto nao e json valido")

	times, err := GetBestTimes()
	if err != nil {
		t.Fatalf("GetBestTimes returned error: %v", err)
	}
	assertAllNil(t, times)
}

func TestSubmitTimeFirstRecord(t *testing.T) {
	useTempConfigDir(t)

	times, err := SubmitTime(DifficultyMedium, 87)
	if err != nil {
		t.Fatalf("SubmitTime returned error: %v", err)
	}
	assertRecord(t, times.Medium, 87, "Medium")
	if times.Easy != nil || times.Hard != nil {
		t.Errorf("other difficulties changed: %+v", times)
	}
}

func TestSubmitTimeImprovesRecord(t *testing.T) {
	useTempConfigDir(t)

	if _, err := SubmitTime(DifficultyHard, 200); err != nil {
		t.Fatalf("first SubmitTime returned error: %v", err)
	}

	times, err := SubmitTime(DifficultyHard, 187)
	if err != nil {
		t.Fatalf("second SubmitTime returned error: %v", err)
	}
	assertRecord(t, times.Hard, 187, "Hard")
}

func TestSubmitTimeDoesNotWorsenRecord(t *testing.T) {
	useTempConfigDir(t)

	if _, err := SubmitTime(DifficultyEasy, 42); err != nil {
		t.Fatalf("first SubmitTime returned error: %v", err)
	}

	times, err := SubmitTime(DifficultyEasy, 99)
	if err != nil {
		t.Fatalf("second SubmitTime returned error: %v", err)
	}
	assertRecord(t, times.Easy, 42, "Easy")

	persisted, err := GetBestTimes()
	if err != nil {
		t.Fatalf("GetBestTimes returned error: %v", err)
	}
	assertRecord(t, persisted.Easy, 42, "persisted Easy")
}

func TestSubmitTimeEqualTimeDoesNotOverwrite(t *testing.T) {
	useTempConfigDir(t)

	if _, err := SubmitTime(DifficultyEasy, 42); err != nil {
		t.Fatalf("first SubmitTime returned error: %v", err)
	}

	times, err := SubmitTime(DifficultyEasy, 42)
	if err != nil {
		t.Fatalf("second SubmitTime returned error: %v", err)
	}
	assertRecord(t, times.Easy, 42, "Easy")
}

func TestSubmitTimeInvalidDifficulty(t *testing.T) {
	dir := useTempConfigDir(t)

	times, err := SubmitTime(Difficulty("impossible"), 42)
	if times != nil {
		t.Errorf("SubmitTime returned times: %+v", times)
	}
	if !errors.Is(err, ErrInvalidDifficulty) {
		t.Errorf("SubmitTime returned %v, want %v", err, ErrInvalidDifficulty)
	}
	if _, statErr := os.Stat(recordPath(t, dir)); !os.IsNotExist(statErr) {
		t.Error("record file was created despite invalid difficulty")
	}
}

func TestSubmitTimeNegativeSeconds(t *testing.T) {
	dir := useTempConfigDir(t)

	times, err := SubmitTime(DifficultyEasy, -1)
	if times != nil {
		t.Errorf("SubmitTime returned times: %+v", times)
	}
	if !errors.Is(err, ErrInvalidTime) {
		t.Errorf("SubmitTime returned %v, want %v", err, ErrInvalidTime)
	}
	if _, statErr := os.Stat(recordPath(t, dir)); !os.IsNotExist(statErr) {
		t.Error("record file was created despite invalid seconds")
	}
}

func TestSubmitTimeZeroSecondsIsValid(t *testing.T) {
	useTempConfigDir(t)

	times, err := SubmitTime(DifficultyEasy, 0)
	if err != nil {
		t.Fatalf("SubmitTime returned error: %v", err)
	}
	assertRecord(t, times.Easy, 0, "Easy")
}

func TestPersistenceAcrossLoads(t *testing.T) {
	dir := useTempConfigDir(t)

	for difficulty, seconds := range map[Difficulty]int{
		DifficultyEasy:   42,
		DifficultyMedium: 87,
		DifficultyHard:   187,
	} {
		if _, err := SubmitTime(difficulty, seconds); err != nil {
			t.Fatalf("SubmitTime(%q) returned error: %v", difficulty, err)
		}
	}

	if _, err := os.Stat(recordPath(t, dir)); err != nil {
		t.Fatalf("record file was not written: %v", err)
	}

	times, err := GetBestTimes()
	if err != nil {
		t.Fatalf("GetBestTimes returned error: %v", err)
	}
	assertRecord(t, times.Easy, 42, "Easy")
	assertRecord(t, times.Medium, 87, "Medium")
	assertRecord(t, times.Hard, 187, "Hard")
}

func TestSubmitTimeLeavesNoTempFile(t *testing.T) {
	dir := useTempConfigDir(t)

	if _, err := SubmitTime(DifficultyEasy, 42); err != nil {
		t.Fatalf("SubmitTime returned error: %v", err)
	}

	if _, err := os.Stat(recordPath(t, dir) + ".tmp"); !os.IsNotExist(err) {
		t.Error("temp file was left behind after write")
	}
}

func TestGetBestTimesUnresolvableConfigDir(t *testing.T) {
	original := userConfigDir
	userConfigDir = func() (string, error) { return "", errors.New("sem diretório de configuração") }
	t.Cleanup(func() { userConfigDir = original })

	times, err := GetBestTimes()
	if err != nil {
		t.Fatalf("GetBestTimes returned error: %v", err)
	}
	assertAllNil(t, times)
}

func TestSubmitTimeSurvivesWriteFailure(t *testing.T) {
	original := userConfigDir
	userConfigDir = func() (string, error) { return "", errors.New("sem diretório de configuração") }
	t.Cleanup(func() { userConfigDir = original })

	times, err := SubmitTime(DifficultyEasy, 42)
	if err != nil {
		t.Fatalf("SubmitTime returned error on write failure: %v", err)
	}
	assertRecord(t, times.Easy, 42, "Easy")
}

func TestSubmitTimeSurvivesUnwritableRecordDir(t *testing.T) {
	dir := useTempConfigDir(t)

	blocker := filepath.Join(dir, recordDirName)
	if err := os.WriteFile(blocker, []byte("nao sou um diretório"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	times, err := SubmitTime(DifficultyHard, 187)
	if err != nil {
		t.Fatalf("SubmitTime returned error on unwritable dir: %v", err)
	}
	assertRecord(t, times.Hard, 187, "Hard")

	reloaded, err := GetBestTimes()
	if err != nil {
		t.Fatalf("GetBestTimes returned error: %v", err)
	}
	assertAllNil(t, reloaded)
}
