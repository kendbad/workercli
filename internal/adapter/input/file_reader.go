package input

import (
	"bufio"
	"os"
	"strconv"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type Reader interface {
	ReadTasks(filePath string) ([]model.Task, error)
}

type FileReader struct {
	logger *utils.Logger
}

func NewFileReader(logger *utils.Logger) *FileReader {
	return &FileReader{logger: logger}
}

func (r *FileReader) ReadTasks(filePath string) ([]model.Task, error) {
	filePath = utils.AutoPath(filePath)
	r.logger.Debugf("Đọc file đầu vào: %s", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		r.logger.Errorf("Lỗi mở file %s: %v", filePath, err)
		return nil, err
	}
	defer file.Close()

	var tasks []model.Task
	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()
		tasks = append(tasks, model.Task{
			ID:   "task-" + strconv.Itoa(i),
			Data: line,
		})
	}

	if err := scanner.Err(); err != nil {
		r.logger.Errorf("Lỗi đọc file %s: %v", filePath, err)
		return nil, err
	}

	r.logger.Infof("Đọc được %d task từ file", len(tasks))
	return tasks, nil
}
