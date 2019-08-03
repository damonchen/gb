package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// HistoryRecord history record
type HistoryRecord struct {
	ID int
	// project name
	ProjectPath string
	// git branch name
	FromName string
	// git branch name
	ToName string
	// git switch branch time
	Occur time.Time
	// from branch commit info
	FromCommit string
	// from commit info
	FromSubject string
}

func (h HistoryRecord) String() string {
	return fmt.Sprintf("%d %s", h.ID, h.Normalize())
}

func (h HistoryRecord) Normalize() string {
	t := h.Occur.Format("2006-01-02 15:04:15")
	shortCommit := h.FromCommit[:20]
	return fmt.Sprintf("%s %s %s %s -> %s", t, shortCommit, h.FromSubject, h.FromName, h.ToName)
}

func getHistoryRoot() string {
	homeDir, _ := os.UserHomeDir()
	rootDir := filepath.Join(homeDir, ".gb")
	return rootDir
}

func getHistoryFileName() string {
	rootDir := getHistoryRoot()
	historyFile := filepath.Join(rootDir, "history")
	return historyFile
}

// History history command
type History struct {
	fileName string
	db       *gorm.DB
}

// NewHistory new history
func NewHistory(fileName string) (*History, error) {
	db, err := gorm.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&HistoryRecord{})

	return &History{
		fileName: fileName,
		db:       db,
	}, nil
}

// Close close history
func (h *History) Close() error {
	if h.db != nil {
		return h.db.Close()
	}
	return nil
}

//GetAllHistoryRecord get all history
func (h *History) GetAllHistoryRecord() ([]HistoryRecord, error) {
	var historyRecords []HistoryRecord
	h.db.Order("occur desc").Where("project_path = ?", projectPath).Find(&historyRecords)
	if h.db.Error != nil {
		return nil, h.db.Error
	}

	return historyRecords, nil
}

// AddNewHistoryRecord add new history record
func (h *History) AddNewHistoryRecord(historyRecord HistoryRecord) (int, error) {
	h.db.Create(&historyRecord)
	if h.db.Error != nil {
		return 0, h.db.Error
	}

	return historyRecord.ID, nil
}

func (h *History) RemoveAllHistoryRecord() error {
	record := HistoryRecord{}
	h.db.Where("project_path = ?", projectPath).Delete(record)
	return h.db.Error
}

func (h *History) RemoveHistoryRecord(branchName string) error {
	h.db.Where("to_name = ? and project_path = ?", branchName, projectPath).Delete(&HistoryRecord{})
	return h.db.Error
}
