//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azblob

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type serviceClient struct {
	con *connection
}

// FilterBlobs - The Filter Blobs operation enables callers to list blobs across all containers whose tags match a given search expression. Filter blobs
// searches across all containers within a storage account but can
// be scoped within the expression to a single container.
// If the operation fails it returns the *StorageError error type.
func (client *serviceClient) FilterBlobs(ctx context.Context, options *ServiceFilterBlobsOptions) (ServiceFilterBlobsResponse, error) {
	req, err := client.filterBlobsCreateRequest(ctx, options)
	if err != nil {
		return ServiceFilterBlobsResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return ServiceFilterBlobsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServiceFilterBlobsResponse{}, client.filterBlobsHandleError(resp)
	}
	return client.filterBlobsHandleResponse(resp)
}

// filterBlobsCreateRequest creates the FilterBlobs request.
func (client *serviceClient) filterBlobsCreateRequest(ctx context.Context, options *ServiceFilterBlobsOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.con.Endpoint())
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("comp", "blobs")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	if options != nil && options.Where != nil {
		reqQP.Set("where", *options.Where)
	}
	if options != nil && options.Marker != nil {
		reqQP.Set("marker", *options.Marker)
	}
	if options != nil && options.Maxresults != nil {
		reqQP.Set("maxresults", strconv.FormatInt(int64(*options.Maxresults), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("x-ms-version", "2019-12-12")
	if options != nil && options.RequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.RequestID)
	}
	req.Raw().Header.Set("Accept", "application/xml")
	return req, nil
}

// filterBlobsHandleResponse handles the FilterBlobs response.
func (client *serviceClient) filterBlobsHandleResponse(resp *http.Response) (ServiceFilterBlobsResponse, error) {
	result := ServiceFilterBlobsResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return ServiceFilterBlobsResponse{}, err
		}
		result.Date = &date
	}
	if err := runtime.UnmarshalAsXML(resp, &result.FilterBlobSegment); err != nil {
		return ServiceFilterBlobsResponse{}, err
	}
	return result, nil
}

// filterBlobsHandleError handles the FilterBlobs error response.
func (client *serviceClient) filterBlobsHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := StorageError{raw: string(body)}
	if err := runtime.UnmarshalAsXML(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// GetAccountInfo - Returns the sku name and account kind
// If the operation fails it returns the *StorageError error type.
func (client *serviceClient) GetAccountInfo(ctx context.Context, options *ServiceGetAccountInfoOptions) (ServiceGetAccountInfoResponse, error) {
	req, err := client.getAccountInfoCreateRequest(ctx, options)
	if err != nil {
		return ServiceGetAccountInfoResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return ServiceGetAccountInfoResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServiceGetAccountInfoResponse{}, client.getAccountInfoHandleError(resp)
	}
	return client.getAccountInfoHandleResponse(resp)
}

// getAccountInfoCreateRequest creates the GetAccountInfo request.
func (client *serviceClient) getAccountInfoCreateRequest(ctx context.Context, options *ServiceGetAccountInfoOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.con.Endpoint())
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "account")
	reqQP.Set("comp", "properties")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("x-ms-version", "2019-12-12")
	req.Raw().Header.Set("Accept", "application/xml")
	return req, nil
}

// getAccountInfoHandleResponse handles the GetAccountInfo response.
func (client *serviceClient) getAccountInfoHandleResponse(resp *http.Response) (ServiceGetAccountInfoResponse, error) {
	result := ServiceGetAccountInfoResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return ServiceGetAccountInfoResponse{}, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("x-ms-sku-name"); val != "" {
		result.SKUName = (*SKUName)(&val)
	}
	if val := resp.Header.Get("x-ms-account-kind"); val != "" {
		result.AccountKind = (*AccountKind)(&val)
	}
	if val := resp.Header.Get("x-ms-is-hns-enabled"); val != "" {
		isHierarchicalNamespaceEnabled, err := strconv.ParseBool(val)
		if err != nil {
			return ServiceGetAccountInfoResponse{}, err
		}
		result.IsHierarchicalNamespaceEnabled = &isHierarchicalNamespaceEnabled
	}
	return result, nil
}

