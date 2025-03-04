// Copyright 2022 The Howijd.Network Authors
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//   or LICENSE file in repository root.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package koios

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	// UTxO model holds inputs and outputs for given UTxO.
	UTxO struct {
		/// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash"`

		// Inputs An array with details about inputs used in a transaction.
		Inputs []TxInput `json:"inputs"`
		// Outputs An array with details about outputs from the transaction.
		Outputs []TxOutput `json:"outputs"`
	}

	// TxMetalabels defines model for tx_metalabels.
	TxMetalabel struct {
		// A distinct known metalabel
		Metalabel uint64 `json:"metalabel"`
	}

	// TxInput an transaxtion input.
	TxInput struct {
		// An array of assets contained on input UTxO.
		AssetList []Asset `json:"asset_list"`

		// input UTxO.
		PaymentAddr PaymentAddr `json:"payment_addr"`

		// StakeAddress for transaction's input UTxO.
		StakeAddress StakeAddress `json:"stake_addr,omitempty"`

		// Hash of Transaction for input UTxO.
		TxHash TxHash `json:"tx_hash"`

		// Index of input UTxO on the mentioned address used for input.
		TxIndex int `json:"tx_index"`

		// Balance on the selected input transaction.
		Value Lovelace `json:"value"`
	}

	// TxOutput an transaxtion output.
	TxOutput struct {
		// An array of assets to be included in output UTxO.
		AssetList []Asset `json:"asset_list"`

		// where funds were sent or change to be returned.
		PaymentAddr PaymentAddr `json:"payment_addr"`

		// StakeAddress for transaction's output UTxO.
		StakeAddress StakeAddress `json:"stake_addr,omitempty"`

		// Hash of this transaction.
		TxHash TxHash `json:"tx_hash"`

		// Index of output UTxO.
		TxIndex int `json:"tx_index"`

		// Total sum on the output address.
		Value Lovelace `json:"value"`
	}

	// TxInfoMetadata metadata in transaction info.
	TxInfoMetadata struct {
		// JSON containing details about metadata within transaction.
		JSON map[string]interface{} `json:"json"`

		// Key is metadata (index).
		Key int `json:"key"`
	}

	// TxsWithdrawal withdrawal record in transaction.
	TxsWithdrawal struct {
		// Amount is withdrawal amount in lovelaces.
		Amount Lovelace `json:"amount,omitempty"`
		// StakeAddress fo withdrawal.
		StakeAddress StakeAddress `json:"stake_addr,omitempty"`
	}

	// TxInfo transaction info.
	TxInfo struct {
		// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash"`

		// BlockHash is hash of the block in which transaction was included.
		BlockHash BlockHash `json:"block_hash"`

		// BlockHeight is block number on chain where transaction was included.
		BlockHeight int `json:"block_height"`

		// Epoch number.
		Epoch EpochNo `json:"epoch"`

		// EpochSlot is slot number within epoch.
		EpochSlot int `json:"epoch_slot"`

		// AbsoluteSlot is overall slot number (slots from genesis block of chain).
		AbsoluteSlot int `json:"absolute_slot"`

		// TxTimestamp is timestamp when block containing transaction was created.
		TxTimestamp string `json:"tx_timestamp"`

		// TxBlockIndex is index of transaction within block.
		TxBlockIndex int `json:"tx_block_index"`

		// TxSize is transaction size in bytes.
		TxSize int `json:"tx_size"`

		// TotalOutput is total sum of all transaction outputs (in lovelaces).
		TotalOutput Lovelace `json:"total_output"`

		// Fee is total transaction fee (in lovelaces).
		Fee Lovelace `json:"fee"`

		// Deposit is total deposits included in transaction (for example,
		// if it is registering a pool/key).
		Deposit Lovelace `json:"deposit"`

		// InvalidAfter is slot number after which transaction cannot be validated.
		InvalidAfter int `json:"invalid_after,omitempty"`

		// InvalidBefore is slot number before which transaction cannot be validated.
		// (if supplied, else 0)
		InvalidBefore int `json:"invalid_before,omitempty"`

		// Inputs An array with details about inputs used in a transaction
		Inputs []TxInput `json:"inputs"`

		// Outputs An array with details about outputs from the transaction.
		Outputs []TxOutput `json:"outputs,omitempty"`

		// AssetsMinted An array of minted assets with-in a transaction (if any).
		AssetsMinted []Asset `json:"assets_minted"`

		// Collaterals An array of collateral inputs needed when dealing with smart contracts.
		Collaterals []TxInput `json:"collaterals"`

		// Metadata present with-in a transaction (if any)
		Metadata []TxInfoMetadata `json:"metadata"`

		// Array of withdrawals with-in a transaction (if any)
		Withdrawals []TxsWithdrawal `json:"withdrawals"`

		// Certificates present with-in a transaction (if any)
		Certificates []Certificate `json:"certificates"`
	}

	// TxsInfoResponse represents response from `/tx_info` endpoint.
	TxsInfosResponse struct {
		Response
		Data []TxInfo `json:"response"`
	}

	// TxsInfoResponse represents response from `/tx_info` endpoint.
	// when requesting info about single transaction.
	TxInfoResponse struct {
		Response
		Data *TxInfo `json:"response"`
	}

	// TxUTxOsResponse represents response from `/tx_utxos` endpoint.
	TxUTxOsResponse struct {
		Response
		Data []UTxO `json:"data"`
	}

	// TxMetadata transaction metadata lookup res for `/tx_metadata` endpoint.
	TxMetadata struct {
		// TxHash is hash of transaction.
		TxHash TxHash `json:"tx_hash"`
		// Metadata present with-in a transaction (if any)
		Metadata map[string]interface{} `json:"metadata"`
	}

	// SubmitSignedTxResponse represents response from `/submittx` endpoint.
	SubmitSignedTxResponse struct {
		Response
		Data TxHash `json:"data"`
	}

	// TxBodyJSON used to Unmarshal built transactions.
	TxBodyJSON struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		CborHex     string `json:"cborHex"`
	}
	// TxMetadataResponse represents response from `/tx_metadata` endpoint.
	TxMetadataResponse struct {
		Response
		Data *TxMetadata `json:"data"`
	}

	// TxsMetadataResponse represents response from `/tx_metadata` endpoint.
	TxsMetadataResponse struct {
		Response
		Data []TxMetadata `json:"data"`
	}

	// TxMetaLabelsResponse represents response from `/tx_metalabels` endpoint.
	TxMetaLabelsResponse struct {
		Response
		Data []TxMetalabel `json:"data"`
	}

	// TxStatus is tx_status enpoint response.
	TxStatus struct {
		TxHash           TxHash `json:"tx_hash"`
		NumConfirmations uint64 `json:"num_confirmations"`
	}

	// TxsStatusesResponse represents response from `/tx_status` endpoint.
	TxsStatusesResponse struct {
		Response
		Data []TxStatus `json:"response"`
	}

	// TxStatusResponse represents response from `/tx_status` endpoint.
	TxStatusResponse struct {
		Response
		Data *TxStatus `json:"response"`
	}
)

