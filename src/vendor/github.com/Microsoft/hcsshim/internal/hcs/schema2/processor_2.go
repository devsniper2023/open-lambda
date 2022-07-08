/*
 * HCS API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 2.5
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package hcsschema

type Processor2 struct {
	Count int32 `json:"Count,omitempty"`

	Limit int32 `json:"Limit,omitempty"`

	Weight int32 `json:"Weight,omitempty"`

	ExposeVirtualizationExtensions bool `json:"ExposeVirtualizationExtensions,omitempty"`

	// An optional object that configures the CPU Group to which a Virtual Machine is going to bind to.
	CpuGroup *CpuGroup `json:"CpuGroup,omitempty"`
}
