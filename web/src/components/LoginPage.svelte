<script lang="ts">
  interface AuthStore {
    login: (username: string, password: string) => Promise<boolean>;
    logout: () => void;
  }

  import { auth } from '../stores/auth';
  const typedAuth = auth as AuthStore;

  let username = '';
  let password = '';
  let error = '';
  let isLoading = false;

  async function handleSubmit() {
    isLoading = true;
    error = '';

    try {
      const success = await typedAuth.login(username, password);
      if (success) {
        // Use window.location for client-side navigation
        window.location.href = '/analytics';
      } else {
        error = 'Invalid credentials';
      }
    } catch (e) {
      error = 'Login failed. Please try again.';
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="login-container">
  <div class="login-card">
    <h1>Admin Login</h1>

    <form on:submit|preventDefault={handleSubmit} class="login-form">
      {#if error}
        <div class="error-message">{error}</div>
      {/if}

      <div class="form-group">
        <label for="username">Username</label>
        <input type="text" id="username" bind:value={username} required disabled={isLoading} />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" bind:value={password} required disabled={isLoading} />
      </div>

      <button type="submit" disabled={isLoading}>
        {isLoading ? 'Logging in...' : 'Login'}
      </button>
    </form>
  </div>
</div>

<style>
  .login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f3f4f6;
  }

  .login-card {
    background: white;
    padding: 2rem;
    border-radius: 0.5rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 400px;
  }

  h1 {
    text-align: center;
    color: #111827;
    margin-bottom: 2rem;
    font-size: 1.5rem;
    font-weight: 600;
  }

  .login-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  label {
    font-size: 0.875rem;
    font-weight: 500;
    color: #374151;
  }

  input {
    padding: 0.5rem;
    border: 1px solid #d1d5db;
    border-radius: 0.25rem;
    font-size: 1rem;
  }

  button {
    padding: 0.75rem;
    background-color: #047857;
    color: white;
    border: none;
    border-radius: 0.25rem;
    font-weight: 500;
    cursor: pointer;
  }

  button:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .error-message {
    padding: 0.75rem;
    background-color: #fee2e2;
    color: #dc2626;
    border-radius: 0.25rem;
    font-size: 0.875rem;
  }
</style>
