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

// List Seqeunces command related usage info
const listSequenceCmdLiteral = "sequences"
const listSequenceCmdShortDesc = "List all the Sequences"

var listSequenceCmdLongDesc = "List all the Sequences"

var listSequenceCmdExamples = dedent.Dedent(`
Example:
  ` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + listSequenceCmdLiteral)

// sequencesListCmd represents the list sequences command
var sequencesListCmd = &cobra.Command{
	Use:   listSequenceCmdLiteral,
	Short: listSequenceCmdShortDesc,
	Long:  listSequenceCmdLongDesc + listSequenceCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "List sequences called")
		executeListSequencesCmd()
	},
}

func init() {
	listCmd.AddCommand(sequencesListCmd)
}

func executeListSequencesCmd() {

	count, sequences, err := GetSequenceList()

	if err == nil {
		// Printing the list of available Sequences
		fmt.Println("No. of Sequences:", count)
		if count > 0 {
			utils.PrintList(sequences)
		}
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of Sequences", err)
		fmt.Println("Something went wrong", err)
	}
}

// GetSequenceList
// @return count (no. of Sequences)
// @return array of Sequence Names
// @return error
func GetSequenceList() (int32, []string, error) {

	finalUrl := utils.RESTAPIBase + utils.PrefixSequences

	utils.Logln(utils.LogPrefixInfo+"URL:", finalUrl)

	headers := make(map[string]string)

	resp, err := utils.InvokeGETRequest(finalUrl, headers)

	if err != nil {
		utils.HandleErrorAndExit("Unable to connect to "+finalUrl, err)
	}

	utils.Logln(utils.LogPrefixInfo+"Response:", resp.Status())

	if resp.StatusCode() == http.StatusOK {
		serviceListResponse := &utils.ListResponse{}
		unmarshalError := xml.Unmarshal([]byte(resp.Body()), &serviceListResponse)

		if unmarshalError != nil {
			utils.HandleErrorAndExit(utils.LogPrefixError+"invalid XML response", unmarshalError)
		}
		return serviceListResponse.Count, serviceListResponse.List, nil
	} else {
		return 0, nil, errors.New(resp.Status())
	}
}
