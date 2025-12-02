# Contributing to LaTeX PPT Generator

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/latex_ppt_by_claude.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes thoroughly
6. Commit and push to your fork
7. Create a Pull Request

## Development Setup

See [SETUP.md](SETUP.md) for detailed setup instructions.

Quick start for development:
```bash
# Start infrastructure services
make dev

# In separate terminals:
make backend-dev
make frontend-dev
```

## Project Structure

```
├── backend/           # Go backend
│   ├── cmd/          # Application entrypoints
│   ├── internal/     # Private application code
│   │   ├── api/     # HTTP handlers and routing
│   │   ├── service/ # Business logic
│   │   ├── repository/ # Data access
│   │   ├── model/   # Data models
│   │   └── config/  # Configuration
│   └── pkg/         # Public libraries
│       ├── ai/      # AI client implementations
│       ├── embedding/ # Embedding generation
│       ├── vectordb/ # Vector database client
│       ├── parser/  # Document parsers
│       └── latex/   # LaTeX compilation
├── frontend/        # Vue3 frontend
│   └── src/
│       ├── api/     # API clients
│       ├── components/ # Vue components
│       ├── views/   # Page components
│       ├── store/   # Pinia stores
│       ├── router/  # Vue Router config
│       └── utils/   # Utility functions
└── docker/          # Docker configurations
```

## Coding Standards

### Backend (Go)

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format code
- Run `go vet` before committing
- Add comments for exported functions and types
- Keep functions focused and small
- Use meaningful variable names

Example:
```go
// GeneratePPT creates a new PPT based on user prompt and optional knowledge base
func (s *PPTService) GeneratePPT(ctx context.Context, userID uint, prompt string) (*model.PPTRecord, error) {
    // Implementation
}
```

### Frontend (Vue3/TypeScript)

- Follow Vue 3 Composition API best practices
- Use TypeScript for type safety
- Use Composition API with `<script setup>`
- Keep components small and focused
- Use Pinia for state management
- Follow Element Plus component patterns

Example:
```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { User } from '@/types'

const user = ref<User | null>(null)

onMounted(async () => {
  // Implementation
})
</script>
```

### General

- Write clear commit messages
- Add comments for complex logic
- Update documentation when changing functionality
- Write tests for new features (when applicable)

## Adding New Features

### Backend Feature

1. **Add Models** (if needed)
   - Create in `backend/internal/model/`
   - Add GORM tags for database mapping

2. **Add Repository Methods**
   - Create/update repository in `backend/internal/repository/`
   - Add CRUD operations

3. **Add Service Logic**
   - Create/update service in `backend/internal/service/`
   - Implement business logic

4. **Add API Handler**
   - Create handler in `backend/internal/api/handler/`
   - Handle HTTP requests/responses

5. **Update Router**
   - Add route in `backend/internal/api/router.go`

6. **Update Documentation**
   - Add API endpoint to `API.md`

### Frontend Feature

1. **Add Types**
   - Define TypeScript types in `frontend/src/types/`

2. **Add API Client**
   - Create API functions in `frontend/src/api/`

3. **Add Store** (if needed)
   - Create Pinia store in `frontend/src/store/`

4. **Add Component/View**
   - Create Vue component in appropriate directory
   - Use Composition API with TypeScript

5. **Update Router** (if new page)
   - Add route in `frontend/src/router/index.ts`

6. **Update Documentation**
   - Update `README.md` with feature description

## Testing

### Backend Testing

```bash
cd backend
go test ./...
```

### Frontend Testing

Currently, no automated tests are set up. Manual testing is required:
1. Test all UI interactions
2. Test API integration
3. Test error handling
4. Test on different screen sizes

### Integration Testing

1. Start all services: `make up`
2. Test complete workflows:
   - User registration and login
   - Document upload and processing
   - PPT generation
   - PDF compilation and download

## Common Development Tasks

### Adding a New AI Provider

1. Create client in `backend/pkg/ai/`
2. Implement `GenerateLaTeX` method
3. Update `AIService` to use new provider
4. Add configuration in `.env.example`
5. Update documentation

### Adding a New Document Parser

1. Create parser in `backend/pkg/parser/`
2. Implement `Parse` method
3. Update `GetParser` function
4. Add file type to upload validation
5. Update documentation

### Adding a New LaTeX Template

1. Add template string in `backend/pkg/latex/templates.go`
2. Add template name to `ListTemplates()`
3. Update frontend template selector
4. Document template features

## Pull Request Guidelines

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] Code compiles without errors
- [ ] Manual testing completed
- [ ] Documentation updated
- [ ] Commit messages are clear

### PR Description

Include:
1. **What**: Description of changes
2. **Why**: Reason for changes
3. **How**: Implementation approach
4. **Testing**: How you tested the changes
5. **Screenshots**: For UI changes

Example:
```markdown
## Add support for Markdown tables in PPT generation

### What
Added support for parsing and rendering Markdown tables in generated LaTeX presentations.

### Why
Users requested ability to include tabular data in presentations.

### How
- Extended Markdown parser to detect tables
- Added LaTeX `tabular` environment generation
- Updated AI prompt to include table formatting instructions

### Testing
- Tested with various table sizes
- Verified LaTeX compilation
- Checked PDF output quality

### Screenshots
[Attach screenshots of generated tables]
```

## Code Review Process

1. Maintainer reviews PR
2. Changes requested or approved
3. After approval, PR is merged
4. Branch is deleted

## Reporting Bugs

Use GitHub Issues with the following information:

```markdown
**Description**
Clear description of the bug

**Steps to Reproduce**
1. Go to '...'
2. Click on '...'
3. See error

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- OS: [e.g., Ubuntu 22.04]
- Docker version: [e.g., 24.0.0]
- Browser: [e.g., Chrome 120]

**Logs**
Relevant error logs or screenshots
```

## Requesting Features

Use GitHub Issues with:

```markdown
**Feature Description**
Clear description of the proposed feature

**Use Case**
Why this feature would be useful

**Proposed Implementation**
(Optional) Ideas on how to implement

**Alternatives Considered**
Other solutions you've thought about
```

## Community

- Be respectful and constructive
- Help others when you can
- Share knowledge and experiences
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md) (if exists)

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT License).

## Questions?

- Check existing issues and documentation
- Ask in GitHub Discussions (if enabled)
- Create a new issue with the "question" label

## Thank You!

Your contributions help make this project better for everyone!
