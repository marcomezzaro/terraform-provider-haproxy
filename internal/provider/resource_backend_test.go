package provider

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceBackend(t *testing.T) {
	var backendId string
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBackendConfig("tfacc-backend1", "roundrobin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("haproxy_backend.test", "name", "tfacc-backend1"),

					// extract backendId for future use
					func(s *terraform.State) error {
						backendId = s.RootModule().Resources["haproxy_backend.test"].Primary.Attributes["id"]
						return nil
					},
				),
			},
			importStep("haproxy_backend.test"),
			{
				ResourceName: "haproxy_backend.test",
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return backendId, nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBackendConfig(name string, algo string) string {
	return fmt.Sprintf(`
resource "haproxy_backend" "test" {
	name = "%[1]s"
	balance_algorithm = "%[2]s"
}	
`, name, algo)
}
