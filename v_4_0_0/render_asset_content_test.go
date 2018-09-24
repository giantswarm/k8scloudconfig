package v_4_0_0

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTemplate = `foo: {{.Foo}}`
	testMultilineTemplate = `foo: {{.Foo}}
secondline`
)

type FakeParams struct {
	Foo string
}

func TestRenderAssetContent(t *testing.T) {
	tests := []struct {
		assetContent    string
		params          FakeParams
		expectedContent []string
	}{
		{
			assetContent:    testTemplate,
			params:          FakeParams{Foo: "bar"},
			expectedContent: []string{"foo: bar"},
		},
	}

	for _, tc := range tests {
		content, err := RenderAssetContent(tc.assetContent, tc.params)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tc.expectedContent, content, "content should be equal")
	}
}


func TestRenderFileAssetContent(t *testing.T) {
	tests := []struct {
		assetContent    string
		params          FakeParams
		expectedContent string
	}{
		{
			assetContent:    testMultilineTemplate,
			params:          FakeParams{Foo: "bar"},
			expectedContent: "foo: bar\nsecondline",
		},
	}

	for _, tc := range tests {
		content, err := RenderFileAssetContent(tc.assetContent, tc.params)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tc.expectedContent, content, "content should be equal")
	}
}