// getAccountInfoHandleError handles the GetAccountInfo error response.
func (client *serviceClient) getAccountInfoHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := StorageError{raw: string(body)}
	if err := runtime.UnmarshalAsXML(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// GetProperties - gets the properties of a storage account's Blob service, including properties for Storage Analytics and CORS (Cross-Origin Resource Sharing)
// rules.
// If the operation fails it returns the *StorageError error type.
func (client *serviceClient) GetProperties(ctx context.Context, options *ServiceGetPropertiesOptions) (ServiceGetPropertiesResponse, error) {
	req, err := client.getPropertiesCreateRequest(ctx, options)
	if err != nil {
		return ServiceGetPropertiesResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return ServiceGetPropertiesResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServiceGetPropertiesResponse{}, client.getPropertiesHandleError(resp)
	}
	return client.getPropertiesHandleResponse(resp)
}

// getPropertiesCreateRequest creates the GetProperties request.
func (client *serviceClient) getPropertiesCreateRequest(ctx context.Context, options *ServiceGetPropertiesOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.con.Endpoint())
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "service")
	reqQP.Set("comp", "properties")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("x-ms-version", "2019-12-12")
	if options != nil && options.RequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.RequestID)
	}
	req.Raw().Header.Set("Accept", "application/xml")
	return req, nil
}

// getPropertiesHandleResponse handles the GetProperties response.
func (client *serviceClient) getPropertiesHandleResponse(resp *http.Response) (ServiceGetPropertiesResponse, error) {
	result := ServiceGetPropertiesResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if err := runtime.UnmarshalAsXML(resp, &result.StorageServiceProperties); err != nil {
		return ServiceGetPropertiesResponse{}, err
	}
	return result, nil
}

// getPropertiesHandleError handles the GetProperties error response.
func (client *serviceClient) getPropertiesHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := StorageError{raw: string(body)}
	if err := runtime.UnmarshalAsXML(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// GetStatistics - Retrieves statistics related to replication for the Blob service. It is only available on the secondary location endpoint when read-access
// geo-redundant replication is enabled for the storage account.
// If the operation fails it returns the *StorageError error type.
func (client *serviceClient) GetStatistics(ctx context.Context, options *ServiceGetStatisticsOptions) (ServiceGetStatisticsResponse, error) {
	req, err := client.getStatisticsCreateRequest(ctx, options)
	if err != nil {
		return ServiceGetStatisticsResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return ServiceGetStatisticsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServiceGetStatisticsResponse{}, client.getStatisticsHandleError(resp)
	}
	return client.getStatisticsHandleResponse(resp)
}

// getStatisticsCreateRequest creates the GetStatistics request.
func (client *serviceClient) getStatisticsCreateRequest(ctx context.Context, options *ServiceGetStatisticsOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.con.Endpoint())
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "service")
	reqQP.Set("comp", "stats")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("x-ms-version", "2019-12-12")
	if options != nil && options.RequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.RequestID)
	}
	req.Raw().Header.Set("Accept", "application/xml")
	return req, nil
}

// getStatisticsHandleResponse handles the GetStatistics response.
func (client *serviceClient) getStatisticsHandleResponse(resp *http.Response) (ServiceGetStatisticsResponse, error) {
	result := ServiceGetStatisticsResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return ServiceGetStatisticsResponse{}, err
		}
		result.Date = &date
	}
	if err := runtime.UnmarshalAsXML(resp, &result.StorageServiceStats); err != nil {
		return ServiceGetStatisticsResponse{}, err
	}
	return result, nil
}

