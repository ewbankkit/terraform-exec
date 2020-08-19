package tfexec

import "context"

func (tf *Terraform) WorkspaceSelect(ctx context.Context, workspace string) error {
	return tf.runTerraformCmd(tf.buildTerraformCmd(ctx, "workspace", "select", workspace))
}
