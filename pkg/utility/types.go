package utility

import (
	"strconv"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//IDtoInt convert pulumi.ID to pulumi.IntOutput
func IDtoInt(crs pulumi.CustomResourceState) pulumi.IntOutput {
	return crs.ID().ApplyT(func(id pulumi.ID) int {
		number, err := strconv.Atoi(string(id))
		if err != nil {
			panic(err)
		}
		return number
	}).(pulumi.IntOutput)
}
