package mockutils

func UpdateEntity[T interface{}](oldEntities []T, newEntity T, getIdentifier func(e T) string, id string) []T {
	return updateEntities(oldEntities, newEntity, getIdentifier, id, true)
}

func DeleteEntity[T interface{}](oldEntities []T, getIdentifier func(e T) string, id string) []T {
	return updateEntities(oldEntities, oldEntities[0], getIdentifier, id, true)
}

func updateEntities[T interface{}](oldEntities []T, newEntity T, getIdentifier func(e T) string, id string, shouldAppend bool) []T {
	updatedEntities := []T{}
	for _, entity := range oldEntities {
		if getIdentifier(entity) != id {
			updatedEntities = append(updatedEntities, entity)
		}
	}
	if shouldAppend {
		updatedEntities = append(
			updatedEntities, newEntity,
		)
	}
	return updatedEntities
}
