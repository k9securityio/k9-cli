package core

const (
	DIFF_DELETED = `deleted`
	DIFF_ADDED   = `added`
	DIFF_CHANGED = `changed`
)

// PrincipalsReportItemDifference represents the differences between two versions
// of the same PrincipalsReportItem (correlated by PrincipalARN).
type PrincipalsReportItemDifference struct {
	Type                           string `csv:"type"`
	PrincipalARN                   string `csv:"principal_arn"`
	BeforePrincipalName            string `csv:"before_principal_name"`
	BeforePrincipalType            string `csv:"before_principal_type"`
	BeforePrincipalIsIAMAdmin      bool   `csv:"before_principal_is_iam_admin"`
	BeforePrincipalLastUsed        string `csv:"before_principal_last_used"`
	BeforePrincipalTagBusinessUnit string `csv:"before_principal_tag_business_unit"`
	BeforePrincipalTagEnvironment  string `csv:"before_principal_tag_environment"`
	BeforePrincipalTagUsedBy       string `csv:"before_principal_tag_used_by"`
	BeforePrincipalTags            string `csv:"before_principal_tags"`
	BeforePasswordLastUsed         string `csv:"before_password_last_used"`
	BeforePasswordLastRotated      string `csv:"before_password_last_rotated"`
	BeforePasswordState            string `csv:"before_password_state"`
	BeforeAccessKey1LastUsed       string `csv:"before_access_key_1_last_used"`
	BeforeAccessKey1LastRotated    string `csv:"before_access_key_1_last_rotated"`
	BeforeAccessKey1State          string `csv:"before_access_key_1_state"`
	BeforeAccessKey2LastUsed       string `csv:"before_access_key_2_last_used"`
	BeforeAccessKey2LastRotated    string `csv:"before_access_key_2_last_rotated"`
	BeforeAccessKey2State          string `csv:"before_access_key_2_state"`
	AfterPrincipalName             string `csv:"after_principal_name"`
	AfterPrincipalType             string `csv:"after_principal_type"`
	AfterPrincipalIsIAMAdmin       bool   `csv:"after_principal_is_iam_admin"`
	AfterPrincipalLastUsed         string `csv:"after_principal_last_used"`
	AfterPrincipalTagBusinessUnit  string `csv:"after_principal_tag_business_unit"`
	AfterPrincipalTagEnvironment   string `csv:"after_principal_tag_environment"`
	AfterPrincipalTagUsedBy        string `csv:"after_principal_tag_used_by"`
	AfterPrincipalTags             string `csv:"after_principal_tags"`
	AfterPasswordLastUsed          string `csv:"after_password_last_used"`
	AfterPasswordLastRotated       string `csv:"after_password_last_rotated"`
	AfterPasswordState             string `csv:"after_password_state"`
	AfterAccessKey1LastUsed        string `csv:"after_access_key_1_last_used"`
	AfterAccessKey1LastRotated     string `csv:"after_access_key_1_last_rotated"`
	AfterAccessKey1State           string `csv:"after_access_key_1_state"`
	AfterAccessKey2LastUsed        string `csv:"after_access_key_2_last_used"`
	AfterAccessKey2LastRotated     string `csv:"after_access_key_2_last_rotated"`
	AfterAccessKey2State           string `csv:"after_access_key_2_state"`
}

// AddedDiff produces a new PrincipalsReportItemDifference with fields set from the
// receiver PrincipalsReportItem in the "after" columns, and the type set to
// DIFF_ADDED.
func (i PrincipalsReportItem) AddedDiff() PrincipalsReportItemDifference {
	return PrincipalsReportItemDifference{
		PrincipalARN:                  i.PrincipalARN,
		Type:                          DIFF_ADDED,
		AfterPrincipalName:            i.PrincipalName,
		AfterPrincipalType:            i.PrincipalType,
		AfterPrincipalIsIAMAdmin:      i.PrincipalIsIAMAdmin,
		AfterPrincipalLastUsed:        i.PrincipalLastUsed,
		AfterPrincipalTagBusinessUnit: i.PrincipalTagBusinessUnit,
		AfterPrincipalTagEnvironment:  i.PrincipalTagEnvironment,
		AfterPrincipalTagUsedBy:       i.PrincipalTagUsedBy,
		AfterPrincipalTags:            i.PrincipalTags,
		AfterPasswordLastUsed:         i.PasswordLastUsed,
		AfterPasswordLastRotated:      i.PasswordLastRotated,
		AfterPasswordState:            i.PasswordState,
		AfterAccessKey1LastUsed:       i.AccessKey1LastUsed,
		AfterAccessKey1LastRotated:    i.AccessKey1LastRotated,
		AfterAccessKey1State:          i.AccessKey1State,
		AfterAccessKey2LastUsed:       i.AccessKey2LastUsed,
		AfterAccessKey2LastRotated:    i.AccessKey2LastRotated,
		AfterAccessKey2State:          i.AccessKey2State,
	}
}

