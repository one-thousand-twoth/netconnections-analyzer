## Agent skills

### Issue tracker

Issues and specs live in GitHub Issues for `one-thousand-twoth/netconnections-analyzer`. See `docs/agents/issue-tracker.md`.

### Triage labels

Use the five default triage roles as GitHub labels. See `docs/agents/triage-labels.md`.

### Domain docs

Use the single-context layout: root `GLOSSARY.md` and `docs/adr/`. See `docs/agents/domain.md`.

## Commits

When an AI agent creates or substantively changes code, tests, configuration, or documentation included in a commit, append a trailer for that agent to the commit message:

```text
Co-authored-by: <Agent name> <agent-email@example.com>
```

Use the agent's project-designated co-author name and email (for Codex: `Co-authored-by: Codex <noreply@openai.com>`). If no attribution email is known, ask the user for one rather than inventing it. Separate the trailer block from the message body with a blank line. Include each contributing co-author identity exactly once. When editing a commit message, preserve existing co-author signatures and all other trailers. Keep the usual Git identity for the primary author and committer.

Attribute only changes the agent actually contributed to in that commit; merely being present in a session or running checks does not qualify. Whether to create a commit is governed by the user's instructions and the project's other rules.

When squash-merging a PR, inspect the final commit message before confirming the merge. Carry over the co-author identities of agents that contributed to changes in the resulting commit, once per identity, preserving other co-author signatures and trailers. See `docs/agents/issue-tracker.md` for the PR workflow.
