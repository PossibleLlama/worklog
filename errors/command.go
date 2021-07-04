package errors

// RootRepoType error value for invalid type of repo
const RootRepoType = "invalid repo, either \"local\" or \"legacy\" must be specified"

// ConfigureArgsMinimum error value when not enough args
const ConfigureArgsMinimum = "overrideDefaults requires at least one argument"

// EditID error value when requires an ID
const EditID = "edit requires a single ID of an existing worklog"

// PrintID error value when requires an ID
const PrintID = "no ids provided"

// PrintArgsMinimum error value when not enough args
const PrintArgsMinimum = "one flag is required"

// Format error value when wrong format
const Format = "format is not valid"
