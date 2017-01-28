package loki

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/openzipkin/zipkin-go-opentracing/_thrift/gen-go/zipkincore"
	"github.com/pmezard/go-difflib/difflib"
)

func TestCollectorSpans(t *testing.T) {
	collector := NewCollector(5)

	// These two should be overwritten
	if err := collector.Collect(zipkincore.NewSpan()); err != nil {
		t.Fatal(err)
	}
	if err := collector.Collect(zipkincore.NewSpan()); err != nil {
		t.Fatal(err)
	}

	want := []*zipkincore.Span{}
	for i := 0; i < 5; i++ {
		span := zipkincore.NewSpan()
		span.Name = fmt.Sprintf("span %d", i)
		want = append(want, span)

		if err := collector.Collect(span); err != nil {
			t.Fatal(err)
		}
	}

	if have := collector.gather(); !reflect.DeepEqual(want, have) {
		t.Fatalf("%s", Diff(want, have))
	}

	if want, have := []*zipkincore.Span{}, collector.gather(); !reflect.DeepEqual(want, have) {
		t.Fatalf("%s", Diff(want, have))
	}
}

func TestCodec(t *testing.T) {
	want := []*zipkincore.Span{}
	for i := 0; i < 5; i++ {
		span := zipkincore.NewSpan()
		span.Name = fmt.Sprintf("span %d", i)
		span.Annotations = []*zipkincore.Annotation{}
		span.BinaryAnnotations = []*zipkincore.BinaryAnnotation{}
		want = append(want, span)
	}

	var buf bytes.Buffer
	if err := WriteSpans(want, &buf); err != nil {
		t.Fatal(err)
	}

	have, err := ReadSpans(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, have) {
		t.Fatalf("%s", Diff(want, have))
	}
}

// Diff diffs two arbitrary data structures, giving human-readable output.
func Diff(want, have interface{}) string {
	config := spew.NewDefaultConfig()
	config.ContinueOnMethod = true
	config.SortKeys = true
	config.SpewKeys = true
	text, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(config.Sdump(want)),
		B:        difflib.SplitLines(config.Sdump(have)),
		FromFile: "want",
		ToFile:   "have",
		Context:  3,
	})
	return "\n" + text
}
