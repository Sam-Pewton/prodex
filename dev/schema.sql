CREATE TABLE IF NOT EXISTS jira (
    issue_key TEXT NOT NULL PRIMARY KEY, -- key
    url TEXT, -- self
    issue_type TEXT, -- fields.issuetype.name
    issue_summary TEXT, --fields.summary
    issue_description TEXT, --fields.description.content
    issue_status TEXT, --fields.status.name
    issue_priority TEXT, -- fields.priority.name
    issue_resolution TEXT, -- fields.resolution.name
    assignee TEXT, -- fields.assignee.displayName
    creator TEXT, -- fields.creator.displayName
    reporter TEXT, -- fields.reporter.displayName
    parent_key TEXT, -- fields.parent.key
    parent_type TEXT, -- fields.parent.fields.issuetype.name
    parent_summary TEXT, -- fields.parent.fields.summary
    parent_status TEXT, -- fields.parent.fields.status.name
    parent_priority TEXT, -- fields.parent.fields.priority.name
    project_key TEXT, -- fields.project.key
    project_name TEXT, -- fields.project.name
    issue_created_date DATETIME, -- fields.created
    issue_updated_date DATETIME, -- fields.updated
    issue_resolution_date DATETIME, -- fields.resolutiondate
    last_status_update_date DATETIME -- fields.statuscategorychangedate
);

CREATE TABLE IF NOT EXISTS confluence (
    id INTEGER NOT NULL PRIMARY KEY,
    url TEXT,
    title TEXT,
    space TEXT NOT NULL,
    created_date DATETIME NOT NULL,
    page_owner TEXT NOT NULL -- ownerId maps to the author
);

CREATE TABLE IF NOT EXISTS github (
    id INTEGER NOT NULL PRIMARY KEY
);
