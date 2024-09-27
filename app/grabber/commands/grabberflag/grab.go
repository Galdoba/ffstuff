package grabberflag

import (
	"github.com/Galdoba/ffstuff/app/grabber/internal/validation"
	"github.com/Galdoba/ffstuff/pkg/logman"
	"github.com/urfave/cli/v2"
)

func ValidateGrabFlags(c *cli.Context) error {
	destination_value := c.String(DESTINATION)
	switch destination_value {
	default:
		if err := validation.DirectoryValidation(destination_value); err != nil {
			return logman.Errorf("invalid -%v flag: %v", DESTINATION, err)
		}
	case "":
	}
	copy_value := c.String(COPY)
	switch copy_value {
	default:
		return logman.Errorf("invalid -%v flag value: '%v'", COPY, copy_value)
	case VALUE_COPY_SKIP, VALUE_COPY_RENAME, VALUE_COPY_OVERWRITE:
	}
	delete_value := c.String(DELETE)
	switch delete_value {
	default:
		return logman.Errorf("invalid -%v flag value: '%v'", DELETE, delete_value)
	case VALUE_DELETE_NONE, VALUE_DELETE_MARKER, VALUE_DELETE_ALL:
	}
	sort_value := c.String(SORT)
	switch sort_value {
	default:
		return logman.Errorf("invalid -%v flag value: '%v'", SORT, sort_value)
	case VALUE_SORT_PRIORITY, VALUE_SORT_SIZE, VALUE_SORT_NONE:
	}

	return nil
}
