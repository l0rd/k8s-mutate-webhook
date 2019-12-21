package validate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	v1beta1 "k8s.io/api/admission/v1beta1"
)

func TestValidateJSON(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1",
		"request": {
			"uid": "07ea6264-e624-439b-83d0-f7724106ec16",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "PodExecOptions"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"subResource": "exec",
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "PodExecOptions"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestSubResource": "exec",
			"name": "nginx-deployment-774548f7d4-4mgc2",
			"namespace": "admission-test",
			"operation": "CONNECT",
			"userInfo": {
				"username": "system:serviceaccount:che:che-workspace",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"kind": "PodExecOptions",
				"apiVersion": "v1",
				"stdin": true,
				"stdout": true,
				"tty": true,
				"container": "nginx",
				"command": [
					"bash"
				]
			},
			"oldObject": null,
			"dryRun": false,
			"options": null
		}
	}`
	response, err := Validate([]byte(rawJSON))
	if err != nil {
		t.Errorf("failed to validate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.True(t, rr.Allowed)

}
