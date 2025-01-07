# Development Standards

## Technology Stack

- Frontend: Astro with Svelte 5 islands
- Styling: Vanilla CSS (no Tailwind)
- Backend: Go
- Types: TypeScript with strict typing

## TypeScript Standards

- Use strict type checking
- Avoid `any` type
- Generate interfaces for missing types
- Use proper type guards
- Prefer type inference where possible

## Code Organization

- Use feature-based folder structure
- Keep components single-responsibility
- Maintain clear type definitions
- Follow proper naming conventions

## Style Guidelines

- Use vanilla CSS with CSS variables
- Maintain consistent class naming
- Follow BEM methodology
- Use responsive design patterns

## Best Practices

- Always use #codebase context
- Write self-documenting code
- Include proper error handling
- Add type documentation
- Follow SOLID principles

## Example Patterns

```typescript
// Component props interface
interface ComponentProps {
  data: DataType;
  onAction: (id: string) => Promise<void>;
}

// Type guard example
function isValidResponse(data: unknown): data is ApiResponse {
  return typeof data === "object" && data !== null && "id" in data;
}
```
