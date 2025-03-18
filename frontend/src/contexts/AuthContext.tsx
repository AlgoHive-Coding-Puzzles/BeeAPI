import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";
import AuthService from "../services/AuthService";

interface AuthContextType {
  isAuthenticated: boolean;
  username: string | null;
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  deleteAccount: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState(
    AuthService.isAuthenticated()
  );
  const [username, setUsername] = useState<string | null>(
    AuthService.getUsername()
  );
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkAuthentication = async () => {
      try {
        // Only verify with backend if we have a token
        if (AuthService.getToken()) {
          const authenticated = await AuthService.checkAuth();
          if (!authenticated) {
            // If backend says token is invalid, clear it
            AuthService.logout();
            setIsAuthenticated(false);
            setUsername(null);
          }
        }
      } catch (error) {
        console.error("Auth check failed:", error);
        setIsAuthenticated(false);
        setUsername(null);
      } finally {
        setLoading(false);
      }
    };

    checkAuthentication();
  }, []);

  const login = async (username: string, password: string) => {
    await AuthService.login(username, password);
    setIsAuthenticated(true);
    setUsername(username);
  };

  const register = async (username: string, password: string) => {
    await AuthService.register(username, password);
    await login(username, password);
  };

  const logout = async () => {
    await AuthService.logout();
    setIsAuthenticated(false);
    setUsername(null);
  };

  const deleteAccount = async () => {
    await AuthService.deleteAccount();
    setIsAuthenticated(false);
    setUsername(null);
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        username,
        login,
        register,
        logout,
        deleteAccount,
      }}
    >
      {!loading && children}
    </AuthContext.Provider>
  );
};
