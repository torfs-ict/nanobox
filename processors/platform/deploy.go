package platform

import (
  "fmt"
  
  "github.com/nanobox-io/nanobox/models"
)

// Deploy provisions platform components needed for an app deploy
func Deploy(a *models.App) error {
  
  for _, component := range deployComponents {
    if err := provisionComponent(a, component); err != nil {
      return fmt.Errorf("failed to provision platform component: %s", err.Error())
    }
  }

  return nil
}