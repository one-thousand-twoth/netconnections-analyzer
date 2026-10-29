## Agent skills

### Issue tracker

Issues and specs live in GitHub Issues for `one-thousand-twoth/netconnections-analyzer`. See `docs/agents/issue-tracker.md`.

### Triage labels

Use the five default triage roles as GitHub labels. See `docs/agents/triage-labels.md`.

### Domain docs

Use the single-context layout: root `GLOSSARY.md` and `docs/adr/`. See `docs/agents/domain.md`.

## Commits

When Codex creates or substantively changes code, tests, configuration, or documentation included in a commit, append this trailer to that commit message:

```text
Co-authored-by: Codex <noreply@openai.com>
```

Separate the trailer block from the message body with a blank line. Include the Codex trailer exactly once. When editing a commit message, preserve existing co-author signatures and all other trailers. Keep the usual Git identity for the primary author and committer.

Attribute only changes Codex actually contributed to in that commit; merely being present in a session or running checks does not qualify. Whether to create a commit is governed by the user's instructions and the project's other rules.

When squash-merging a PR, inspect the final commit message before confirming the merge. If Codex contributed to changes in the resulting commit, carry its trailer into that message once, preserving other co-author signatures and trailers. See `docs/agents/issue-tracker.md` for the PR workflow.