// GetTxInfo returns detailed information about transaction.
func (c *Client) GetTxInfo(ctx context.Context, tx TxHash) (res *TxInfoResponse, err error) {
	res = &TxInfoResponse{}
	rsp, err := c.GetTxsInfos(ctx, []TxHash{tx})
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	}
	return
}

// GetTxsInfos returns detailed information about transaction(s).
func (c *Client) GetTxsInfos(ctx context.Context, txs []TxHash) (res *TxsInfosResponse, err error) {
	res = &TxsInfosResponse{}
	if len(txs) == 0 {
		err = ErrNoTxHash
		res.applyError(nil, err)
		return
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_info", txHashesPL(txs), nil, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}
	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}
	res.ready()
	return res, nil
}

// GetTxsUTxOs returns UTxO set (inputs/outputs) of transactions.
func (c *Client) GetTxsUTxOs(ctx context.Context, txs []TxHash) (res *TxUTxOsResponse, err error) {
	res = &TxUTxOsResponse{}
	if len(txs) == 0 {
		err = ErrNoTxHash
		res.applyError(nil, err)
		return
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_utxos", txHashesPL(txs), nil, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}
	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}
	res.ready()
	return res, nil
}

// GetTxMetadata returns metadata information (if any) for given transaction.
func (c *Client) GetTxMetadata(ctx context.Context, tx TxHash) (res *TxMetadataResponse, err error) {
	res = &TxMetadataResponse{}
	rsp, err := c.GetTxsMetadata(ctx, []TxHash{tx})
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	}
	return
}

