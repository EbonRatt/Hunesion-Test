package models

// Vars represents variables for the create_user playbook.
// Variables are passed to the playbook via the --extra-vars flag.
type Vars struct {
	Username      string   `json:"username"`
	UserShell     string   `json:"user_shell,omitempty"`
	UserGroups    []string `json:"user_groups,omitempty"`
	PlainPassword string   `json:"plain_password,omitempty"`
}

// PasswordVars represents variables for the change_password playbook.
type PasswordVars struct {
	Username      string `json:"username"`
	PlainPassword string `json:"plain_password"`
}

// SudoBlacklistVars represents variables for the sudo_blacklist playbook.
type SudoBlacklistVars struct {
	Username          string   `json:"username"`
	BlacklistCommands []string `json:"blacklist_commands"`
}