// getStatisticsHandleError handles the GetStatistics error response.
func (client *serviceClient) getStatisticsHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := StorageError{raw: string(body)}
	if err := runtime.UnmarshalAsXML(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// GetUserDelegationKey - Retrieves a user delegation key for the Blob service. This is only a valid operation when using bearer token authentication.
// If the operation fails it returns the *StorageError error type.
func (client *serviceClient) GetUserDelegationKey(ctx context.Context, keyInfo KeyInfo, options *ServiceGetUserDelegationKeyOptions) (ServiceGetUserDelegationKeyResponse, error) {
	req, err := client.getUserDelegationKeyCreateRequest(ctx, keyInfo, options)
	if err != nil {
		return ServiceGetUserDelegationKeyResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return ServiceGetUserDelegationKeyResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServiceGetUserDelegationKeyResponse{}, client.getUserDelegationKeyHandleError(resp)
	}
	return client.getUserDelegationKeyHandleResponse(resp)
}

// getUserDelegationKeyCreateRequest creates the GetUserDelegationKey request.
func (client *serviceClient) getUserDelegationKeyCreateRequest(ctx context.Context, keyInfo KeyInfo, options *ServiceGetUserDelegationKeyOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPost, client.con.Endpoint())
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "service")
	reqQP.Set("comp", "userdelegationkey")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("x-ms-version", "2019-12-12")
	if options != nil && options.RequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.RequestID)
	}
	req.Raw().Header.Set("Accept", "application/xml")
	return req, runtime.MarshalAsXML(req, keyInfo)
}

// getUserDelegationKeyHandleResponse handles the GetUserDelegationKey response.
func (client *serviceClient) getUserDelegationKeyHandleResponse(resp *http.Response) (ServiceGetUserDelegationKeyResponse, error) {
	result := ServiceGetUserDelegationKeyResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return ServiceGetUserDelegationKeyResponse{}, err
		}
		result.Date = &date
	}
	if err := runtime.UnmarshalAsXML(resp, &result.UserDelegationKey); err != nil {
		return ServiceGetUserDelegationKeyResponse{}, err
	}
	return result, nil
}

// getUserDelegationKeyHandleError handles the GetUserDelegationKey error response.
func (client *serviceClient) getUserDelegationKeyHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := StorageError{raw: string(body)}
	if err := runtime.UnmarshalAsXML(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// ListContainersSegment - The List Containers Segment operation returns a list of the containers under the specified account
// If the operation fails it returns the *StorageError error type.
func (client *serviceClient) ListContainersSegment(options *ServiceListContainersSegmentOptions) *ServiceListContainersSegmentPager {
	return &ServiceListContainersSegmentPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listContainersSegmentCreateRequest(ctx, options)
		},
		advancer: func(ctx context.Context, resp ServiceListContainersSegmentResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.ListContainersSegmentResponse.NextMarker)
		},
	}
}

// listContainersSegmentCreateRequest creates the ListContainersSegment request.
func (client *serviceClient) listContainersSegmentCreateRequest(ctx context.Context, options *ServiceListContainersSegmentOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodGet, client.con.Endpoint())
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("comp", "list")
	if options != nil && options.Prefix != nil {
		reqQP.Set("prefix", *options.Prefix)
	}
	if options != nil && options.Marker != nil {
		reqQP.Set("marker", *options.Marker)
	}
	if options != nil && options.Maxresults != nil {
		reqQP.Set("maxresults", strconv.FormatInt(int64(*options.Maxresults), 10))
	}
	if options != nil && options.Include != nil {
		reqQP.Set("include", strings.Join(strings.Fields(strings.Trim(fmt.Sprint(options.Include), "[]")), ","))
	}
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("x-ms-version", "2019-12-12")
	if options != nil && options.RequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.RequestID)
	}
	req.Raw().Header.Set("Accept", "application/xml")
	return req, nil
}

// listContainersSegmentHandleResponse handles the ListContainersSegment response.
func (client *serviceClient) listContainersSegmentHandleResponse(resp *http.Response) (ServiceListContainersSegmentResponse, error) {
	result := ServiceListContainersSegmentResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if err := runtime.UnmarshalAsXML(resp, &result.ListContainersSegmentResponse); err != nil {
		return ServiceListContainersSegmentResponse{}, err
	}
	return result, nil
}

