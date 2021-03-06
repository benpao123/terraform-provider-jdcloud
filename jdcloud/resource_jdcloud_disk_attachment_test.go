package jdcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
	"testing"
	"time"
)

const TestAccDiskAttachmentConfig = `
resource "jdcloud_disk_attachment" "disk-attachment-TEST-1"{
	instance_id = "i-96ef5rv62n" 
	disk_id = "vol-39dmz9csj6"
}
`

func TestAccJDCloudDiskAttachment_basic(t *testing.T) {

	var instanceId, diskId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDiskAttachmentDestroy(&instanceId, &diskId),
		Steps: []resource.TestStep{
			{
				Config: TestAccDiskAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(

					testAccIfDiskAttachmentExists("jdcloud_disk_attachment.disk-attachment-TEST-1", &instanceId, &diskId),
				),
			},
		},
	})
}

//-------------------------- Customized check functions

func testAccIfDiskAttachmentExists(resourceName string, resourceId, diskId *string) resource.TestCheckFunc {

	return func(stateInfo *terraform.State) error {

		time.Sleep(time.Second * 15)

		infoStoredLocally, ok := stateInfo.RootModule().Resources[resourceName]
		if ok == false {
			return fmt.Errorf("we can not find a resouce namely:{%s} in terraform.State", resourceName)
		}
		if infoStoredLocally.Primary.ID == "" {
			return fmt.Errorf("operation failed, resource:%s is created but ID not set", resourceName)
		}
		*resourceId = infoStoredLocally.Primary.Attributes["instance_id"]
		*diskId = infoStoredLocally.Primary.Attributes["disk_id"]

		config := testAccProvider.Meta().(*JDCloudConfig)
		vmClient := client.NewVmClient(config.Credential)

		req := apis.NewDescribeInstanceRequest(config.Region, *resourceId)
		resp, err := vmClient.DescribeInstance(req)

		if err != nil {
			return err
		}

		expectedCloudDiskNotFound := true
		for _, aDisk := range resp.Result.Instance.DataDisks {
			if aDisk.CloudDisk.DiskId == *diskId {
				expectedCloudDiskNotFound = false
			}
		}
		if expectedCloudDiskNotFound {
			return fmt.Errorf("resource not found remotely")
		}

		return nil
	}
}

func testAccDiskAttachmentDestroy(resourceId *string, diskId *string) resource.TestCheckFunc {

	return func(stateInfo *terraform.State) error {

		config := testAccProvider.Meta().(*JDCloudConfig)
		vmClient := client.NewVmClient(config.Credential)

		req := apis.NewDescribeInstanceRequest(config.Region, *resourceId)
		resp, err := vmClient.DescribeInstance(req)

		if err != nil {
			return err
		}

		if resp.Result.Instance.DataDisks[0].Status == "detached" || resp.Result.Instance.DataDisks[0].Status == "detaching" {
			return nil
		}

		return fmt.Errorf("data disk failed in detatching ")
	}
}
