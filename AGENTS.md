# Project instructions for every Active Agent

catty uses the two-party, model-neutral protocol in
[`docs/COLLABORATION.md`](docs/COLLABORATION.md). The two parties are LizBing
(Project Owner) and one Active Agent. A model name is not a project role.

Before doing project work, read in this order:

1. [`docs/PROJECT_STATUS.md`](docs/PROJECT_STATUS.md)
2. [`docs/COLLABORATION.md`](docs/COLLABORATION.md)
3. The active workstream linked by `PROJECT_STATUS.md`, if one exists
4. Governing Accepted ADRs
5. [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) and relevant development docs

The Owner's current request determines whether the session is discussion,
read-only review, implementation, or integration. Do not expand that authority.
Without an Accepted workstream, do not begin non-trivial implementation.
Proposed ADRs and Roadmap entries do not authorize implementation.

Repository evidence outranks chat history. Do not report a capability complete
without its acceptance evidence. Do not commit, merge, push, publish, rewrite
remote state, or discard unrelated changes unless the Owner explicitly authorizes
that integration action.
