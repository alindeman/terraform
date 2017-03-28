package heroku

import (
	"context"
	"fmt"
	"os"
	"testing"

	heroku "github.com/cyberdelia/heroku-go/v3"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHerokuSpace_Basic(t *testing.T) {
	var space *heroku.Space
	org := os.Getenv("HEROKU_ORGANIZATION")
	enterprise := os.Getenv("HEROKU_ENTERPRISE")
	spaceName := fmt.Sprintf("tftest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if org == "" {
				t.Skip("HEROKU_ORGANIZATION is not set; skipping test")
			}
			if enterprise == "" {
				t.Skip("HEROKU_ENTERPRISE is not set; skipping test")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHerokuSpaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHerokuSpaceConfig(organizationName, spaceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHerokuSpaceExists("heroku_space.foobar", &space),
				),
			},
		},
	})
}

func testAccCheckHerokuSpaceExists(n string, space *heroku.Space) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No Space ID is set")
		}

		client := testAccProvider.Meta().(*heroku.Service)

		foundSpace, err := client.SpaceInfo(context.TODO(), id)
		if err != nil {
			return err
		}

		*space = *foundSpace
		return nil
	}
}

func testAccHerokuSpaceConfig(organizationName string, spaceName string) string {
	return fmt.Sprintf(`
  resource "heroku_space" "foobar" {
    name = "%s"
    organization = "%s"
    region = "us"
    shield = false
  }
`, spaceName, organizationName)
}

func testAccCheckHerokuSpaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*heroku.Service)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "heroku_space" {
			continue
		}

		id := rs.Primary.ID
		space, err := client.SpaceInfo(context.TODO(), id)
		if err == nil && space.ID == id {
			return fmt.Errorf("space '%s' still exists", space.ID)
		}
		return err
	}
	return nil
}
