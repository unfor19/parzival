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
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set SSM Parameters by path according to the output of `get`",
	Long: `Set up to 1000 SSM Parameters in a single command.
Behind the scenes, parzival is making multiple API requests. For example:
	
parzival set --region "us-east-1" --parameters-path "/myapp/stg/" --prefix-to-replace "/myapp/dev" --intput-file-path ".dev_parameters.json"
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
		inputFilePath, err := cmd.Flags().GetString("input-file-path")
		if err != nil {
			logger.Fatalln(err)
		}
		prefixToReplace, err := cmd.Flags().GetString("prefix-to-replace")
		if err != nil {
			logger.Fatalln(err)
		}
		kmsKeyId, err := cmd.Flags().GetString("kms-key-id")
		if err != nil {
			logger.Fatalln(err)
		}
		overwriteValues, err := cmd.Flags().GetBool("overwrite-values")
		if err != nil {
			logger.Fatalln(err)
		}
		awsConfig := InitAwsConfig(useLocalStack, parametersRegion)
		svc := ssm.NewFromConfig(awsConfig)
		file, _ := ioutil.ReadFile(inputFilePath)
		var data SsmParameterGroups
		_ = json.Unmarshal([]byte(file), &data)
		putParameterMiddleware(*svc, data.String, prefixToReplace, parametersPath, kmsKeyId, overwriteValues)
		putParameterMiddleware(*svc, data.SecureString, prefixToReplace, parametersPath, kmsKeyId, overwriteValues)
		putParameterMiddleware(*svc, data.StringList, prefixToReplace, parametersPath, kmsKeyId, overwriteValues)
	},
}

func putParameterMiddleware(svc ssm.Client, data []SsmParameter, prefixToReplace string, parametersPath string, kmsKeyId string, overwriteValues bool) {
	for _, p := range data {
		logger.Debug(p.Name)
		targetName := aws.String(strings.Replace(p.Name, prefixToReplace, parametersPath, -1))
		params := null
		if (p.Type == "SecureString"){
			params = ssm.PutParameterInput{
				Name:      targetName,
				Value:     &p.Value,
				KeyId:     aws.String(kmsKeyId),
				Overwrite: overwriteValues,
				Type:      types.ParameterType(p.Type),
			}			
		} else {
			params = ssm.PutParameterInput{
				Name:      targetName,
				Value:     &p.Value,
				Overwrite: overwriteValues,
				Type:      types.ParameterType(p.Type),
			}			
		}
		resp, err := svc.PutParameter(context.TODO(), &params)
		if err != nil {
			logger.Warn(err)
			logger.Fatalln("Failed to set", params.Type, params.Name)
		}
		logger.Infoln("Updated variable", *params.Name, "to version", resp.Version)
	}
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.PersistentFlags().StringP("input-file-path", "i", "", "Input file path that was generated by `get`, e.g. .parzival.json")
	setCmd.PersistentFlags().StringP("prefix-to-replace", "s", "", "Set which prefix to replace from the input file, e.g. /myapp/dev/")
	setCmd.PersistentFlags().StringP("kms-key-id", "k", "alias/aws/ssm", "The AWS KMS key to use for encrypting SecureString(s)")
	setCmd.PersistentFlags().BoolP("overwrite-values", "w", false, "Overwrite values for existing parameters")
}
