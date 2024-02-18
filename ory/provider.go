package ory

type Provider interface {
	Kratos() Kratos
}
