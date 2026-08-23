package domain

import "fmt"

type ExportJob struct {
	ID     string
	Format string
	Status string
	Rows   int
	Error  string
}

func (e ExportJob) CanDownload() bool { return e.Status == "completed" && e.Rows > 0 }
func (e ExportJob) Filename(prefix string) string {
	extension := e.Format
	if extension == "" {
		extension = "json"
	}
	return fmt.Sprintf("%s-%s.%s", prefix, e.ID, extension)
}
func (e *ExportJob) Start() error {
	if e.Status != "queued" {
		return ErrConflict
	}
	e.Status = "running"
	return nil
}
func (e *ExportJob) Complete(rows int) error {
	if e.Status != "running" {
		return ErrConflict
	}
	if rows < 0 {
		return ErrConflict
	}
	e.Rows = rows
	e.Status = "completed"
	return nil
}
func (e *ExportJob) Fail(err string) error {
	if e.Status != "running" {
		return ErrConflict
	}
	e.Error = err
	e.Status = "failed"
	return nil
}
