package jdcloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/client"
	"strings"
)

func resourceJDCloudAcl() *schema.Resource {

	return &schema.Resource{

		Create: resourceJDCloudAclCreate,
		Read:   resourceJDCloudAclRead,
		Update: resourceJDCloudAclUpdate,
		Delete: resourceJDCloudAclDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Required:    true,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"network_aclname": {
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceJDCloudAclCreate(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*JDCloudConfig)
	vpcClient := client.NewVpcClient(config.Credential)

	//description := d.Get("description").(string)
	vpcId := d.Get("vpc_id").(string)
	networkAclName := strings.Trim(d.Get("network_aclname").(string), " ")
	//description := d.Get("description").(string)

	rq := apis.NewCreateNetworkAclRequest(config.Region, vpcId, networkAclName)

	if _, ok := d.GetOk("description"); ok {
		rq.Description = GetStringAddr(d, "description")
	}

	//发送请求
	resp, err := vpcClient.CreateNetworkAcl(rq)

	if err != nil {
		return fmt.Errorf("[ERROR] resourceVpcCreate failed %s ", err.Error())
	}

	if resp.Error.Code != 0 {
		return fmt.Errorf("[ERROR] resourceVpcCreate failed  code:%d staus:%s message:%s ", resp.Error.Code, resp.Error.Status, resp.Error.Message)
	}
	d.SetId(resp.Result.NetworkAclId)
	//d.Set("acl_id", resp.Result.NetworkAclId)
	return nil
}

func resourceJDCloudAclRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceJDCloudAclUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceJDCloudAclDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*JDCloudConfig)
	vpcClient := client.NewVpcClient(config.Credential)

	//构造请求
	NetworkAclId := d.Id()
	rq := apis.NewDeleteNetworkAclRequest(config.Region, NetworkAclId)

	//return fmt.Errorf("acl_id %s",NetworkAclId)
	//发送请求
	resp, err := vpcClient.DeleteNetworkAcl(rq)

	if err != nil {
		return fmt.Errorf("[ERROR] resourceVpcCreate failed %s ", err.Error())
	}

	if resp.Error.Code != 0 {
		return fmt.Errorf("[ERROR] resourceVpcCreate failed  code:%d staus:%s message:%s ", resp.Error.Code, resp.Error.Status, resp.Error.Message)
	}

	return nil
}