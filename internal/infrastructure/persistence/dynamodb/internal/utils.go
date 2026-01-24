package internal

import "fmt"

func CreateLinkMetaPartitionKey(id string) PrimaryKey {
	pk := fmt.Sprintf("LINK#%s", id)
	sk := "META"

	return PrimaryKey{PK: pk, SK: sk}
}
