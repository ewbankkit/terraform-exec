package e2etest

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/hashicorp/terraform-exec/tfexec"
)

const defaultWorkspace = "default"

func TestWorkspaceList_default(t *testing.T) {
	runTest(t, "basic", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		workspaces, current, err := tf.WorkspaceList(context.Background())
		if err != nil {
			t.Fatalf("got error querying workspace list: %s", err)
		}
		if current != defaultWorkspace {
			t.Fatalf("expected %q workspace to be selected, got %q", defaultWorkspace, current)
		}
		if len(workspaces) != 1 || workspaces[0] != defaultWorkspace {
			t.Fatalf("expected workspace list to only contain %q, got %#v", defaultWorkspace, workspaces)
		}
	})
}

func TestWorkspaceList_multiple(t *testing.T) {
	runTest(t, "workspaces", func(t *testing.T, tfv *version.Version, tf *tfexec.Terraform) {
		workspaces, current, err := tf.WorkspaceList(context.Background())
		if err != nil {
			t.Fatalf("got error querying workspace list: %s", err)
		}
		if current != "foo" {
			t.Fatalf("expected %q workspace to be selected, got %q", "foo", current)
		}
		if !reflect.DeepEqual([]string{defaultWorkspace, "foo"}, workspaces) {
			t.Fatalf("expected %#v, got %#v", []string{defaultWorkspace, "foo"}, workspaces)
		}
	})
}
