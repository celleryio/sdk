/*
 * Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"github.com/spf13/cobra"

	"github.com/cellery-io/sdk/components/cli/pkg/commands"
)

// newExtractResourcesCommand creates a command which can be invoked to extract the cell
// image resources to a specific directory.
func newExtractResourcesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract-resources [CELL_IMAGE] [OUTPUT_DIR]",
		Short: "Extract the resource files of a pulled image to the provided location",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				cmd.Help()
				return nil
			}
			err := commands.RunExtractResources(args[0], args[1])
			if err != nil {
				cmd.Help()
				return err
			}
			return nil
		},
	}
	return cmd
}
