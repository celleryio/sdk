/*
 * Copyright (c) 2018 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package constants

const DOMAIN_NAME_PATTERN = "[a-z0-9]+((-|.)[a-z0-9]+)*"
const CELLERY_ID_PATTERN = "[a-z0-9]+(-[a-z0-9]+)*"
const IMAGE_VERSION_PATTERN = "[0-9]+\\.[0-9]+\\.[0-9]+"
const CELL_IMAGE_PATTERN = CELLERY_ID_PATTERN + "\\/" + CELLERY_ID_PATTERN + ":" + IMAGE_VERSION_PATTERN
const CELL_IMAGE_WITH_REGISTRY_PATTERN = "(" + DOMAIN_NAME_PATTERN + "\\/)?" + CELL_IMAGE_PATTERN

const GROUP_NAME = "mesh.cellery.io"

const CELL_IMAGE_EXT = ".zip"

const CENTRAL_REGISTRY_HOST = "registry-1.docker.io"

const EMPTY_STRING = ""

const CELLERY_SETUP_MANAGE = "Manage"
const CELLERY_SETUP_CREATE = "Create"
const CELLERY_SETUP_SWITCH = "Switch"
const CELLERY_SETUP_BACK = "BACK"
const CELLERY_SETUP_EXIT = "EXIT"

const CELLERY_CREATE_LOCAL = "Local"
const CELLERY_CREATE_KUBEADM = "kubeadm"
const CELLERY_CREATE_GCP = "GCP"

const CELLERY_MANAGE_STOP = "stop"
const CELLERY_MANAGE_START = "start"
const CELLERY_MANAGE_CLEANUP = "cleanup"

const AWS_S3_BUCKET = "cellery-runtime-installation"
const AWS_S3_ITEM_VM_COMPLETE = "Cellery_Complete_Runtime_0.2.0.tar.gz"
const AWS_S3_ITEM_VM_MINIMAL = "Cellery_Minimal_Runtime_0.2.0.tar.gz"
const AWS_S3_ITEM_CONFIG_MINIMAL = "config-minimum-runtime.0.2.0"
const AWS_S3_ITEM_CONFIG_COMPLETE = "config-complete-runtime-0.2.0"
const AWS_REGION = "ap-south-1"

const VM_NAME = "ubuntu-18.04"
const VM_FILE_NAME = "ubuntu-18.04.ovf"
const VM_DISK_NAME = "ubuntu-18.04-disk001.vmdk"

const GCP_CLUSTER_NAME = "cellery-cluster"
const GCP_DB_INSTANCE_NAME = "cellery-sql"
const GCP_BUCKET_NAME = "cellery-gcp-bucket"
const GCP_SQL_USER_NAME = "cellery-sql-user"
const GCP_SQL_PASSWORD = "cellery-sql-user"
const GCP_NFS_SERVER_INSTANCE = "nfs-server"

const GCP_NFS_CONFIG_NAME = "data"
const GCP_NFS_CONFIG_CAPACITY = 1024
const GCP_CLUSTER_IMAGE_TYPE = "cos_containerd"
const GCP_CLUSTER_MACHINE_TYPE = "n1-standard-4"
const GCP_CLUSTER_DISK_SIZE_GB = 100
const GCP_SQL_TIER = "db-n1-standard-1"
const GCP_SQL_DISK_SIZE_GB = 20

const ZIP_BALLERINA_SOURCE = "src"
const ZIP_ARTIFACTS = "artifacts"

const CELLERY_HOME_ENV_VAR = "CELLERY_HOME"
const CELLERY_HOME_DOCS_VIEW_DIR = "docs-view"

const CELLERY_HOME = ".cellery"
const GCP = "gcp"
const ARTIFACTS = "artifacts"
const K8S_ARTIFACTS = "k8s-artefacts"
const GLOBAL_APIM = "global-apim"
const VBOX_MANAGE = "VBoxManage"
const VM = "vm"
const ARTIFACTS_OLD = "artifacts-old"
const OBSERVABILITY = "observability"
const MYSQL = "mysql"
const CONF = "conf"
const DATA_SOURCES = "datasources"
const DEPLOYMENT_YAML = "deployment.yaml"
const MASTER_DATA_SOURCES_XML = "master-datasources.xml"
const DATABASE_USERNAME = "DATABASE_USERNAME"
const DATABASE_PASSWORD = "DATABASE_PASSWORD"
const MYSQL_DATABASE_HOST = "MYSQL_DATABASE_HOST"
const SP = "sp"
const ARTIFACTS_PERSISTENT_VOLUME_YAML = "artifacts-persistent-volume.yaml"
const DB_SCRIPTS = "dbscripts"
const INIT_SQL = "init.sql"
const ERROR_REPLACING_MASTER_DATASOURCES_XML = "Error replacing in file /global-apim/conf/datasources/master-datasources.xml"
const ERROR_REPLACING_OBSERVABILITY_YAML = "Error replacing in file /observability/sp/conf/deployment.yaml"

const KUBECTL = "kubectl"
const CREATE = "create"
const APPLY = "apply"
const CONFIG_MAP = "configmap"
const KUBECTL_FLAG = "-f"

const CELLERY_HUB_KEYRING_NAME = "hubcelleryio"

const BASIC = "Basic"
const COMPLETE = "Complete"

const K8S_ARTIFACTS_PATH_MAC = "/Library/Cellery/artifacts"
const K8S_ARTIFACTS_PATH_UBUNTU = "/usr/share/cellery/artifacts"
