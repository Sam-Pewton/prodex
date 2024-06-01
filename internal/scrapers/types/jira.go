package types

type JiraIssueResponse struct {
	Issues []JiraIssue `json:"issues"`
}

type JiraUser struct {
	AccountId   string `json:"accountId"`
	DisplayName string `json:"displayName"` // needed `issue_assignee` `issue_creator` `issue_reporter`
}

type JiraIssue struct {
	Id     string          `json:"id"`   // needed, Jira identifier `issue_id`
	Key    string          `json:"key"`  // needed, "SDWAN-x" `issue_key`
	Self   string          `json:"self"` // needed, the url of the story? to check `issue_self`
	Fields JiraIssueFields `json:"fields"`
}

type JiraParent struct {
	Id     string           `json:"id"`  // needed, Jira identifier `parent_id`
	Key    string           `json:"key"` // needed, "SDWAN-x" `parent_key`
	Fields JiraParentFields `json:"fields"`
}

type JiraPriority struct {
	Name string `json:"name"` // needed, "Medium" `issue_priority` `parent_priority`
}

type JiraParentFields struct {
	Summary   string        `json:"summary"` // needed `parent_summary`
	Status    JiraStatus    `json:"status"`
	Priority  JiraPriority  `json:"priority"`
	IssueType JiraIssueType `json:"issuetype"`
}

type JiraIssueFields struct {
	StatusCategoryChangeDate string          `json:"statuscategorychangedate"` // needed, the last status change `last_status_update_date`
	IssueType                JiraIssueType   `json:"issuetype"`
	Parent                   JiraParent      `json:"parent"`
	Project                  JiraProject     `json:"project"`
	Created                  string          `json:"created"` // needed, the date this was created `created_date`
	Priority                 JiraPriority    `json:"priority"`
	IssueLinks               []string        `json:"issuelinks"` // needed, the linked issues - might not work as a string `issue_links`
	Assignee                 JiraUser        `json:"assignee"`
	Updated                  string          `json:"updated"` // needed `last_updated_date`
	Status                   JiraStatus      `json:"status"`
	Summary                  string          `json:"summary"` //needed `issue_summary`
	Creator                  JiraUser        `json:"creator"`
	Reporter                 JiraUser        `json:"reporter"`
	Resolution               JiraResolution  `json:"resolution"`
	ResolutionDate           string          `json:"resolutiondate"` // needed, `resolution_date`
	Description              JiraDescription `json:"description"`
}

type JiraIssueType struct {
	Id          string `json:"id"`          // I think this is the ID of the issue type `issue_type_id`
	Description string `json:"description"` // could add a table describing these things `issue_description` -- TO DO
	Name        string `json:"name"`        // needed, "Story, Task, Bug" `issue_type` `parent_type`
	// Subtask bool `json:"subtask"` // maybe
	// EntityID string `json:"entityId"`
	// HierarchyLevel uint8 `json:"hierarchyLevel"`
}

type JiraProject struct {
	Id   string `json:"id"`   // needed `project_id`
	Key  string `json:"key"`  // needed, the project key `project_key`
	Name string `json:"name"` //needed, the project name `project_name`
	// ProjectTypeKey string `json:"projectTypeKey"`
}

type JiraResolution struct {
	Name string `json:"name"` // needed, `issue_resolution`
}

type JiraStatus struct {
	Name string `json:"name"` // needed, `parent_status` `issue_status`
}

type JiraDescription struct {
	Type    string   `json:"type"`
	Content []string `json:"description"`
}