// DeletedDiff produces a new PrincipalsReportItemDifference with fields set from the
// receiver PrincipalsReportItem in the "before" columns, and the type set to
// DIFF_DELETED.
func (i PrincipalsReportItem) DeletedDiff() PrincipalsReportItemDifference {
	return PrincipalsReportItemDifference{
		PrincipalARN:                   i.PrincipalARN,
		Type:                           DIFF_DELETED,
		BeforePrincipalName:            i.PrincipalName,
		BeforePrincipalType:            i.PrincipalType,
		BeforePrincipalIsIAMAdmin:      i.PrincipalIsIAMAdmin,
		BeforePrincipalLastUsed:        i.PrincipalLastUsed,
		BeforePrincipalTagBusinessUnit: i.PrincipalTagBusinessUnit,
		BeforePrincipalTagEnvironment:  i.PrincipalTagEnvironment,
		BeforePrincipalTagUsedBy:       i.PrincipalTagUsedBy,
		BeforePrincipalTags:            i.PrincipalTags,
		BeforePasswordLastUsed:         i.PasswordLastUsed,
		BeforePasswordLastRotated:      i.PasswordLastRotated,
		BeforePasswordState:            i.PasswordState,
		BeforeAccessKey1LastUsed:       i.AccessKey1LastUsed,
		BeforeAccessKey1LastRotated:    i.AccessKey1LastRotated,
		BeforeAccessKey1State:          i.AccessKey1State,
		BeforeAccessKey2LastUsed:       i.AccessKey2LastUsed,
		BeforeAccessKey2LastRotated:    i.AccessKey2LastRotated,
		BeforeAccessKey2State:          i.AccessKey2State,
	}
}

