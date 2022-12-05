package provider

import (
	"fmt"
	"testing"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceServer(t *testing.T) {
	// var serverId string
	// var backendId string
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServerConfig("tfacc-server1", 9090, "127.0.0.1", "tfacc-backend2", "enabled", "roundrobin"), // backend name related to backend test
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("haproxy_server.tfacc-server1", "name", "tfacc-server1"),
					// extract serverId for future use
					func(s *terraform.State) error {						
						// serverId = s.RootModule().Resources["haproxy_server.tfacc-server1"].Primary.Attributes["id"]
						// backendId = s.RootModule().Resources["haproxy_backend.tfacc-backend2"].Primary.Attributes["id"]
						return nil
					},
				),
			},
		},
	})
}

func testAccServerConfig(name string, port int, address string, parent_name string, check string, algo string) string {
	backend := testAccFrontendServerConfig(parent_name, algo)
	server := fmt.Sprintf(`
resource "haproxy_server" "%[1]s" {
	name = "%[1]s"
	port = "%[2]d"
	address = "%[3]s"
	parent_name = "%[4]s"
	check = "%[5]s"
}`, name, port, address, parent_name, check)

return strings.Join([]string{backend, server}, "\n")
}

func testAccFrontendServerConfig(name string, algo string) string {
	return fmt.Sprintf(`
resource "haproxy_backend" "%[1]s" {
	name = "%[1]s"
	balance_algorithm =  "%[2]s"
}	
`, name, algo)
}
