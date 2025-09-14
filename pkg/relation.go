package pkg

var Relations = map[string]map[string][]string{
	"post": {
		"owner":  {"read", "write", "delete"},
		"editor": {"read", "write"},
		"viewer": {"read"},
	},
}

func PermissionsInRelation(namespace, relation string) []string {
	_, ok := Relations[namespace]
	if !ok {
		return []string{}
	}

	_, ok = Relations[namespace][relation]
	if !ok {
		return []string{}
	}

	return Relations[namespace][relation]
}

func RelationsForPermission(namespace string, permission string) []string {
	permissions := Permissions()
	_, ok := permissions[namespace]
	if !ok {
		return []string{}
	}

	_, ok = permissions[namespace][permission]
	if !ok {
		return []string{}
	}

	return permissions[namespace][permission]
}

func Permissions() map[string]map[string][]string {
	ps := map[string]map[string][]string{}
	for namespace, relationToPermissions := range Relations {
		ps[namespace] = map[string][]string{}
		for relation, permissions := range relationToPermissions {
			for _, p := range permissions {
				ps[namespace][p] = append(ps[namespace][p], relation)
			}
		}
	}
	return ps
}