func (i PrincipalsReportItem) Diff(original PrincipalsReportItem) PrincipalsReportItemDifference {
	if i.PrincipalARN != original.PrincipalARN {
		panic(`comparing two different PrincipalReportItems`)
	}
	diff := PrincipalsReportItemDifference{
		PrincipalARN: i.PrincipalARN,
		Type:         DIFF_CHANGED}

	if i.PrincipalName != original.PrincipalName {
		diff.AfterPrincipalName = i.PrincipalName
		diff.BeforePrincipalName = original.PrincipalName
	}
	if i.PrincipalType != original.PrincipalType {
		diff.AfterPrincipalType = i.PrincipalType
		diff.BeforePrincipalType = original.PrincipalType
	}
	if i.PrincipalIsIAMAdmin != original.PrincipalIsIAMAdmin {
		diff.AfterPrincipalIsIAMAdmin = i.PrincipalIsIAMAdmin
		diff.BeforePrincipalIsIAMAdmin = original.PrincipalIsIAMAdmin
	}
	if i.PrincipalLastUsed != original.PrincipalLastUsed {
		diff.AfterPrincipalLastUsed = i.PrincipalLastUsed
		diff.BeforePrincipalLastUsed = original.PrincipalLastUsed
	}
	if i.PrincipalTagBusinessUnit != original.PrincipalTagBusinessUnit {
		diff.AfterPrincipalTagBusinessUnit = i.PrincipalTagBusinessUnit
		diff.BeforePrincipalTagBusinessUnit = original.PrincipalTagBusinessUnit
	}
	if i.PrincipalTagEnvironment != original.PrincipalTagEnvironment {
		diff.AfterPrincipalTagEnvironment = i.PrincipalTagEnvironment
		diff.BeforePrincipalTagEnvironment = original.PrincipalTagEnvironment
	}
	if i.PrincipalTagUsedBy != original.PrincipalTagUsedBy {
		diff.AfterPrincipalTagUsedBy = i.PrincipalTagUsedBy
		diff.BeforePrincipalTagUsedBy = original.PrincipalTagUsedBy
	}
	if i.PrincipalTags != original.PrincipalTags {
		diff.AfterPrincipalTags = i.PrincipalTags
		diff.BeforePrincipalTags = original.PrincipalTags
	}
	if i.PasswordLastUsed != original.PasswordLastUsed {
		diff.AfterPasswordLastUsed = i.PasswordLastUsed
		diff.BeforePasswordLastUsed = original.PasswordLastUsed
	}
	if i.PasswordLastRotated != original.PasswordLastRotated {
		diff.AfterPasswordLastRotated = i.PasswordLastRotated
		diff.BeforePasswordLastRotated = original.PasswordLastRotated
	}
	if i.PasswordState != original.PasswordState {
		diff.AfterPasswordState = i.PasswordState
		diff.BeforePasswordState = original.PasswordState
	}
	if i.AccessKey1LastUsed != original.AccessKey1LastUsed {
		diff.AfterAccessKey1LastUsed = i.AccessKey1LastUsed
		diff.BeforeAccessKey1LastUsed = original.AccessKey1LastUsed
	}
	if i.AccessKey1LastRotated != original.AccessKey1LastRotated {
		diff.AfterAccessKey1LastRotated = i.AccessKey1LastRotated
		diff.BeforeAccessKey1LastRotated = original.AccessKey1LastRotated
	}
	if i.AccessKey1State != original.AccessKey1State {
		diff.AfterAccessKey1State = i.AccessKey1State
		diff.BeforeAccessKey1State = original.AccessKey1State
	}
	if i.AccessKey2LastUsed != original.AccessKey2LastUsed {
		diff.AfterAccessKey2LastUsed = i.AccessKey2LastUsed
		diff.BeforeAccessKey2LastUsed = original.AccessKey2LastUsed
	}
	if i.AccessKey2LastRotated != original.AccessKey2LastRotated {
		diff.AfterAccessKey2LastRotated = i.AccessKey2LastRotated
		diff.BeforeAccessKey2LastRotated = original.AccessKey2LastRotated
	}
	if i.AccessKey2State != original.AccessKey2State {
		diff.AfterAccessKey2State = i.AccessKey2State
		diff.BeforeAccessKey2State = original.AccessKey2State
	}
	return diff
}

// ResourcesReportItemDifference represents the differences between two versions
// of the same ResourcesReportItem (correlated by ResourceARN).
type ResourcesReportItemDifference struct {
	Type        string `csv:"type"`
	ResourceARN string `csv:"resource_arn"`

	BeforeResourceName               string `csv:"before_resource_name"`
	BeforeResourceType               string `csv:"before_resource_type"`
	BeforeResourceTagBusinessUnit    string `csv:"before_resource_tag_business_unit"`
	BeforeResourceTagEnvironment     string `csv:"before_resource_tag_environment"`
	BeforeResourceTagOwner           string `csv:"before_resource_tag_owner"`
	BeforeResourceTagConfidentiality string `csv:"before_resource_tag_confidentiality"`
	BeforeResourceTagIntegrity       string `csv:"before_resource_tag_integrity"`
	BeforeResourceTagAvailability    string `csv:"before_resource_tag_availability"`
	BeforeResourceTags               string `csv:"before_resource_tags"`

	AfterResourceName               string `csv:"after_resource_name"`
	AfterResourceType               string `csv:"after_resource_type"`
	AfterResourceTagBusinessUnit    string `csv:"after_resource_tag_business_unit"`
	AfterResourceTagEnvironment     string `csv:"after_resource_tag_environment"`
	AfterResourceTagOwner           string `csv:"after_resource_tag_owner"`
	AfterResourceTagConfidentiality string `csv:"after_resource_tag_confidentiality"`
	AfterResourceTagIntegrity       string `csv:"after_resource_tag_integrity"`
	AfterResourceTagAvailability    string `csv:"after_resource_tag_availability"`
	AfterResourceTags               string `csv:"after_resource_tags"`
}

