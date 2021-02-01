package entities

func NewSaga(ID, name string) *Saga {
	return &Saga{
		ID:   ID,
		Name: name,
	}
}

type Saga struct {
	// todo replace with uuid
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SagaStep struct {
	ID     string `json:"id"`
	SagaID string `json:"saga_id"`

	// NextOnSuccess is an identifier of the step called when the current one completes successfully.
	NextOnSuccess *string `json:"next_on_success"`

	// NextOnFailure is an identifier of the step called when the current one completes with an error.
	NextOnFailure *string `json:"next_on_failure"`

	// IsInitial is an attribute that holds true if the current step is first in saga or false otherwise.
	IsInitial bool `json:"is_initial"`

	Name string `json:"name"`

	// Endpoint holds an address of service should be called.
	Endpoint string `json:"endpoint"`
}

func NewSagaCall(ID, idempotencyKey, sagaID, content string) *SagaCall {
	return &SagaCall{
		ID:             ID,
		IdempotencyKey: idempotencyKey,
		SagaID:         sagaID,
		Content:        content,
	}
}

type SagaCall struct {
	ID             string
	IdempotencyKey string `json:"idempotency_key"`
	SagaID         string `json:"saga_id"`
	Content        string `json:"content"`
}

func NewUnenqueuedSagaCallStep(sagaStepID, sagaCallID, payload string) *SagaCallStep {
	return &SagaCallStep{
		SagaStepID: sagaStepID,
		SagaCallID: sagaCallID,
		Payload:    payload,
	}
}

type SagaCallStep struct {
	SagaStepID string
	SagaCallID string
	Payload    string
}
