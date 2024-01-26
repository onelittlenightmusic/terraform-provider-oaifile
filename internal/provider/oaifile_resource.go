package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource = &OAIFileResource{}
    _ resource.ResourceWithConfigure = &OAIFileResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewOAIFileResource() resource.Resource {
    return &OAIFileResource{}
}

// OAIFileResource is the resource implementation.
type OAIFileResource struct {
	client *OAIFileClient
}

type OAIFileResourceModel struct {
	ID          types.Int64 `tfsdk:"id"`
	FilePath 			types.String `tfsdk:"filepath"`
    FileID 			types.String `tfsdk:"file_id"`
    Name 			types.String `tfsdk:"name"`
}

// Metadata returns the resource type name.
func (r *OAIFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_file"
}

// Schema defines the schema for the resource.
func (r *OAIFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.Int64Attribute{
                Computed: true,
            },
            "filepath": schema.StringAttribute{
                Optional: true,
            },
            "file_id": schema.StringAttribute{
                Computed: true,
            },
            "name": schema.StringAttribute{
                Optional: true,
            },
        },

    }
}

// Create creates the resource and sets the initial Terraform state.
func (r *OAIFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan OAIFileResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    requestBody := OAIFileRequest{
        FilePath: plan.FilePath.ValueString(),
        Name: plan.Name.ValueString(),
    }
    // Create new oaifile
    oaifile, err := r.client.CreateOAIFile(ctx, requestBody)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating oaifile",
            "Could not create oaifile, unexpected error: "+err.Error(),
        )
        return
    }

    plan.FilePath = types.StringValue(oaifile.FilePath)
    plan.FileID = types.StringValue(oaifile.FileID)
    plan.Name = types.StringValue(oaifile.Name)
    plan.ID = types.Int64Value(int64(oaifile.ID))

    // Set state to fully populated data
    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read refreshes the Terraform state with the latest data.
func (r *OAIFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
    var state OAIFileResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
            return
    }

    // Get refreshed oaifile value from HashiCups
    oaifile, err := r.client.GetOAIFile(ctx, int(state.ID.ValueInt64()))
    if err != nil {
            resp.Diagnostics.AddError(
                    "Error Reading HashiCups ",
                    fmt.Sprintf("Could not read HashiCups oaifile ID %d: %v", state.ID.ValueInt64(), err.Error()))
            return
    }

    state.FilePath = types.StringValue(oaifile.FilePath)
    state.FileID = types.StringValue(oaifile.FileID)
    state.Name = types.StringValue(oaifile.Name)
    // state.ID = types.Int64Value(int64(oaifile.ID))

    // Set refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
            return
    }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *OAIFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Retrieve values from plan
    var plan OAIFileResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Update existing oaifile
    _, err := r.client.UpdateOAIFile(ctx, OAIFile{
        ID: int(plan.ID.ValueInt64()),
        FilePath: plan.FilePath.ValueString(),
        FileID: plan.FileID.ValueString(),
        Name: plan.Name.ValueString(),
    })
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Updating HashiCups ",
            "Could not update oaifile, unexpected error: "+err.Error(),
        )
        return
    }

    // Fetch updated items from Get as Update items are not
    // populated.
    oaifile, err := r.client.GetOAIFile(ctx, int(plan.ID.ValueInt64()))
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Reading HashiCups ",
            fmt.Sprintf("Could not read HashiCups oaifile ID %d: %v", plan.ID.ValueInt64(), err.Error()),
        )
        return
    }

    plan.FilePath = types.StringValue(oaifile.FilePath)
    plan.FileID = types.StringValue(oaifile.FileID)
    plan.Name = types.StringValue(oaifile.Name)
    plan.ID = types.Int64Value(int64(oaifile.ID))

    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *OAIFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state OAIFileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
			return
	}

	// Delete existing oaifile
	err := r.client.DeleteOAIFile(ctx, int(state.ID.ValueInt64()))
	if err != nil {
			resp.Diagnostics.AddError(
					"Error Deleting HashiCups Order",
					"Could not delete oaifile, unexpected error: "+err.Error(),
			)
			return
	}
}

func (r *OAIFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*OAIFileClient)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *OAIFileClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}