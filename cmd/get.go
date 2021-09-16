/*
Copyright © 2021 Meir Gabay  <unfor19@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get SSM Parameters by path",
	Long: `Get up to 1000 SSM Parameters in a single command.
Behind the scenes, parzival is making multiple API requests. For example:

parzival get --region "us-east-1" --parameters-path "/myapp/dev/" --output-file-path ".dev_parameters.json"
`,
	Run: func(cmd *cobra.Command, args []string) {
		useLocalStack, err := cmd.Flags().GetBool("localstack")
		if err != nil {
			logger.Fatalln(err)
		}
		parametersRegion, err := cmd.Flags().GetString("region")
		if err != nil {
			logger.Fatalln(err)
		}
		parametersPath, err := cmd.Flags().GetString("parameters-path")
		if err != nil {
			logger.Fatalln(err)
		}
		maxPageResults, err := cmd.Flags().GetInt32("max-page-results")
		if err != nil {
			logger.Fatalln(err)
		}
		if maxPageResults < 1 || maxPageResults > 10 {
			logger.Warnln("Invalid value for max-page-results", maxPageResults)
			maxPageResults = 10
			logger.Warnln("Using", maxPageResults, "instead")
		}
		outputFilePath, err := cmd.Flags().GetString("output-file-path")
		if err != nil {
			logger.Fatalln(err)
		}
		awsConfig := InitAwsConfig(useLocalStack, parametersRegion)
		svc := ssm.NewFromConfig(awsConfig)
		params := &ssm.GetParametersByPathInput{
			Path:           aws.String(parametersPath),
			MaxResults:     maxPageResults,
			Recursive:      true,
			WithDecryption: true,
		}

		pagniator := ssm.NewGetParametersByPathPaginator(svc, params)
		var ssmParametersGroup SsmParameterGroups
		// Iterate through the SSM Parameters pages.
		for pagniator.HasMorePages() {
			page, err := pagniator.NextPage(context.TODO())
			if err != nil {
				logger.Fatalln("failed to get a page ", err)
			}
			for _, param := range page.Parameters {
				p := &SsmParameter{
					ARN:              *param.ARN,
					Name:             *param.Name,
					Type:             string(param.Type),
					LastModifiedDate: param.LastModifiedDate.Unix(),
					Value:            *param.Value,
					Version:          int32(param.Version),
				}
				switch pType := p.Type; pType {
				case "String":
					ssmParametersGroup.String = append(ssmParametersGroup.String, *p)
				case "SecureString":
					ssmParametersGroup.SecureString = append(ssmParametersGroup.SecureString, *p)
				case "StringList":
					ssmParametersGroup.StringList = append(ssmParametersGroup.StringList, *p)
				default:
					logger.Fatalln("Unknown Parameter Type: ", p.Type, " For ", p.Name)
				}
			}
			ssmParametersJson, err := json.MarshalIndent(ssmParametersGroup, "", " ")
			if err != nil {
				fmt.Println(err)
				return
			}
			logger.Debugln("Saved file to", outputFilePath)
			_ = ioutil.WriteFile(outputFilePath, ssmParametersJson, 0644)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	now := time.Now() // used in the output file name, to avoid the loss of existing sensitive data
	sec := now.Format("20060102150405")
	getCmd.PersistentFlags().StringP("output-file-path", "o", ".parameters-"+sec+".json", "Output file path")
	getCmd.PersistentFlags().Int32P("max-page-results", "m", 10, "Max results per query, 10 is the maximum")
}