func (i ResourcesReportItem) Diff(original ResourcesReportItem) ResourcesReportItemDifference {
	if i.ResourceARN != original.ResourceARN {
		panic(`comparing two different ResourceReportItems`)
	}
	diff := ResourcesReportItemDifference{
		ResourceARN: i.ResourceARN,
		Type:        DIFF_CHANGED,
	}

	if i.ResourceName != original.ResourceName {
		diff.AfterResourceName = i.ResourceName
		diff.BeforeResourceName = original.ResourceName
	}
	if i.ResourceType != original.ResourceType {
		diff.AfterResourceType = i.ResourceType
		diff.BeforeResourceType = original.ResourceType
	}
	if i.ResourceTagBusinessUnit != original.ResourceTagBusinessUnit {
		diff.AfterResourceTagBusinessUnit = i.ResourceTagBusinessUnit
		diff.BeforeResourceTagBusinessUnit = original.ResourceTagBusinessUnit
	}
	if i.ResourceTagEnvironment != original.ResourceTagEnvironment {
		diff.AfterResourceTagEnvironment = i.ResourceTagEnvironment
		diff.BeforeResourceTagEnvironment = original.ResourceTagEnvironment
	}
	if i.ResourceTagOwner != original.ResourceTagOwner {
		diff.AfterResourceTagOwner = i.ResourceTagOwner
		diff.BeforeResourceTagOwner = original.ResourceTagOwner
	}
	if i.ResourceTagConfidentiality != original.ResourceTagConfidentiality {
		diff.AfterResourceTagConfidentiality = i.ResourceTagConfidentiality
		diff.BeforeResourceTagConfidentiality = original.ResourceTagConfidentiality
	}
	if i.ResourceTagIntegrity != original.ResourceTagIntegrity {
		diff.AfterResourceTagIntegrity = i.ResourceTagIntegrity
		diff.BeforeResourceTagIntegrity = original.ResourceTagIntegrity
	}
	if i.ResourceTagAvailability != original.ResourceTagAvailability {
		diff.AfterResourceTagAvailability = i.ResourceTagAvailability
		diff.BeforeResourceTagAvailability = original.ResourceTagAvailability
	}
	if i.ResourceTags != original.ResourceTags {
		diff.AfterResourceTags = i.ResourceTags
		diff.BeforeResourceTags = original.ResourceTags
	}
	return diff
}

// DeletedDiff produces a new ResourceReportItemDifference with fields set from the
// receiver ResourcesReportItem in the "before" columns, and the type set to
// DIFF_DELETED.
func (i ResourcesReportItem) DeletedDiff() ResourcesReportItemDifference {
	return ResourcesReportItemDifference{
		Type:                             DIFF_DELETED,
		ResourceARN:                      i.ResourceARN,
		BeforeResourceName:               i.ResourceName,
		BeforeResourceType:               i.ResourceType,
		BeforeResourceTagBusinessUnit:    i.ResourceTagBusinessUnit,
		BeforeResourceTagEnvironment:     i.ResourceTagEnvironment,
		BeforeResourceTagOwner:           i.ResourceTagOwner,
		BeforeResourceTagConfidentiality: i.ResourceTagConfidentiality,
		BeforeResourceTagIntegrity:       i.ResourceTagIntegrity,
		BeforeResourceTagAvailability:    i.ResourceTagAvailability,
		BeforeResourceTags:               i.ResourceTags,
	}
}

// AddedDiff produces a new ResourceReportItemDifference with fields set from the
// receiver ResourcesReportItem in the "after" columns, and the type set to
// DIFF_ADDED.
func (i ResourcesReportItem) AddedDiff() ResourcesReportItemDifference {
	return ResourcesReportItemDifference{
		Type:                            DIFF_ADDED,
		ResourceARN:                     i.ResourceARN,
		AfterResourceName:               i.ResourceName,
		AfterResourceType:               i.ResourceType,
		AfterResourceTagBusinessUnit:    i.ResourceTagBusinessUnit,
		AfterResourceTagEnvironment:     i.ResourceTagEnvironment,
		AfterResourceTagOwner:           i.ResourceTagOwner,
		AfterResourceTagConfidentiality: i.ResourceTagConfidentiality,
		AfterResourceTagIntegrity:       i.ResourceTagIntegrity,
		AfterResourceTagAvailability:    i.ResourceTagAvailability,
		AfterResourceTags:               i.ResourceTags,
	}
}
