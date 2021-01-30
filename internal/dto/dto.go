package dto

type SagaCreateDTO struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	StepList []*SagaStepCreateDTO `json:"step_list"`
}

type SagaStepCreateDTO struct {
	ID            string  `json:"id"`
	NextOnSuccess *string `json:"next_on_success"`
	NextOnFailure *string `json:"next_on_failure"`
	IsInitial     bool    `json:"is_initial"`
	Name          string  `json:"name"`
	Endpoint      string  `json:"endpoint"`
}

type SagaCallCreateDTO struct {
	IdempotencyKey string `json:"idempotency_key"`
	Saga           string `json:"saga"`
	Content        string
}
