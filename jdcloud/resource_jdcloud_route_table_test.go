package jdcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/client"
	"testing"
)

const TestAccRouteTableConfig = `
resource "jdcloud_route_table" "route-table-TEST-1"{
	route_table_name = "route_table_test"
	vpc_id = "vpc-npvvk4wr5j"
	description = "test"
}
`

func TestAccJDCloudRouteTable_basic(t *testing.T) {

	// routeTableId is declared but not assigned any values here
	// It will be assigned value in "testAccIfRouteTableExists"
	var routeTableId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccRouteTableDestroy(&routeTableId),
		Steps: []resource.TestStep{
			{
				Config: TestAccRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(

					// ROUTE_TABLE_ID validation
					testAccIfRouteTableExists("jdcloud_route_table.route-table-TEST-1", &routeTableId),
					// Remaining attributes validation
					resource.TestCheckResourceAttr("jdcloud_route_table.route-table-TEST-1", "route_table_name", "route_table_test"),
					resource.TestCheckResourceAttr("jdcloud_route_table.route-table-TEST-1", "vpc_id", "vpc-npvvk4wr5j"),
					resource.TestCheckResourceAttr("jdcloud_route_table.route-table-TEST-1", "description", "test"),
				),
			},
		},
	})
}

func testAccIfRouteTableExists(routeTableName string, routeTableId *string) resource.TestCheckFunc {

	return func(stateInfo *terraform.State) error {

		// STEP-1: Check if RouteTable resource has been created locally
		routeTableInfoStoredLocally, ok := stateInfo.RootModule().Resources[routeTableName]
		if ok == false {
			return fmt.Errorf("we can not find a RouteTable namely:{%s} in terraform.State", routeTableName)
		}
		if routeTableInfoStoredLocally.Primary.ID == "" {
			return fmt.Errorf("operation failed, RouteTable is created but ID not set")
		}
		routeTableIdStoredLocally := routeTableInfoStoredLocally.Primary.ID

		// STEP-2 : Check if RouteTable resource has been created remotely
		routeTableconfig := testAccProvider.Meta().(*JDCloudConfig)
		routeTableClient := client.NewVpcClient(routeTableconfig.Credential)

		requestOnRouteTable := apis.NewDescribeRouteTableRequest(routeTableconfig.Region, routeTableIdStoredLocally)
		responseOnRouteTable, err := routeTableClient.DescribeRouteTable(requestOnRouteTable)

		if err != nil {
			return err
		}
		if responseOnRouteTable.Error.Code != 0 {
			return fmt.Errorf("according to the ID stored locally,we cannot find any RouteTable created remotely")
		}

		// RouteTable ID has been validated
		// We are going to validate the remaining attributes - name,vpc_id,description
		*routeTableId = routeTableIdStoredLocally
		return nil
	}
}

func testAccRouteTableDestroy(routeTableIdStoredLocally *string) resource.TestCheckFunc {

	return func(stateInfo *terraform.State) error {

		//  routeTableId is not supposed to be empty
		if *routeTableIdStoredLocally == "" {
			return fmt.Errorf("route Table Id appears to be empty")
		}

		routeTableConfig := testAccProvider.Meta().(*JDCloudConfig)
		routeTableClient := client.NewVpcClient(routeTableConfig.Credential)

		routeTableRegion := routeTableConfig.Region
		requestOnRouteTable := apis.NewDescribeRouteTableRequest(routeTableRegion, *routeTableIdStoredLocally)
		responseOnRouteTable, err := routeTableClient.DescribeRouteTable(requestOnRouteTable)

		// Error.Code is supposed to be 404 since RouteTable was actually deleted
		// Meanwhile turns out to be 0, successfully queried. Indicating delete error
		if err != nil {
			return err
		}
		if responseOnRouteTable.Error.Code == 0 {
			return fmt.Errorf("routeTable resource still exists,check position-4")
		}

		return nil
	}
}
