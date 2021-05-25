package errors

const Unexpected = "An unexpected error occurred"

const (
	RepoCreateDirectory  = "Unable to create directory"
	RepoCreateFile       = "unable to create file"
	RepoConfigFileCreate = "unable to create configuration file"
	RepoConfigFileSave   = "unable to save config"
	RepoSaveFile         = "unable to save worklog"
)

const (
	RepoGetFiles     = "unable to get all files"
	RepoGetFilesRead = "error reading file"
)
