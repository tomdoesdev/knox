## Development Workflow: Spec → Code

THESE INSTRUCTIONS ARE CRITICAL!

They dramatically improve the quality of the work you create.

### Phase 1: Requirements First

When asked to implement any feature or make changes, ALWAYS start by asking:
"Should I create a Spec for this task first?"

IFF user agrees:

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

ONLY after user types "GO!" or explicitly approves:

- Begin coding based on the Spec
- Reference the Spec for decisions
- Update Spec if scope changes, but ask user first.

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

**Remember: Think first, ask clarifying questions, _then_ code. The Spec is your north star.**

(source: https://lukebechtel.com/blog/vibe-speccing)