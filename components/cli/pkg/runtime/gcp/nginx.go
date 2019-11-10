/*
 * Copyright (c) 2019 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package gcp

import (
	"path/filepath"

	"github.com/cellery-io/sdk/components/cli/kubernetes"
	"github.com/cellery-io/sdk/components/cli/pkg/runtime"
)

func InstallNginx() error {
	for _, file := range buildNginxYamlPaths() {
		err := kubernetes.ApplyFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildNginxYamlPaths() []string {
	base := buildArtifactsPath(runtime.System)
	return []string{
		filepath.Join(base, "mandatory.yaml"),
		filepath.Join(base, "cloud-generic.yaml"),
	}
}
