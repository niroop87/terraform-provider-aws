package aws

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAwsAutoScalingGroupDataSource_basic(t *testing.T) {
	datasourceName := "data.aws_autoscaling_group.good_match"
	resourceName := "aws_autoscaling_group.foo"
	rName := fmt.Sprintf("tf-test-asg-%d", acctest.RandInt())
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingGroupDataResourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(datasourceName, "availability_zones.#", resourceName, "availability_zones.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "default_cooldown", resourceName, "default_cooldown"),
					resource.TestCheckResourceAttrPair(datasourceName, "desired_capacity", resourceName, "desired_capacity"),
					resource.TestCheckResourceAttrPair(datasourceName, "health_check_grace_period", resourceName, "health_check_grace_period"),
					resource.TestCheckResourceAttrPair(datasourceName, "health_check_type", resourceName, "health_check_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "launch_configuration", resourceName, "launch_configuration"),
					resource.TestCheckResourceAttrPair(datasourceName, "load_balancers.#", resourceName, "load_balancers.#"),
					resource.TestCheckResourceAttr(datasourceName, "new_instances_protected_from_scale_in", "false"),
					resource.TestCheckResourceAttrPair(datasourceName, "max_size", resourceName, "max_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "min_size", resourceName, "min_size"),
					resource.TestCheckResourceAttrPair(datasourceName, "target_group_arns.#", resourceName, "target_group_arns.#"),
					resource.TestCheckResourceAttr(datasourceName, "vpc_zone_identifier", ""),
				),
			},
		},
	})
}

// Lookup based on AutoScalingGroupName
func testAccAutoScalingGroupDataResourceConfig(rName string) string {
	return composeConfig(
		testAccLatestAmazonLinuxHvmEbsAmiConfig(),
		testAccAvailableAZsNoOptInConfig(),
		testAccAvailableEc2InstanceTypeForAvailabilityZone("data.aws_availability_zones.available.names[0]", "t3.micro", "t2.micro"),
		fmt.Sprintf(`
resource "aws_launch_configuration" "data_source_aws_autoscaling_group_test" {
  name          = "%[1]s"
  image_id      = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = data.aws_ec2_instance_type_offering.available.instance_type
}

resource "aws_autoscaling_group" "foo" {
  name                      = "%[1]s_foo"
  max_size                  = 0
  min_size                  = 0
  health_check_grace_period = 300
  health_check_type         = "ELB"
  desired_capacity          = 0
  force_delete              = true
  launch_configuration      = aws_launch_configuration.data_source_aws_autoscaling_group_test.name
  availability_zones        = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]
}

resource "aws_autoscaling_group" "bar" {
  name                      = "%[1]s_bar"
  max_size                  = 0
  min_size                  = 0
  health_check_grace_period = 300
  health_check_type         = "ELB"
  desired_capacity          = 0
  force_delete              = true
  launch_configuration      = aws_launch_configuration.data_source_aws_autoscaling_group_test.name
  availability_zones        = [data.aws_availability_zones.available.names[0], data.aws_availability_zones.available.names[1]]
}

data "aws_autoscaling_group" "good_match" {
  name       = aws_autoscaling_group.foo.name
  depends_on = [aws_autoscaling_group.foo]
}
`, rName))
}
