package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*carResource)(nil)

func NewCarResource(baseURL string, client *http.Client) resource.Resource {
	log.Printf("[DEBUG] Creating new car resource with baseURL: %s", baseURL)
	return &carResource{
		baseURL: baseURL,
		client:  client,
	}
}

type carResource struct {
	baseURL string
	client  *http.Client
}

type carResourceModel struct {
	Id    types.String `tfsdk:"id"`
	Model types.String `tfsdk:"model"`
	Year  types.Int64  `tfsdk:"year"`
}

func (r *carResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	log.Printf("[DEBUG] Setting resource type name to: %s_car", req.ProviderTypeName)
	resp.TypeName = req.ProviderTypeName + "_car"
}

func (r *carResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	log.Printf("[DEBUG] Configuring resource schema")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"model": schema.StringAttribute{
				Required: true,
			},
			"year": schema.Int64Attribute{
				Required: true,
			},
		},
	}
}

func (r *carResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	log.Printf("[DEBUG] Beginning car resource creation")
	var data carResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		log.Printf("[ERROR] Failed to read plan data: %v", resp.Diagnostics)
		return
	}

	// Create request body
	reqBody := struct {
		Model string `json:"model"`
		Year  int64  `json:"year"`
	}{
		Model: data.Model.ValueString(),
		Year:  data.Year.ValueInt64(),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal request body: %v", err)
		resp.Diagnostics.AddError(
			"Error creating car",
			fmt.Sprintf("Could not marshal request body: %s", err),
		)
		return
	}

	// Create API call
	apiReq, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/cars", r.baseURL), bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("[ERROR] Failed to create HTTP request: %v", err)
		resp.Diagnostics.AddError(
			"Error creating car",
			fmt.Sprintf("Could not create HTTP request: %s", err),
		)
		return
	}
	apiReq.Header.Set("Content-Type", "application/json")

	log.Printf("[DEBUG] Making POST request to %s", apiReq.URL.String())
	apiResp, err := r.client.Do(apiReq)
	if err != nil {
		log.Printf("[ERROR] Failed to make HTTP request: %v", err)
		resp.Diagnostics.AddError(
			"Error creating car",
			fmt.Sprintf("Could not make HTTP request: %s", err),
		)
		return
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusCreated {
		log.Printf("[ERROR] Unexpected status code: %d", apiResp.StatusCode)
		resp.Diagnostics.AddError(
			"Error creating car",
			fmt.Sprintf("Unexpected status code: %d", apiResp.StatusCode),
		)
		return
	}

	// Parse response
	var apiResponse struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(apiResp.Body).Decode(&apiResponse); err != nil {
		log.Printf("[ERROR] Failed to decode response: %v", err)
		resp.Diagnostics.AddError(
			"Error creating car",
			fmt.Sprintf("Could not decode response: %s", err),
		)
		return
	}

	data.Id = types.StringValue(apiResponse.ID)
	log.Printf("[DEBUG] Created car with ID: %s", apiResponse.ID)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *carResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	log.Printf("[DEBUG] Beginning car resource read")
	var data carResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		log.Printf("[ERROR] Failed to read state data: %v", resp.Diagnostics)
		return
	}

	// Read API call
	apiReq, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/cars/%s", r.baseURL, data.Id.ValueString()), nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create HTTP request: %v", err)
		resp.Diagnostics.AddError(
			"Error reading car",
			fmt.Sprintf("Could not create HTTP request: %s", err),
		)
		return
	}
	apiReq.Header.Set("Content-Type", "application/json")

	log.Printf("[DEBUG] Making GET request to %s", apiReq.URL.String())
	apiResp, err := r.client.Do(apiReq)
	if err != nil {
		log.Printf("[ERROR] Failed to make HTTP request: %v", err)
		resp.Diagnostics.AddError(
			"Error reading car",
			fmt.Sprintf("Could not make HTTP request: %s", err),
		)
		return
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] Car not found, removing from state")
		resp.State.RemoveResource(ctx)
		return
	}

	if apiResp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Unexpected status code: %d", apiResp.StatusCode)
		resp.Diagnostics.AddError(
			"Error reading car",
			fmt.Sprintf("Unexpected status code: %d", apiResp.StatusCode),
		)
		return
	}

	// Parse response
	var apiResponse struct {
		ID    string `json:"id"`
		Model string `json:"model"`
		Year  int64  `json:"year"`
	}
	if err := json.NewDecoder(apiResp.Body).Decode(&apiResponse); err != nil {
		log.Printf("[ERROR] Failed to decode response: %v", err)
		resp.Diagnostics.AddError(
			"Error reading car",
			fmt.Sprintf("Could not decode response: %s", err),
		)
		return
	}

	data.Id = types.StringValue(apiResponse.ID)
	data.Model = types.StringValue(apiResponse.Model)
	data.Year = types.Int64Value(apiResponse.Year)
	log.Printf("[DEBUG] Read car with ID: %s, Model: %s, Year: %d", apiResponse.ID, apiResponse.Model, apiResponse.Year)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *carResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	log.Printf("[DEBUG] Beginning car resource update")
	var data carResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		log.Printf("[ERROR] Failed to read plan data: %v", resp.Diagnostics)
		return
	}

	// Create request body
	reqBody := struct {
		Model string `json:"model"`
		Year  int64  `json:"year"`
	}{
		Model: data.Model.ValueString(),
		Year:  data.Year.ValueInt64(),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal request body: %v", err)
		resp.Diagnostics.AddError(
			"Error updating car",
			fmt.Sprintf("Could not marshal request body: %s", err),
		)
		return
	}

	// Update API call
	apiReq, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/cars/%s", r.baseURL, data.Id.ValueString()), bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("[ERROR] Failed to create HTTP request: %v", err)
		resp.Diagnostics.AddError(
			"Error updating car",
			fmt.Sprintf("Could not create HTTP request: %s", err),
		)
		return
	}
	apiReq.Header.Set("Content-Type", "application/json")

	log.Printf("[DEBUG] Making PUT request to %s", apiReq.URL.String())
	apiResp, err := r.client.Do(apiReq)
	if err != nil {
		log.Printf("[ERROR] Failed to make HTTP request: %v", err)
		resp.Diagnostics.AddError(
			"Error updating car",
			fmt.Sprintf("Could not make HTTP request: %s", err),
		)
		return
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Unexpected status code: %d", apiResp.StatusCode)
		resp.Diagnostics.AddError(
			"Error updating car",
			fmt.Sprintf("Unexpected status code: %d", apiResp.StatusCode),
		)
		return
	}

	// Parse response and update data
	var apiResponse struct {
		ID    string `json:"id"`
		Model string `json:"model"`
		Year  int64  `json:"year"`
	}
	if err := json.NewDecoder(apiResp.Body).Decode(&apiResponse); err != nil {
		log.Printf("[ERROR] Failed to decode response: %v", err)
		resp.Diagnostics.AddError(
			"Error updating car",
			fmt.Sprintf("Could not decode response: %s", err),
		)
		return
	}

	data.Id = types.StringValue(apiResponse.ID)
	data.Model = types.StringValue(apiResponse.Model)
	data.Year = types.Int64Value(apiResponse.Year)
	log.Printf("[DEBUG] Updated car with ID: %s, Model: %s, Year: %d", apiResponse.ID, apiResponse.Model, apiResponse.Year)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *carResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	log.Printf("[DEBUG] Beginning car resource deletion")
	var data carResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		log.Printf("[ERROR] Failed to read state data: %v", resp.Diagnostics)
		return
	}

	// Delete API call
	apiReq, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/cars/%s", r.baseURL, data.Id.ValueString()), nil)
	if err != nil {
		log.Printf("[ERROR] Failed to create HTTP request: %v", err)
		resp.Diagnostics.AddError(
			"Error deleting car",
			fmt.Sprintf("Could not create HTTP request: %s", err),
		)
		return
	}

	log.Printf("[DEBUG] Making DELETE request to %s", apiReq.URL.String())
	apiResp, err := r.client.Do(apiReq)
	if err != nil {
		log.Printf("[ERROR] Failed to make HTTP request: %v", err)
		resp.Diagnostics.AddError(
			"Error deleting car",
			fmt.Sprintf("Could not make HTTP request: %s", err),
		)
		return
	}
	defer apiResp.Body.Close()

	if apiResp.StatusCode != http.StatusNoContent && apiResp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Unexpected status code: %d", apiResp.StatusCode)
		resp.Diagnostics.AddError(
			"Error deleting car",
			fmt.Sprintf("Unexpected status code: %d", apiResp.StatusCode),
		)
		return
	}
	log.Printf("[DEBUG] Successfully deleted car with ID: %s", data.Id.ValueString())
}
