package helpers

import (
	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"
)

// CheckHTML uses machine learning to detect fake
// websites such as pishing ones
func CheckHTML(html string) float32 {
	model := tg.LoadModel("models/phishing", []string{"serve"}, nil)

	input, _ := tf.NewTensor(html)

	results := model.Exec([]tf.Output{
		model.Op("StatefulPartitionedCall", 0),
	}, map[tf.Output]*tf.Tensor{
		model.Op("serving_default_inputs_input", 0): input,
	})

	return results[0].Value().(float32)
}
