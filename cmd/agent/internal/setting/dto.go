package setting

type SetMaxWorkerDTO struct {
	MaxWorker int32 `json:"max_worker" validate:"min=0"`
}
