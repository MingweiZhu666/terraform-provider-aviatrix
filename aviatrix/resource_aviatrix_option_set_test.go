package aviatrix

import (
	"context"
	"fmt"
	"github.com/AviatrixSystems/terraform-provider-aviatrix/v2/goaviatrix"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

func TestAccAviatrixOptionSet_basic(t *testing.T) {

	skipAcc := os.Getenv("SKIP_OPTION_SET")
	if skipAcc == "yes" {
		t.Skip("Skipping Option Set tests as SKIP_OPTION_SET is set.")
	}
	msgCommon := "Set SKIP_OPTION_SET to yes to skip Option Set tests."
	resourceName := "aviatrix_option_set.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			preAccountCheck(t, msgCommon)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAviatrixOptionSetBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccAviatrixOptionSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "option_set_Acc_test"),
					resource.TestCheckResourceAttr(resourceName, "search_domains.0.search_domain", "internal.compute"),
					resource.TestCheckResourceAttr(resourceName, "name_servers.0.transport", "udp"),
				),
			},
		},
	})
}

func testAccAviatrixOptionSetBasic() string {
	return fmt.Sprintf(`
resource "aviatrix_option_set" "test" {
  name       = "option_set_Acc_test"
  search_domains {
    search_domain = "internal.compute"
    server        = "test.123"
  }
  search_domains {
    search_domain = "compute"
    server        = "test"
  }
  name_servers {
    server       = "8.8.8.9"
    dot          = false
    transport    = "udp"
  }
  name_servers {
    server       = "8.9.9.9"
    dot          = false
    transport    = "udp"
  }
}
`)
}

func testAccAviatrixOptionSetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Option Set Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Option Set ID is set")
		}
		client := testAccProvider.Meta().(*goaviatrix.Client)

		optionSetName := rs.Primary.Attributes["name"]
		_, err := client.GetOptionSet(context.Background(), optionSetName)
		if err != nil {
			return err
		}
		return nil
	}
}
