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

package org.cellery.components.test.scenarios.reviews;

import io.cellery.models.Cell;
import org.ballerinax.kubernetes.exceptions.KubernetesPluginException;
import org.ballerinax.kubernetes.utils.KubernetesUtils;
import org.cellery.components.test.models.CellImageInfo;
import org.cellery.components.test.utils.CelleryUtils;
import org.cellery.components.test.utils.LangTestUtils;
import org.testng.Assert;
import org.testng.annotations.AfterClass;
import org.testng.annotations.BeforeClass;
import org.testng.annotations.Test;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;

import static org.cellery.components.test.utils.CelleryTestConstants.BAL;
import static org.cellery.components.test.utils.CelleryTestConstants.CELLERY;
import static org.cellery.components.test.utils.CelleryTestConstants.CELLERY_IMAGE_NAME;
import static org.cellery.components.test.utils.CelleryTestConstants.CELLERY_IMAGE_ORG;
import static org.cellery.components.test.utils.CelleryTestConstants.CELLERY_IMAGE_VERSION;
import static org.cellery.components.test.utils.CelleryTestConstants.PRODUCT_REVIEW;
import static org.cellery.components.test.utils.CelleryTestConstants.TARGET;
import static org.cellery.components.test.utils.CelleryTestConstants.YAML;

public class CustomerProductTest {

    private static final Path SAMPLE_DIR = Paths.get(System.getProperty("sample.dir"));
    private static final Path SOURCE_DIR_PATH =
            SAMPLE_DIR.resolve(PRODUCT_REVIEW + File.separator + CELLERY + File.separator +
                    "customer-products");
    private static final Path TARGET_PATH = SOURCE_DIR_PATH.resolve(TARGET);
    private static final Path CELLERY_PATH = TARGET_PATH.resolve(CELLERY);
    private Cell cell;
    private CellImageInfo cellImageInfo = new CellImageInfo("myorg", "products", "1.0.0", "cust-inst");
    private Map<String, CellImageInfo> dependencyCells = new HashMap<>();

    @BeforeClass
    public void compileSample() throws IOException, InterruptedException {
        Assert.assertEquals(LangTestUtils.compileCellBuildFunction(SOURCE_DIR_PATH, "customer-products" + BAL
                , cellImageInfo), 0);
        Assert.assertEquals(LangTestUtils.compileCellRunFunction(SOURCE_DIR_PATH, "customer-products" + BAL
                , cellImageInfo, dependencyCells), 0);
        File artifactYaml = CELLERY_PATH.resolve(cellImageInfo.getName() + YAML).toFile();
        Assert.assertTrue(artifactYaml.exists());
        cell = CelleryUtils.getInstance(CELLERY_PATH.resolve(cellImageInfo.getName() + YAML).toString());
    }

    @Test
    public void validateCellAvailability() {
        Assert.assertNotNull(cell);
    }

    @Test
    public void validateAPIVersion() {
        Assert.assertEquals(cell.getApiVersion(), "mesh.cellery.io/v1alpha1");
    }

    @Test
    public void validateMetaData() {
        Assert.assertEquals(cell.getMetadata().getName(), cellImageInfo.getName());
        Assert.assertEquals(cell.getMetadata().getAnnotations().get(CELLERY_IMAGE_ORG),
                cellImageInfo.getOrg());
        Assert.assertEquals(cell.getMetadata().getAnnotations().get(CELLERY_IMAGE_NAME),
                cellImageInfo.getName());
        Assert.assertEquals(cell.getMetadata().getAnnotations().get(CELLERY_IMAGE_VERSION),
                cellImageInfo.getVer());
    }

    @Test
    public void validateGatewayTemplate() {

        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getGrpc().get(0).getBackendHost(),
                "categories");
        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getGrpc().get(0).getBackendPort(),
                8000);
        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getGrpc().get(0).getPort(),
                8000);

        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(0).getBackend(),
                "customers");
        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(0).getContext(),
                "customers-1");
        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(0).getDefinitions().get(0).
                getMethod(), "GET");
        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(0).getDefinitions().get(0).
                getPath(), "/*");
        Assert.assertTrue(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(0).isAuthenticate());

        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(1).getBackend(),
                "products");
        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(1).getContext(),
                "products-1");
        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(1).getDefinitions().get(0).
                getMethod(), "GET");
        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(1).getDefinitions().get(0).
                getPath(), "/*");
        Assert.assertTrue(cell.getSpec().getGatewayTemplate().getSpec().getHttp().get(0).isAuthenticate());

        Assert.assertEquals(cell.getSpec().getGatewayTemplate().getSpec().getType(), "MicroGateway");
    }

    @Test
    public void validateServiceTemplates() {
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(0).getMetadata().getName(), "customers");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(0).getSpec().getContainer().getEnv().get(0).
                getName(), "PORT");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(0).getSpec().getContainer().getEnv().get(0).
                getValue(), "8080");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(0).getSpec().getContainer().getImage(),
                "celleryio/samples-productreview-customers");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(0).getSpec().getContainer().getPorts().get(0).
                getContainerPort().intValue(), 8080);
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(0).getSpec().getReplicas(), 1);
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(0).getSpec().getServicePort(), 80);

        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(1).getMetadata().getName(), "categories");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(1).getSpec().getContainer().getEnv().get(0).
                getName(), "PORT");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(1).getSpec().getContainer().getEnv().get(0).
                getValue(), "8000");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(1).getSpec().getContainer().getImage(),
                "celleryio/samples-productreview-categories");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(1).getSpec().getContainer().getPorts().get(0).
                getContainerPort().intValue(), 8000);
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(1).getSpec().getReplicas(), 1);
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(1).getSpec().getServicePort(), 8000);
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(1).getSpec().getProtocol(), "GRPC");

        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getMetadata().getName(), "products");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getContainer().getEnv().get(0).
                getName(), "CATEGORIES_PORT");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getContainer().getEnv().get(0).
                getValue(), "8000");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getContainer().getEnv().get(1).
                getName(), "PORT");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getContainer().getEnv().get(1).
                getValue(), "8080");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getContainer().getEnv().get(2).
                getName(), "CATEGORIES_HOST");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getContainer().getEnv().get(2).
                getValue(), "");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getContainer().getImage(),
                "celleryio/samples-productreview-products");
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getContainer().getPorts().get(0).
                getContainerPort().intValue(), 8080);
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getReplicas(), 1);
        Assert.assertEquals(cell.getSpec().getServicesTemplates().get(2).getSpec().getServicePort(), 80);
    }

    @AfterClass
    public void cleanUp() throws KubernetesPluginException {
        KubernetesUtils.deleteDirectory(String.valueOf(TARGET_PATH));
    }
}