// GetTxsInfos returns detailed information about transaction(s).
func (c *Client) GetTxsMetadata(ctx context.Context, txs []TxHash) (res *TxsMetadataResponse, err error) {
	res = &TxsMetadataResponse{}
	if len(txs) == 0 {
		err = ErrNoTxHash
		res.applyError(nil, err)
		return
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_metadata", txHashesPL(txs), nil, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}
	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}
	res.ready()
	return res, nil
}

// GetTxMetadataLabels retruns a list of all transaction metalabels.
func (c *Client) GetTxMetaLabels(ctx context.Context) (res *TxMetaLabelsResponse, err error) {
	res = &TxMetaLabelsResponse{}
	rsp, err := c.request(ctx, &res.Response, "GET", "/tx_metalabels", nil, nil, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}
	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}
	res.ready()
	return res, nil
}

// SubmitSignedTx Submit an transaction to the network.
func (c *Client) SubmitSignedTx(ctx context.Context, stx TxBodyJSON) (res *SubmitSignedTxResponse, err error) {
	var cborb []byte
	res = &SubmitSignedTxResponse{}

	cborb, err = hex.DecodeString(stx.CborHex)
	if err != nil {
		return
	}

	h := http.Header{}
	h.Set("Content-Type", "application/cbor")
	h.Set("Content-Length", fmt.Sprint(len(cborb)))
	rsp, err := c.request(ctx, &res.Response, "POST", "/submittx", bytes.NewBuffer(cborb), nil, h)

	if err != nil {
		res.applyError(nil, err)
		return
	}
	body, err := readResponseBody(rsp)

	if err != nil {
		res.applyError(body, err)
		return
	}

	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}

	res.ready()
	return res, nil
}

// GetTxInfo returns detailed information about transaction.
func (c *Client) GetTxStatus(ctx context.Context, tx TxHash) (res *TxStatusResponse, err error) {
	res = &TxStatusResponse{}
	rsp, err := c.GetTxsStatuses(ctx, []TxHash{tx})
	res.Response = rsp.Response
	if len(rsp.Data) == 1 {
		res.Data = &rsp.Data[0]
	}
	return
}

// GetTxsInfos returns detailed information about transaction(s).
func (c *Client) GetTxsStatuses(ctx context.Context, txs []TxHash) (res *TxsStatusesResponse, err error) {
	res = &TxsStatusesResponse{}
	if len(txs) == 0 {
		err = ErrNoTxHash
		res.applyError(nil, err)
		return
	}

	rsp, err := c.request(ctx, &res.Response, "POST", "/tx_status", txHashesPL(txs), nil, nil)
	if err != nil {
		res.applyError(nil, err)
		return
	}

	body, err := readResponseBody(rsp)
	if err != nil {
		res.applyError(body, err)
		return
	}
	if err = json.Unmarshal(body, &res.Data); err != nil {
		res.applyError(body, err)
		return
	}
	res.ready()
	return res, nil
}

func txHashesPL(txs []TxHash) io.Reader {
	var payload = struct {
		TxHashes []TxHash `json:"_tx_hashes"`
	}{txs}
	rpipe, w := io.Pipe()
	go func() {
		_ = json.NewEncoder(w).Encode(payload)
		defer w.Close()
	}()
	return rpipe
}