// listContainersSegmentHandleError handles the ListContainersSegment error response.
func (client *serviceClient) listContainersSegmentHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := StorageError{raw: string(body)}
	if err := runtime.UnmarshalAsXML(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// SetProperties - Sets properties for a storage account's Blob service endpoint, including properties for Storage Analytics and CORS (Cross-Origin Resource
// Sharing) rules
// If the operation fails it returns the *StorageError error type.
func (client *serviceClient) SetProperties(ctx context.Context, storageServiceProperties StorageServiceProperties, options *ServiceSetPropertiesOptions) (ServiceSetPropertiesResponse, error) {
	req, err := client.setPropertiesCreateRequest(ctx, storageServiceProperties, options)
	if err != nil {
		return ServiceSetPropertiesResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return ServiceSetPropertiesResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted) {
		return ServiceSetPropertiesResponse{}, client.setPropertiesHandleError(resp)
	}
	return client.setPropertiesHandleResponse(resp)
}

// setPropertiesCreateRequest creates the SetProperties request.
func (client *serviceClient) setPropertiesCreateRequest(ctx context.Context, storageServiceProperties StorageServiceProperties, options *ServiceSetPropertiesOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPut, client.con.Endpoint())
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "service")
	reqQP.Set("comp", "properties")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("x-ms-version", "2019-12-12")
	if options != nil && options.RequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.RequestID)
	}
	req.Raw().Header.Set("Accept", "application/xml")
	return req, runtime.MarshalAsXML(req, storageServiceProperties)
}

// setPropertiesHandleResponse handles the SetProperties response.
func (client *serviceClient) setPropertiesHandleResponse(resp *http.Response) (ServiceSetPropertiesResponse, error) {
	result := ServiceSetPropertiesResponse{RawResponse: resp}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

// setPropertiesHandleError handles the SetProperties error response.
func (client *serviceClient) setPropertiesHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := StorageError{raw: string(body)}
	if err := runtime.UnmarshalAsXML(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// SubmitBatch - The Batch operation allows multiple API calls to be embedded into a single HTTP request.
// If the operation fails it returns the *StorageError error type.
func (client *serviceClient) SubmitBatch(ctx context.Context, contentLength int64, multipartContentType string, body io.ReadSeekCloser, options *ServiceSubmitBatchOptions) (ServiceSubmitBatchResponse, error) {
	req, err := client.submitBatchCreateRequest(ctx, contentLength, multipartContentType, body, options)
	if err != nil {
		return ServiceSubmitBatchResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return ServiceSubmitBatchResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServiceSubmitBatchResponse{}, client.submitBatchHandleError(resp)
	}
	return client.submitBatchHandleResponse(resp)
}

// submitBatchCreateRequest creates the SubmitBatch request.
func (client *serviceClient) submitBatchCreateRequest(ctx context.Context, contentLength int64, multipartContentType string, body io.ReadSeekCloser, options *ServiceSubmitBatchOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPost, client.con.Endpoint())
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("comp", "batch")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.SkipBodyDownload()
	req.Raw().Header.Set("Content-Length", strconv.FormatInt(contentLength, 10))
	req.Raw().Header.Set("Content-Type", multipartContentType)
	req.Raw().Header.Set("x-ms-version", "2019-12-12")
	if options != nil && options.RequestID != nil {
		req.Raw().Header.Set("x-ms-client-request-id", *options.RequestID)
	}
	req.Raw().Header.Set("Accept", "application/xml")
	return req, runtime.MarshalAsXML(req, body)
}

// submitBatchHandleResponse handles the SubmitBatch response.
func (client *serviceClient) submitBatchHandleResponse(resp *http.Response) (ServiceSubmitBatchResponse, error) {
	result := ServiceSubmitBatchResponse{RawResponse: resp}
	if val := resp.Header.Get("Content-Type"); val != "" {
		result.ContentType = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	return result, nil
}

// submitBatchHandleError handles the SubmitBatch error response.
func (client *serviceClient) submitBatchHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := StorageError{raw: string(body)}
	if err := runtime.UnmarshalAsXML(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
