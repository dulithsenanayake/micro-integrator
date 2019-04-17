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
	"fmt"
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/micro-integrator/cmd/utils"
	"net/http"
)

// List APIs command related usage info
const listAPICmdLiteral = "apis"
const listAPICmdShortDesc = "List all the APIs"

var listAPICmdLongDesc = "List all the APIs\n"

var listAPICmdExamples = dedent.Dedent(`
Example:
  ` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + listAPICmdLiteral)

// apisListCmd represents the list apis command
var apisListCmd = &cobra.Command{
	Use:   listAPICmdLiteral,
	Short: listAPICmdShortDesc,
	Long:  listAPICmdLongDesc + listAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "List APIs called")
		executeListAPIsCmd()
	},
}

func init() {
	listCmd.AddCommand(apisListCmd)
}

func executeListAPIsCmd() {

	count, apis, err := GetAPIList()

	if err == nil {
		// Printing the list of available APIs
		fmt.Println("No. of APIs:", count)
		if count > 0 {
			utils.PrintList(apis)
		}
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of APIs", err)
	}
}

// GetAPIList
// @return count (no. of APIs)
// @return array of API names
// @return error
func GetAPIList() (int32, []string, error) {

	finalUrl := utils.RESTAPIBase + utils.PrefixAPIs

	utils.Logln(utils.LogPrefixInfo+"URL:", finalUrl)

	headers := make(map[string]string)

	resp, err := utils.InvokeGETRequest(finalUrl, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+finalUrl, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		apiListResponse := &utils.ListResponse{}
		unmarshalError := xml.Unmarshal([]byte(resp.Body()), &apiListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid XML response", unmarshalError)
		}
		return apiListResponse.Count, apiListResponse.List, nil
	} else {
		return 0, nil, errors.New(resp.Status())
	}
}
