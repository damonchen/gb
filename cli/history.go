package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// ProjectRecord project record
type ProjectRecord struct {
	ID int
	// project name
	ProjectPath string `gorm:"type:varchar(512);unique_index"`
}

// BranchSwitchRecord history record
type BranchSwitchRecord struct {
	ID int
	// project name
	ProjectID int
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

func (h BranchSwitchRecord) String() string {
	return fmt.Sprintf("%d %s", h.ID, h.Normalize())
}

func (h BranchSwitchRecord) Normalize() string {
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
	fileName      string
	projectPath   string
	projectRecord *ProjectRecord
	db            *gorm.DB
}

// NewHistory new history
func NewHistory(projectPath string, fileName string) (*History, error) {
	db, err := gorm.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&BranchSwitchRecord{})
	db.AutoMigrate(&ProjectRecord{})

	h := &History{
		fileName:    fileName,
		projectPath: projectPath,
		db:          db,
	}
	projectRecord, err := h.GetOrCreateProjectRecord()
	if err != nil {
		return nil, err
	}

	h.projectRecord = projectRecord
	return h, nil
}

func (h *History) GetOrCreateProjectRecord() (*ProjectRecord, error) {
	projectRecord := ProjectRecord{
		ProjectPath: h.projectPath,
	}
	h.db.Where("project_path = ?", h.projectPath).First(&projectRecord)
	if h.db.Error != nil {
		if gorm.IsRecordNotFoundError(h.db.Error) {
			projectRecord.ProjectPath = h.projectPath
			h.db.Create(&projectRecord)
		}

		if h.db.Error != nil {
			return nil, h.db.Error
		}

		return &projectRecord, h.db.Error
	}

	return &projectRecord, nil
}

// Close close history
func (h *History) Close() error {
	if h.db != nil {
		return h.db.Close()
	}
	return nil
}

//GetProjectBranchSwitchRecords get all history
func (h *History) GetProjectBranchSwitchRecords() ([]BranchSwitchRecord, error) {
	var historyRecords []BranchSwitchRecord
	h.db.Order("occur desc").Where("project_id = ?", h.projectRecord.ID).Find(&historyRecords)
	if h.db.Error != nil {
		return nil, h.db.Error
	}

	return historyRecords, nil
}

// AddNewProjectBranchSwitchRecord add new history record
func (h *History) AddNewProjectBranchSwitchRecord(historyRecord BranchSwitchRecord) (int, error) {
	historyRecord.ProjectID = h.projectRecord.ID
	h.db.Create(&historyRecord)
	if h.db.Error != nil {
		return 0, h.db.Error
	}

	return historyRecord.ID, nil
}

// RemoveProjectBranchSwitchRecords remove project branch records
func (h *History) RemoveProjectBranchSwitchRecords() error {
	record := BranchSwitchRecord{}
	h.db.Where("project_id = ?", h.projectRecord.ID).Delete(record)
	return h.db.Error
}

// RemoveProjectBranchSwitchRecord remove one branch of project
func (h *History) RemoveProjectBranchSwitchRecord(branchName string) error {
	h.db.Where("to_name = ? and project_id = ?", branchName, h.projectRecord.ID).Delete(&BranchSwitchRecord{})
	return h.db.Error
}
