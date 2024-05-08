package dto

import (
	"time"
)

type Hook struct {
	HookName     string     `json:"hook_name"`
	Password     string     `json:"password"`
	HookID       int        `json:"hook_id"`
	HookURL      string     `json:"hook_url"`
	Timestamp    string     `json:"timestamp"`
	Sign         string     `json:"sign"`
	Comment      Comment    `json:"comment"`
	NoteableType string     `json:"noteable_type"`
	Issue        Issue      `json:"issue"`
	Repository   Repository `json:"repository"`
	Sender       Sender     `json:"sender"`
	Enterprise   Enterprise `json:"enterprise"`
}
type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Comment struct {
	HTMLURL   string    `json:"html_url"`
	ID        int       `json:"id"`
	Body      string    `json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Labels struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
type Assignee struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Milestone struct {
	HTMLURL        string      `json:"html_url"`
	ID             int         `json:"id"`
	Number         int         `json:"number"`
	Title          string      `json:"title"`
	Description    interface{} `json:"description"`
	OpenIssues     int         `json:"open_issues"`
	StartedIssues  int         `json:"started_issues"`
	ClosedIssues   int         `json:"closed_issues"`
	ApprovedIssues int         `json:"approved_issues"`
	State          string      `json:"state"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	DueOn          interface{} `json:"due_on"`
}
type Issue struct {
	HTMLURL   string    `json:"html_url"`
	ID        int       `json:"id"`
	Number    string    `json:"number"`
	Title     string    `json:"title"`
	User      User      `json:"user"`
	Labels    []Labels  `json:"labels"`
	State     string    `json:"state"`
	Assignee  Assignee  `json:"assignee"`
	Milestone Milestone `json:"milestone"`
	Comments  int       `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
}
type Owner struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Repository struct {
	ID                int         `json:"id"`
	Name              string      `json:"name"`
	Path              string      `json:"path"`
	FullName          string      `json:"full_name"`
	Owner             Owner       `json:"owner"`
	Private           bool        `json:"private"`
	HTMLURL           string      `json:"html_url"`
	URL               string      `json:"url"`
	Description       string      `json:"description"`
	Fork              bool        `json:"fork"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	PushedAt          time.Time   `json:"pushed_at"`
	GitURL            string      `json:"git_url"`
	SSHURL            string      `json:"ssh_url"`
	CloneURL          string      `json:"clone_url"`
	SvnURL            string      `json:"svn_url"`
	GitHTTPURL        string      `json:"git_http_url"`
	GitSSHURL         string      `json:"git_ssh_url"`
	GitSvnURL         string      `json:"git_svn_url"`
	Homepage          interface{} `json:"homepage"`
	StargazersCount   int         `json:"stargazers_count"`
	WatchersCount     int         `json:"watchers_count"`
	ForksCount        int         `json:"forks_count"`
	Language          string      `json:"language"`
	HasIssues         bool        `json:"has_issues"`
	HasWiki           bool        `json:"has_wiki"`
	HasPages          bool        `json:"has_pages"`
	License           interface{} `json:"license"`
	OpenIssuesCount   int         `json:"open_issues_count"`
	DefaultBranch     string      `json:"default_branch"`
	Namespace         string      `json:"namespace"`
	NameWithNamespace string      `json:"name_with_namespace"`
	PathWithNamespace string      `json:"path_with_namespace"`
}
type Sender struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	UserName  string `json:"user_name"`
	URL       string `json:"url"`
}
type Enterprise struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
