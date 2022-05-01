package experiment

type Experiment struct {
	Parameter string `yaml:"parameter"`
	Value     string `yaml:"value"`
	Key       string `yaml:"key"`
	Applied   bool
}

//func(e *Experiment) Execute(kubeletConfig map[string]interface{}) error {
//	val := kubeletConfig[e.Parameter]
//}