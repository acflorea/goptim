package generators

// The distributions
type Distribution int

// Map with distribution by name
var Distributions = map[string]Distribution{
	"Uniform":     Uniform,
	"Exponential": Exponential,
	"Discrete":    Discrete,
}

// Types of distributions
const (
	// Uniform
	Uniform Distribution = iota
	// Exponential
	Exponential
	// Discrete
	Discrete
)

// Algorithm labels
var distributions = [...]string{
	"Uiform",
	"Exponential",
	"Discrete",
}

// The random generation algorithm
type Algorithm int

// Map with functions by name
var Algorithms = map[string]Algorithm{
	"ManagerWorker":   ManagerWorker,
	"Leapfrog":        Leapfrog,
	"SeqSplit":        SeqSplit,
	"Parametrization": Parametrization,
}

// Types of parallel random generators
const (
	// A single generator, the master generates the values and pushes them to workers
	ManagerWorker Algorithm = iota
	// For each sequence of random numbers x[r]...x[r+p]...x[r+2p]... the process p takes every p-th value
	Leapfrog
	// Block allocation of data to tasks
	SeqSplit
	// Each woker has it's own parametrized generator
	Parametrization
)

// Algorithm labels
var algorithms = [...]string{
	"ManagerWorker",
	"FebrLeapfroguary",
	"SeqSplit",
	"Parametrization",
}
