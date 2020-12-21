package utility

import (
	"strconv"

	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

//IDtoInt convert pulumi.ID to pulumi.IntOutput
func IDtoInt(crs pulumi.CustomResourceState) pulumi.IntOutput {
	return crs.ID().ApplyInt(func(id pulumi.ID) int {
		number, err := strconv.Atoi(string(id))
		if err != nil {
			panic(err)
		}
		return number
	})
}
