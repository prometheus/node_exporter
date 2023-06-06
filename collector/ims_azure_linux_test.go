// Copyright 2021 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const azureMockData = `
{
	"azEnvironment": "AZUREPUBLICCLOUD",
    "additionalCapabilities": {
        "hibernationEnabled": "true"
    },
    "hostGroup": {
      "id": "testHostGroupId"
    },
    "extendedLocation": {
        "type": "edgeZone",
        "name": "microsoftlosangeles"
    },
    "evictionPolicy": "",
    "isHostCompatibilityLayerVm": "true",
    "licenseType":  "",
    "location": "westus",
	"provider": "Microsoft.Compute",
    "name": "examplevmname",
    "offer": "UbuntuServer",
    "osProfile": {
        "adminUsername": "admin",
        "computerName": "examplevmname",
        "disablePasswordAuthentication": "true"
    },
    "osType": "Linux",
    "placementGroupId": "f67c14ab-e92c-408c-ae2d-da15866ec79a",
    "plan": {
        "name": "planName",
        "product": "planProduct",
        "publisher": "planPublisher"
    },
    "platformFaultDomain": "36",
    "platformSubFaultDomain": "",
    "platformUpdateDomain": "42",
    "priority": "Regular",
    "publicKeys": [{
            "keyData": "ssh-rsa 0",
            "path": "/home/user/.ssh/authorized_keys0"
        },
        {
            "keyData": "ssh-rsa 1",
            "path": "/home/user/.ssh/authorized_keys1"
        }
    ],
    "publisher": "Canonical",
    "resourceGroupName": "macikgo-test-may-23",
    "resourceId": "/subscriptions/xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx/resourceGroups/macikgo-test-may-23/providers/Microsoft.Compute/virtualMachines/examplevmname",
    "securityProfile": {
        "secureBootEnabled": "true",
        "virtualTpmEnabled": "false",
        "encryptionAtHost": "true",
        "securityType": "TrustedLaunch"
    },
    "sku": "18.04-LTS",
    "storageProfile": {
        "dataDisks": [{
            "bytesPerSecondThrottle": "979202048",
            "caching": "None",
            "createOption": "Empty",
            "diskCapacityBytes": "274877906944",
            "diskSizeGB": "1024",
            "image": {
              "uri": ""
            },
            "isSharedDisk": "false",
            "isUltraDisk": "true",
            "lun": "0",
            "managedDisk": {
              "id": "/subscriptions/xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx/resourceGroups/macikgo-test-may-23/providers/Microsoft.Compute/disks/exampledatadiskname",
              "storageAccountType": "StandardSSD_LRS"
            },
            "name": "exampledatadiskname",
            "opsPerSecondThrottle": "65280",
            "vhd": {
              "uri": ""
            },
            "writeAcceleratorEnabled": "false"
        }],
        "imageReference": {
            "id": "",
            "offer": "UbuntuServer",
            "publisher": "Canonical",
            "sku": "16.04.0-LTS",
            "version": "latest"
        },
        "osDisk": {
            "caching": "ReadWrite",
            "createOption": "FromImage",
            "diskSizeGB": "30",
            "diffDiskSettings": {
                "option": "Local"
            },
            "encryptionSettings": {
              "enabled": "false",
              "diskEncryptionKey": {
                "sourceVault": {
                  "id": "/subscriptions/test-source-guid/resourceGroups/testrg/providers/Microsoft.KeyVault/vaults/test-kv"
                },
                "secretUrl": "https://test-disk.vault.azure.net/secrets/xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx/xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
              },
              "keyEncryptionKey": {
                "sourceVault": {
                  "id": "/subscriptions/test-key-guid/resourceGroups/testrg/providers/Microsoft.KeyVault/vaults/test-kv"
                },
                "keyUrl": "https://test-key.vault.azure.net/secrets/xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx/xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx"
              }
            },
            "image": {
                "uri": ""
            },
            "managedDisk": {
                "id": "/subscriptions/xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx/resourceGroups/macikgo-test-may-23/providers/Microsoft.Compute/disks/exampleosdiskname",
                "storageAccountType": "StandardSSD_LRS"
            },
            "name": "exampleosdiskname",
            "osType": "Linux",
            "vhd": {
                "uri": ""
            },
            "writeAcceleratorEnabled": "false"
        },
        "resourceDisk": {
            "size": "4096"
        }
    },
    "subscriptionId": "xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
    "tags": "baz:bash;foo:bar",
    "version": "15.05.22",
    "virtualMachineScaleSet": {
        "id": "/subscriptions/xxxxxxxx-xxxxx-xxx-xxx-xxxx/resourceGroups/resource-group-name/providers/Microsoft.Compute/virtualMachineScaleSets/virtual-machine-scale-set-name"
    },
    "vmId": "02aab8a4-74ef-476e-8182-f6d2ba4166a6",
    "vmScaleSetName": "crpteste9vflji9",
    "vmSize": "Standard_A3",
    "zone": "1"
}`

func Test_NewAzureInstanceMetaDataColector(t *testing.T) {
	instance := NewAzureInstanceMetadataCollector("description", "fqdn", nil)

	if len(instance.labels) != 4 {
		t.Fatal("expected 4 labels for azure instance metadata collector")
	}
	if len(instance.mapping) != 4 {
		t.Fatal("expected 4 entries at the mapping configuration for azure instance metadata collector")
	}

	desc := instance.InfoDesc.String()
	if !strings.Contains(desc, "fqName: \"fqdn\"") {
		t.Fatal("expected the required fqdn \"fqdn\"")
	}
}

func Test_GetMetaDataServiceInfo(t *testing.T) {
	azureImsC := &AzureInstanceMetadataCollector{}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, azureMockData)
	}))

	defer mockServer.Close()
	azureImsC.computeAPI = mockServer.URL

	data, err := azureImsC.GetMetaDataServiceInfo()
	if err != nil {
		t.Fatal("expected no error at GetMetaDataServiceInfo")
	}

	if len(data) == 0 {
		t.Fatal("expected some data that returns from mock server")
	}
}

func TestImsAzureResponse(t *testing.T) {
	azureMockServer := func(collector Collector) *httptest.Server {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, azureMockData)
		}))

		collector.(*AzureInstanceMetadataCollector).computeAPI = mockServer.URL
		return mockServer
	}

	testcases := []struct {
		labels      string
		metricsFile string
		setup       func(Collector) *httptest.Server
	}{
		{"", "", nil},
		{"azure", "fixtures/ims/azure_ims_result.txt", azureMockServer},
		{"unknown", "", nil},
	}

	for _, test := range testcases {
		t.Run(test.labels, func(t *testing.T) {
			args := []string{fmt.Sprintf("--collector.ims.provider=%s", test.labels)}

			if _, err := kingpin.CommandLine.Parse(args); err != nil {
				t.Fatal(err)
			}

			collector, err := NewImsCollector(log.NewNopLogger())
			if err != nil {
				t.Fatal(err)
			}

			if test.setup != nil {
				mockServer := test.setup(collector)
				defer mockServer.Close()
			}

			registry := prometheus.NewRegistry()
			registry.MustRegister(miniCollector{c: collector})

			rw := httptest.NewRecorder()
			promhttp.InstrumentMetricHandler(registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{})).ServeHTTP(rw, &http.Request{})

			if len(test.metricsFile) > 0 {
				wantMetrics, err := os.ReadFile(test.metricsFile)
				if err != nil {
					t.Fatalf("unable to read input test file %s: %s", test.metricsFile, err)
				}

				wantLines := strings.Split(string(wantMetrics), "\n")
				gotLines := strings.Split(string(rw.Body.String()), "\n")
				gotLinesIdx := 0

				// Until the Prometheus Go client library offers better testability
				// (https://github.com/prometheus/client_golang/issues/58), we simply compare
				// verbatim text-format metrics outputs, but ignore any lines we don't have
				// in the fixture. Put differently, we are only testing that each line from
				// the fixture is present, in the order given.
			wantLoop:
				for _, want := range wantLines {
					for _, got := range gotLines[gotLinesIdx:] {
						if want == got {
							// this is a line we are interested in, and it is correct
							continue wantLoop
						} else {
							gotLinesIdx++
						}
					}
					// if this point is reached, the line we want was missing
					t.Fatalf("Missing expected output line(s), first missing line is %s", want)
				}
			} else {
				gotLines := strings.Split(string(rw.Body.String()), "\n")
				fmt.Println(gotLines)
				for i := range gotLines {
					if strings.Contains(gotLines[i], "node_ims_info") {
						t.Fatal("Found ims information that is not expected")
					}
				}
			}
		})
	}
}
