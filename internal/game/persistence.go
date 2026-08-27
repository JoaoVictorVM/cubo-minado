package game

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const (
	recordDirName  = "cubo-minado"
	recordFileName = "highscores.json"
)

var userConfigDir = os.UserConfigDir

func recordFilePath() (string, error) {
	dir, err := userConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, recordDirName, recordFileName), nil
}

func loadBestTimes() *BestTimes {
	empty := &BestTimes{}

	path, err := recordFilePath()
	if err != nil {
		return empty
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return empty
	}

	var times BestTimes
	if err := json.Unmarshal(data, &times); err != nil {
		return empty
	}

	return &times
}

func saveBestTimes(times *BestTimes) error {
	path, err := recordFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(times, "", "  ")
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}

	return nil
}

func GetBestTimes() (*BestTimes, error) {
	return loadBestTimes(), nil
}

func SubmitTime(difficulty Difficulty, seconds int) (*BestTimes, error) {
	if _, ok := BoardSizeForDifficulty(difficulty); !ok {
		return nil, ErrInvalidDifficulty
	}
	if seconds < 0 {
		return nil, ErrInvalidTime
	}

	times := loadBestTimes()

	current := recordFor(times, difficulty)
	if current != nil && *current <= seconds {
		return times, nil
	}

	setRecord(times, difficulty, seconds)

	if err := saveBestTimes(times); err != nil {
		log.Printf("cubo-minado: falha ao gravar o recorde de %s: %v", difficulty, err)
	}

	return times, nil
}

func recordFor(times *BestTimes, difficulty Difficulty) *int {
	switch difficulty {
	case DifficultyEasy:
		return times.Easy
	case DifficultyMedium:
		return times.Medium
	case DifficultyHard:
		return times.Hard
	default:
		return nil
	}
}

func setRecord(times *BestTimes, difficulty Difficulty, seconds int) {
	switch difficulty {
	case DifficultyEasy:
		times.Easy = &seconds
	case DifficultyMedium:
		times.Medium = &seconds
	case DifficultyHard:
		times.Hard = &seconds
	}
}
