export interface User {
  username: string;
  createdAt: number;
}

export interface LoginResponse {
  message: string;
  token: string;
  username: string;
}

export interface AuthError {
  error: string;
}

export default class AuthService {
  static isAuthenticated(): boolean {
    return localStorage.getItem("token") !== null;
  }

  static getToken(): string | null {
    return localStorage.getItem("token");
  }

  static getUsername(): string | null {
    return localStorage.getItem("username");
  }

  static async login(
    username: string,
    password: string
  ): Promise<LoginResponse> {
    const response = await fetch("/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || "Login failed");
    }

    localStorage.setItem("token", data.token);
    localStorage.setItem("username", data.username);

    return data;
  }

  static async register(
    username: string,
    password: string
  ): Promise<{ message: string }> {
    const response = await fetch("/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || "Registration failed");
    }

    return data;
  }

  static async logout(): Promise<void> {
    try {
      await fetch("/auth/logout", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${this.getToken()}`,
        },
      });
    } catch (error) {
      console.error("Logout error:", error);
    } finally {
      localStorage.removeItem("token");
      localStorage.removeItem("username");
    }
  }

  static async deleteAccount(): Promise<{ message: string }> {
    const response = await fetch("/auth/delete-account", {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${this.getToken()}`,
      },
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.error || "Failed to delete account");
    }

    localStorage.removeItem("token");
    localStorage.removeItem("username");

    return data;
  }

  static async getCurrentUser(): Promise<User> {
    const response = await fetch("/auth/user", {
      headers: {
        Authorization: `Bearer ${this.getToken()}`,
      },
    });

    if (!response.ok) {
      throw new Error("Failed to get user information");
    }

    return response.json();
  }

  static async checkAuth(): Promise<boolean> {
    try {
      if (!this.getToken()) {
        return false;
      }

      const response = await fetch("/auth/check", {
        headers: {
          Authorization: `Bearer ${this.getToken()}`,
        },
      });

      if (!response.ok) {
        return false;
      }

      const data = await response.json();
      return data.authenticated === true;
    } catch (error) {
      console.error("Auth check error:", error);
      return false;
    }
  }
}
