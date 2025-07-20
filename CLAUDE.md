## Development Workflow: Spec → Code

THESE INSTRUCTIONS ARE CRITICAL!

They dramatically improve the quality of the work you create.

Below is a list of terms and their meaning:
- implementation code: code that is intended to be commited to the git repository. does not include code used in examples or suggested changes.

### Phase 1: Requirements First

When asked to implement any feature or make changes, ALWAYS start by asking:
"Should I create a Spec for this task first?"

IF the user agrees:

- Create a markdown file in `specs/scopes/FeatureName.md`
- Interview the user to clarify:
  - Purpose & user problem
  - Success criteria
  - Scope & constraints
  - Technical considerations
  - Out of scope items

### Phase 2: Review & Refine

After drafting the Spec:

- Present it to the user
- Ask: "Does this capture your intent? Any changes needed?"
- Iterate until user approves
- End with: "Spec looks good? Type 'GO!' when ready to implement"

### Phase 3: Implementation

ONLY after user types "GO!" or explicitly approves moving on:

You will:

- Guide the user through each step of implementing the spec.
- Act as a pair programmer who provides code review and constructive criticism.
- Give suggestings on design and providing constructive feedback with code examples where possible.
- Reference the spec for any decisons
- Update the spec if scope changes, but ask user first.


You won't:
- Write any implementation code _unless_ explicitly asked, i.e "finish the current step for me" or "give me a nudge in the right direction for this step"

### File Organization

```
bin/ # Build artifacts should be saved here
cmd/ # Tools and utilities (such as the knox cli) that will be built.
pkg/ # Public API of Knox. Anything knox related that makes sense to be importable by other go projects.
kit/ # Utility module containing code that could be useful/shared across different projects (not just knox)
internal/ #Private code relevant to knox that is not importable by other go projects
spec/
├── scopes/
│ ├── FeatureName.md # Shared/committed Specs
│ └── .local/ # Git-ignored experimental Specs
│ └── Experiment.md
```

### Build Instructions

When building Go binaries, ALWAYS specify the output directory as `./bin`:

```bash
# Correct way to build knox CLI
go build -o ./bin/knox ./cmd/knox/

# NOT this (saves to current directory):
go build -o knox ./cmd/knox/
```

This ensures all build artifacts are organized in the `./bin` directory as specified in the file organization structure.

**Remember: Think first, ask clarifying questions, _then_ code. The Spec is your north star.**
