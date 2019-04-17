/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package cmd

import (
	"encoding/xml"
	"errors"
	"github.com/lithammer/dedent"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/wso2/micro-integrator/cmd/utils"
	"net/http"
	"os"
)

var proxyServiceName string

// Show ProxyService command related usage info
const showProxyServiceCmdLiteral = "proxyService"
const showProxyServiceCmdShortDesc = "Get information about the specified Proxy Service"

var showProxyServiceCmdLongDesc = "Get information about the Proxy Service specified by the flag --name, -n\n"

var showProxyServiceCmdExamples = dedent.Dedent(`
Example:
  ` + utils.ProjectName + ` ` + showCmdLiteral + ` ` + showProxyServiceCmdLiteral + ` -n TestProxyService
`)

// proxyServiceShowCmd represents the proxyService command
var proxyServiceShowCmd = &cobra.Command{
	Use:   showProxyServiceCmdLiteral,
	Short: showProxyServiceCmdShortDesc,
	Long:  showProxyServiceCmdLongDesc + showProxyServiceCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "Show ProxyService called")
		executeGetProxyServiceCmd(proxyServiceName)
	},
}

func init() {
	showCmd.AddCommand(proxyServiceShowCmd)

	proxyServiceShowCmd.Flags().StringVarP(&proxyServiceName, "name", "n", "", "Name of the Proxy Service")
	proxyServiceShowCmd.MarkFlagRequired("name")
}

func executeGetProxyServiceCmd(proxyServiceName string) {

	proxyService, err := GetProxyServiceInfo(proxyServiceName)

	if err == nil {
		// Printing the details of the Proxy Service
		printProxyServiceInfo(*proxyService)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting Information of InboundEndpoint", err)
	}
}

// GetProxyServiceInfo
// @param name of the proxy service
// @return ProxyService object
// @return error
func GetProxyServiceInfo(name string) (*utils.ProxyService, error) {

	finalUrl := utils.RESTAPIBase + utils.PrefixProxyServices + "?proxyServiceName=" + name

	utils.Logln(utils.LogPrefixInfo+"URL:", finalUrl)

	headers := make(map[string]string)

	resp, err := utils.InvokeGETRequest(finalUrl, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+finalUrl, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		proxyServiceResponse := &utils.ProxyService{}
		unmarshalError := xml.Unmarshal([]byte(resp.Body()), &proxyServiceResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid XML response", unmarshalError)
		}
		return proxyServiceResponse, nil
	} else {
		return nil, errors.New(resp.Status())
	}
}

// printProxyServiceInfo
// @param ProxyService : ProxyService object
func printProxyServiceInfo(proxyService utils.ProxyService) {
	table := tablewriter.NewWriter(os.Stdout)

	row := []string{"NAME", proxyService.Name}
	table.Append(row)

	row = []string{"DESCRIPTION", proxyService.Description}
	table.Append(row)

	row = []string{"IN SEQUENCE", proxyService.InSequence}
	table.Append(row)

	row = []string{"OUT SEQUENCE", proxyService.OutSequence}
	table.Append(row)

	row = []string{"FAULT SEQUENCE", proxyService.FaultSequence}
	table.Append(row)

	row = []string{"ENDPOINT", proxyService.Endpoint}
	table.Append(row)

	for _, transport := range proxyService.Transports {
		row = []string{"TRANSPORTS", transport}
		table.Append(row)
	}

	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	table.Render() // Send output
}
