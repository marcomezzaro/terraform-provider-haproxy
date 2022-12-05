package provider

import (
	"fmt"
	"testing"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceBind(t *testing.T) {
	// var bindId string
	// var frontendId string
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBindConfig("tfacc-bind1", 9090, "127.0.0.1", "tfacc-frontend2"), // frontend name related to frontend test
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("haproxy_bind.tfacc-bind1", "name", "tfacc-bind1"),
					// extract bindId for future use
					func(s *terraform.State) error {						
						// bindId = s.RootModule().Resources["haproxy_bind.tfacc-bind1"].Primary.Attributes["id"]
						// frontendId = s.RootModule().Resources["haproxy_frontend.tfacc-frontend2"].Primary.Attributes["id"]
						return nil
					},
				),
			},
		},
	})
}

func testAccBindConfig(name string, port int, address string, parent_name string) string {
	frontend := testAccFrontendBindConfig(parent_name)
	bind := fmt.Sprintf(`
resource "haproxy_bind" "%[1]s" {
	name = "%[1]s"
	port = "%[2]d"
	address = "%[3]s"
	parent_name = "%[4]s"
}`, name, port, address, parent_name)

return strings.Join([]string{frontend, bind}, "\n")
}

func testAccFrontendBindConfig(name string) string {
	return fmt.Sprintf(`
resource "haproxy_frontend" "%[1]s" {
	name = "%[1]s"
}	
`, name)
}
