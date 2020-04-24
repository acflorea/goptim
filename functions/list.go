package functions

// Map with functions by name
var Functions = map[string]NumericalFunction{
	"F_x_square":   F_x_square,
	"F_constant":   F_constant,
	"F_identity":   F_identity,
	"F_sombrero":   F_sombrero,
	"LIBSVM_optim": LIBSVM_optim,
	"Script":       Script,
	"SparkIt":      SparkIt,
	"K7M":          K7M,
}
