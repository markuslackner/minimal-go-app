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
    releaseRules: [
        // extending the default ReleaseRules: https://github.com/semantic-release/commit-analyzer/blob/master/lib/default-release-rules.js
        {type: 'build', release: 'patch'},
        {type: 'ci', release: 'patch'},
        {type: 'chore', release: 'patch'},
        {type: 'docs', release: 'patch'},
        {type: 'refactor', release: 'patch'},
        {type: 'style', release: 'patch'},
        {type: 'test', release: 'patch'},
    ],
    presetConfig: {
        types: [
            // overriding the default presetConfig: https://github.com/conventional-changelog/conventional-changelog/blob/master/packages/conventional-changelog-conventionalcommits/writer-opts.js#L181
            {type: 'feat', section: 'Features'},
            {type: 'feature', section: 'Features'},
            {type: 'fix', section: 'Bug Fixes'},
            {type: 'perf', section: 'Performance Improvements'},
            {type: 'revert', section: 'Reverts'},
            {type: 'docs', section: 'Documentation'},
            {type: 'style', section: 'Styles'},
            {type: 'chore', section: 'Miscellaneous Chores'},
            {type: 'refactor', section: 'Code Refactoring'},
            {type: 'test', section: 'Tests'},
            {type: 'build', section: 'Build System'},
            {type: 'ci', section: 'Continuous Integration'}
        ]
    }
}