# Contributing

## Getting started

Contributions are made to this repo via [Issues][RepoIssues] and [Pull Requests (PRs)][PRs].

Commits to this repo are not required to be signed, however it is recommended.
For more information on how to do this, please visit [git][GitSigning] to sign
commits, and [Github][GithubSigning] to verify them.

A few general guidelines that to keep in mind:

* Search for existing Issues and PRs before creating your own.
* Use labels to draw attention to why something should be worked on.

### Issues

Issues should be used to record bugs, and request additional features, or discuss
potential changes before it's implemented.

If you find an reported bug that you have experienced and you want fixed, adding
a [reaction][GithubReaction] can help prioritisation.

Once an issue is raised, it should have appropriate labels to allow for easy discovery.
Maintainers will add milestones to the issue to indicate the priority and timeline for the change.

When starting to work on an issue, please assign it to yourself to let others know
that you are working on it.

### Pull requests

PRs will always be welcome, and is a great way to quickly make sure your change makes
its way in.

A PR should:

* Add unit or e2e tests for fixed or changed functionality.
* Include documentation.
* Be accompanied by a completed Pull Request template.
* Bump the version.

Branch names should be descriptive, and be prefixed with the issue number.
All commits related to the issue should also be prefixed with the issue number.

When a PR is merged, it will automatically squash the commits into the main branch.

### Tests

The unit tests should be used to assert that the functions and individual pieces of functionality
logically work as expected.
These will include mocking out the lower levels of the application, for example the service layer
tests will mock the database layer's responses.

End to end tests will run the executable, and assert on the output and changes around that you
can observe from the command line.
These are ran in a linux environment.
To run these tests, you'll need [commander][CommanderCLI] and [jq][JqCLI] installed and accessible
on the path.

[RepoIssues]: https://github.com/PossibleLlama/worklog/issues
[PRs]: https://github.com/PossibleLlama/worklog/pulls
[GitSigning]: https://git-scm.com/book/en/v2/Git-Tools-Signing-Your-Work
[GithubSigning]: https://docs.github.com/en/github/authenticating-to-github/signing-commits
[GithubReaction]: https://github.blog/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/
[CommanderCLI]: https://github.com/commander-cli/commander
[JqCLI]: https://github.com/stedolan/jq
