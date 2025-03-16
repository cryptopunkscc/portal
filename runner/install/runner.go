package install

func NewRunner(outputDir string) Runner {
	return Runner{OutputDir: outputDir}
}

type Runner struct {
	OutputDir string
}
