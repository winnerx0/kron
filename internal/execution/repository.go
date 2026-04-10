package execution

type Repository interface {
	Save(execution Execution) error
	FindByJobID(jobID uint) ([]Execution, error)
	
}