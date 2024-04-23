module.exports = {
    branches: ["main"],
    plugins: [
        // Docs about Plugin Development: https://semantic-release.gitbook.io/semantic-release/developer-guide/plugin
        "@semantic-release/commit-analyzer",
        "@semantic-release/release-notes-generator",
        ["@semantic-release/github", {
            // https://github.com/semantic-release/github
            releasedLabels: false,
            successComment: false,
            failComment: false,
            failTitle: false,
        }],
        // https://github.com/semantic-release/exec
        ["@semantic-release/exec", {
            // verifyRelease.sh:
            //  * located at: dt-uibk-workshop/semantic-release-installer-action@main
            //  * automatically copied by related action.yml
            //  * TODO: create dedicated semantic-release plugin for this scripts
            verifyReleaseCmd: "./.github/verifyRelease.sh ${nextRelease.version} ${lastRelease.version} ${lastRelease.gitTag} ${lastRelease.gitHead} ${process.env.COMMIT_LINTING_ENABLED} ${process.env.VERIFY_DRYRUN_VERSION}",
            generateNotesCmd: "echo \"${nextRelease.notes}\" > .NOTES",
        }],
    ],
    preset: "conventionalcommits", // default types spec: https://github.com/conventional-changelog/conventional-changelog-config-spec/blob/master/versions/2.0.0/README.md#types
}
