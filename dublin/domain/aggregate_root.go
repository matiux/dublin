package domain

// AggregateRoot represents entities that are an aggregate root.
type AggregateRoot interface {

	GetUncommittedEvents() EventStream

	GetAggregateRootId() string
}